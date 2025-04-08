package handler

import (
	"context"
	"geolize/service-protos/generated/geolize/geolize_pb"
	iplocation "geolize/services/geolize/internal/pkg/ip_location"
	"geolize/utilities/logging"
)

type Service struct {
	//geolize_pb.UnimplementedGeolizeServer
	logger     logging.Logger
	ipLocation iplocation.IPGeolocate
}

func (s Service) Ping(ctx context.Context, request *geolize_pb.PingRequest) (*geolize_pb.PingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewService(logger logging.Logger, ipLocation iplocation.IPGeolocate) *Service {
	return &Service{
		logger:     logger,
		ipLocation: ipLocation,
	}
}
