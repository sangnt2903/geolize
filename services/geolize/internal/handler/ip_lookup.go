package handler

import (
	"context"
	"errors"
	"geolize/service-protos/generated/geolize/geolize_pb"
	"geolize/services/geolize/internal/pkg/ip_location/model"
	"geolize/services/geolize/internal/pkg/transform_response"
	"geolize/utilities/logging"
)

func (s Service) LookupIP(ctx context.Context, request *geolize_pb.LookupIPRequest) (*geolize_pb.LookupIPResponse, error) {
	if len(request.Ips) < 1 {
		return nil, errors.New("IPs are required")
	}

	resp, err := s.ipLocation.Lookup(ctx, &model.IPLookupRequest{
		IPs: request.GetIps(),
	})
	if err != nil {
		s.logger.Error(ctx, "ipLocation.Lookup", logging.NewError(err)...)
		return nil, err
	}

	return &geolize_pb.LookupIPResponse{
		Data: transform_response.ToLookupIPsResponse(resp),
	}, nil
}
