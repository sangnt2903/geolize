package cmd

import (
	"geolize/services/geolize/internal/handler"
	iplocation "geolize/services/geolize/internal/pkg/ip_location"
	"geolize/utilities/grpc_service"
	"geolize/utilities/logging"

	"geolize/service-protos/generated/geolize/geolize_pb"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	port int
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Geolize service",
	Long: `Start the Geolize service which provides IP geolocation services.
The grpc_service will listen on the specified port and handle API requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	logger, err := logging.NewLogger(logging.ZapLoggerType)
	if err != nil {
		panic(err)
	}

	ipLocation := iplocation.NewIPGeolocate(logger)
	service := handler.NewService(logger, ipLocation)

	var register grpc_service.GrpcRegister = func(s *grpc.Server) {
		geolize_pb.RegisterGeolizeServer(s, service)
	}

	s := grpc_service.New(logger, register, geolize_pb.RegisterGeolizeHandlerFromEndpoint)

	if err = s.Run(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
