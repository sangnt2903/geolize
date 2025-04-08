package iplocation

import (
	"context"
	"geolize/services/geolize/internal/pkg/ip_location/model"
	"geolize/services/geolize/internal/pkg/ip_location/providers/maxmind"
	"geolize/utilities/logging"
)

type IPGeolocate interface {
	Lookup(ctx context.Context, request *model.IPLookupRequest) ([]*model.IPResult, error)
	Update(ctx context.Context, request *model.IPUpdateRequest) error
}

func NewIPGeolocate(logger logging.Logger) IPGeolocate {
	return maxmind.New(logger)
}
