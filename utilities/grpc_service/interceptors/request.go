package interceptors

import (
	"context"

	"github.com/google/uuid"

	"geolize/utilities/contexts"
	jsonhelper "geolize/utilities/json_helper"
	"geolize/utilities/logging"

	"google.golang.org/grpc"
)

func RequestInterceptor(logger logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		reqCtx := contexts.NewServerContext(ctx)

		if len(contexts.GetServerData(reqCtx).RequestID) < 1 {
			xRequestID, err := uuid.NewV7()
			if err != nil {
				logger.Error(reqCtx, "Error: cannot generate request id with UUIDV7", logging.NewError(err)...)
				xRequestID = uuid.New()
			}
			contexts.SetRequestID(reqCtx, xRequestID.String())
		}

		logger.Info(reqCtx, "Request headers",
			logging.NewKeyVal("in-md", contexts.GetServerData(reqCtx).IncomingHeaders))

		// Create logger with request ID field
		requestLogger := logger.WithFields(
			logging.NewKeyVal("request_id", contexts.GetServerData(reqCtx).RequestID))

		requestLogger.Info(reqCtx, "Incoming Request", logging.NewKeyVal("api", info.FullMethod), logging.NewKeyVal("request", jsonhelper.ToString(req)))

		resp, err := handler(reqCtx, req)
		if err != nil {
			requestLogger.Error(reqCtx, "Error",
				logging.NewError(err)...)
			return nil, err
		}

		requestLogger.Info(reqCtx, "Response", logging.NewKeyVal("respone", jsonhelper.ToString(resp)))

		return resp, nil
	}
}
