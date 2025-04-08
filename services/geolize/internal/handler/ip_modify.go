package handler

import (
	"context"
	"errors"
	"geolize/service-protos/generated/geolize/geolize_pb"
	"geolize/services/geolize/internal/pkg/ip_location/model"
)

func (s Service) ModifyIP(ctx context.Context, request *geolize_pb.ModifyIPRequest) (*geolize_pb.ModifyIPResponse, error) {
	if len(request.Ip) == 0 {
		return nil, errors.New("invalid IP")
	}

	err := s.ipLocation.Update(ctx, &model.IPUpdateRequest{
		IP: request.Ip,
		Continent: func() *model.Continent {
			if request.Continent == nil {
				return nil
			}
			return &model.Continent{
				Code:  request.Continent.Code,
				Names: request.Continent.Names,
			}
		}(),
		Country: func() *model.Country {
			if request.Country == nil {
				return nil
			}
			return &model.Country{
				ISOCode:           request.Country.IsoCode,
				Names:             request.Country.Names,
				IsInEuropeanUnion: request.Country.IsInEuropeanUnion,
			}
		}(),
		Subdivisions: func() []*model.Subdivision {
			if request.GetSubdivisions() == nil {
				return nil
			}
			var subdivisions []*model.Subdivision
			for _, subdivision := range request.Subdivisions {
				subdivisions = append(subdivisions, &model.Subdivision{
					ISOCode: subdivision.IsoCode,
					Names:   subdivision.Names,
				})
			}
			return subdivisions
		}(),
		Location: func() *model.Location {
			if request.Location == nil {
				return nil
			}
			return &model.Location{
				Latitude:       request.Location.Latitude,
				Longitude:      request.Location.Longitude,
				AccuracyRadius: uint16(request.Location.AccuracyRadius),
				TimeZone:       request.Location.TimeZone,
			}
		}(),
		Postal: func() *model.Postal {
			if request.Postal == nil {
				return nil
			}
			return &model.Postal{
				Code: request.Postal.Code,
			}
		}(),
		City: func() *model.City {
			if request.City == nil {
				return nil
			}
			return &model.City{
				Names: request.City.Names,
			}
		}(),
		RepresentedCountry: func() *model.RepresentedCountry {
			if request.RepresentedCountry == nil {
				return nil
			}
			return &model.RepresentedCountry{
				ISOCode:           request.RepresentedCountry.IsoCode,
				Names:             request.RepresentedCountry.Names,
				IsInEuropeanUnion: request.RepresentedCountry.IsInEuropeanUnion,
			}
		}(),
		RegisteredCountry: func() *model.RegisteredCountry {
			if request.RegisteredCountry == nil {
				return nil
			}
			return &model.RegisteredCountry{
				ISOCode:           request.RegisteredCountry.IsoCode,
				Names:             request.RegisteredCountry.Names,
				IsInEuropeanUnion: request.RegisteredCountry.IsInEuropeanUnion,
			}
		}(),
		Traits: func() *model.Traits {
			if request.Traits == nil {
				return nil
			}
			return &model.Traits{
				IsAnonymousProxy:    request.Traits.IsAnonymousProxy,
				IsAnycast:           request.Traits.IsAnycast,
				IsSatelliteProvider: request.Traits.IsSatelliteProvider,
			}
		}(),
	})
	if err != nil {
		return nil, err
	}

	return &geolize_pb.ModifyIPResponse{}, nil
}
