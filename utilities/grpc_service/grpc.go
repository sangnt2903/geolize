package grpc_service

import (
	"context"
	"fmt"
	"geolize/utilities/grpc_service/interceptors"
	"geolize/utilities/logging"
	"geolize/utilities/service"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	gatewayRegister GatewayRegister
	register        GrpcRegister

	httpServer *http.Server
	server     *grpc.Server

	logger logging.Logger
}

type GrpcRegister func(server *grpc.Server)
type GatewayRegister func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

func New(logger logging.Logger, register GrpcRegister, gwRegister GatewayRegister) *Server {
	return &Server{
		gatewayRegister: gwRegister,
		register:        register,
		logger:          logger,
	}
}

const (
	contentTypeName  = "content-type"
	contentTypeValue = "application/grpc"
)

func (s *Server) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", service.GetPort()))
	if err != nil {
		s.logger.Fatal(context.Background(), func() string {
			return "Failed to listen: " + err.Error()
		}())
	}

	if s.gatewayRegister != nil {
		m := cmux.New(l)
		httpL := m.Match(cmux.HTTP1Fast())
		grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings(contentTypeName, contentTypeValue))

		go func() {
			err := s.httpRun(httpL)
			if err != nil {
				return
			}
		}()
		go func() {
			err := s.grpcRun(grpcL)
			if err != nil {
				panic(err)
			}
		}()

		return m.Serve()
	}

	return s.grpcRun(l)
}

func (s *Server) grpcRun(l net.Listener) error {
	var unaryInterceptors = []grpc.UnaryServerInterceptor{
		interceptors.RequestInterceptor(s.logger),
	}
	var streamInterceptors []grpc.StreamServerInterceptor

	s.server = grpc.NewServer(grpc.ChainUnaryInterceptor(unaryInterceptors...), grpc.ChainStreamInterceptor(streamInterceptors...))
	s.register(s.server)

	s.introduce()
	return s.server.Serve(l)
}

func (s *Server) httpRun(l net.Listener) error {
	gwMux := runtime.NewServeMux(
		runtime.WithErrorHandler(ErrorHandler),
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	//swagger
	mux := http.NewServeMux()
	mux.Handle("/", gwMux)

	fs := http.FileServer(http.Dir("/var/lib/swagger-ui"))
	mux.Handle("/docs/", http.StripPrefix("/docs", fs))

	mux.HandleFunc("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./swagger.json")
	})

	err := s.gatewayRegister(context.Background(), gwMux, fmt.Sprintf(":%d", service.GetPort()), opts)
	if err != nil {
		return err
	}

	// Add CORS support
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins, change as needed
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)

	s.httpServer = &http.Server{Handler: corsHandler}
	return s.httpServer.Serve(l)
}

func (s *Server) introduce() {
	s.logger.Info(context.Background(), func() string {
		return fmt.Sprintf("Server is running on :%d", service.GetPort())
	}())
}
