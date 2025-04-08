package transform_response

import (
	"geolize/service-protos/generated/geolize/geolize_pb"
	"geolize/services/geolize/internal/pkg/ip_location/model"
)

func ToLookupIPsResponse(ipResults []*model.IPResult) []*geolize_pb.IPInfo {
	var ipInfos []*geolize_pb.IPInfo
	for _, ipResult := range ipResults {
		ipInfos = append(ipInfos, &geolize_pb.IPInfo{
			Ip:        ipResult.IP,
			DbVersion: ipResult.DBVersion,
			City: &geolize_pb.City{
				Names: ipResult.City.Names,
			},
			Location: &geolize_pb.Location{
				Latitude:       ipResult.Location.Latitude,
				Longitude:      ipResult.Location.Longitude,
				AccuracyRadius: uint32(ipResult.Location.AccuracyRadius),
				TimeZone:       ipResult.Location.TimeZone,
			},
			Continent: &geolize_pb.Continent{
				Code:  ipResult.Continent.Code,
				Names: ipResult.Continent.Names,
			},
			Country: &geolize_pb.Country{
				IsoCode:           ipResult.Country.ISOCode,
				Names:             ipResult.Country.Names,
				IsInEuropeanUnion: ipResult.Country.IsInEuropeanUnion,
			},
			Subdivisions: func() []*geolize_pb.Subdivision {
				var subdivisions []*geolize_pb.Subdivision
				for _, subdivision := range ipResult.Subdivisions {
					subdivisions = append(subdivisions, &geolize_pb.Subdivision{
						IsoCode: subdivision.ISOCode,
						Names:   subdivision.Names,
					})
				}
				return subdivisions
			}(),
			RepresentedCountry: &geolize_pb.RepresentedCountry{
				IsoCode:           ipResult.RepresentedCountry.ISOCode,
				Type:              ipResult.RepresentedCountry.Type,
				Names:             ipResult.RepresentedCountry.Names,
				IsInEuropeanUnion: ipResult.RepresentedCountry.IsInEuropeanUnion,
			},
			RegisteredCountry: &geolize_pb.RegisteredCountry{
				IsoCode:           ipResult.RegisteredCountry.ISOCode,
				Names:             ipResult.RegisteredCountry.Names,
				IsInEuropeanUnion: ipResult.RegisteredCountry.IsInEuropeanUnion,
			},
			Postal: &geolize_pb.Postal{
				Code: ipResult.Postal.Code,
			},
			Traits: &geolize_pb.Traits{
				IsAnonymousProxy:    ipResult.Traits.IsAnonymousProxy,
				IsAnycast:           ipResult.Traits.IsAnycast,
				IsSatelliteProvider: ipResult.Traits.IsSatelliteProvider,
			},
		})
	}

	return ipInfos
}
