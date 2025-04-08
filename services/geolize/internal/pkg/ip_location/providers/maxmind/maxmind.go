package maxmind

import (
	"context"
	"fmt"
	"geolize/services/geolize/internal/pkg/ip_location/model"
	"geolize/utilities/logging"
)

type Maxmind struct {
	logger logging.Logger
	reader *Reader
	writer *Writer
}

func (m *Maxmind) Lookup(ctx context.Context, request *model.IPLookupRequest) ([]*model.IPResult, error) {
	var result []*model.IPResult
	for _, ip := range request.IPs {
		record, err := m.reader.Lookup(ip)
		if err != nil {
			m.logger.Error(ctx, "reader.Lookup", logging.NewError(err)...)
			return nil, err
		}

		result = append(result, &model.IPResult{
			IP:        ip,
			DBVersion: m.reader.Version(),
			Country: &model.Country{
				ISOCode:           record.Country.IsoCode,
				Names:             record.Country.Names,
				IsInEuropeanUnion: record.Country.IsInEuropeanUnion,
			},
			City: &model.City{
				Names: record.City.Names,
			},
			Location: &model.Location{
				Latitude:       record.Location.Latitude,
				Longitude:      record.Location.Longitude,
				AccuracyRadius: record.Location.AccuracyRadius,
				TimeZone:       record.Location.TimeZone,
			},
			Postal: &model.Postal{
				Code: record.Postal.Code,
			},
			Continent: &model.Continent{
				Code:  record.Continent.Code,
				Names: record.Continent.Names,
			},
			Subdivisions: func() []*model.Subdivision {
				var subdivisions []*model.Subdivision
				for _, subdivision := range record.Subdivisions {
					subdivisions = append(subdivisions, &model.Subdivision{
						ISOCode: subdivision.IsoCode,
						Names:   subdivision.Names,
					})
				}
				return subdivisions
			}(),
			RepresentedCountry: &model.RepresentedCountry{
				ISOCode:           record.RepresentedCountry.IsoCode,
				Names:             record.RepresentedCountry.Names,
				Type:              record.RepresentedCountry.Type,
				IsInEuropeanUnion: record.RepresentedCountry.IsInEuropeanUnion,
			},
			RegisteredCountry: &model.RegisteredCountry{
				ISOCode:           record.RegisteredCountry.IsoCode,
				Names:             record.RegisteredCountry.Names,
				IsInEuropeanUnion: record.RegisteredCountry.IsInEuropeanUnion,
			},
			Traits: &model.Traits{
				IsAnonymousProxy:    record.Traits.IsAnonymousProxy,
				IsSatelliteProvider: record.Traits.IsSatelliteProvider,
				IsAnycast:           record.Traits.IsAnycast,
			},
		})
	}
	return result, nil
}

func (m *Maxmind) Update(ctx context.Context, request *model.IPUpdateRequest) error {
	if m.writer == nil {
		return fmt.Errorf("writer is not ready yet")
	}
	err := m.writer.Update(ctx, request)
	if err != nil {
		m.logger.Error(ctx, "writer.Update", logging.NewError(err)...)
		return err
	}
	return nil
}

func New(logger logging.Logger) *Maxmind {
	m := &Maxmind{
		logger: logger,
	}

	reader, err := NewReader(logger)
	if err != nil {
		panic(err)
	}

	m.reader = reader

	go func() {
		writer, err := NewWriter(logger)
		if err != nil {
			panic(err)
		}

		m.writer = writer
	}()

	logger.Info(context.Background(), "Maxmind geolocation provider initialized")

	return m
}
