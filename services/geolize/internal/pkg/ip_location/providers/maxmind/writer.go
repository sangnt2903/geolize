package maxmind

import (
	"context"
	"encoding/json"
	"fmt"
	"geolize/services/geolize/internal/pkg/ip_location/model"
	jsonhelper "geolize/utilities/json_helper"
	"geolize/utilities/logging"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

type Writer struct {
	writer  *mmdbwriter.Tree
	history *versionHistoryManager
	logger  logging.Logger
	once    *sync.Once
}

func (w *Writer) Update(ctx context.Context, request *model.IPUpdateRequest) error {
	file, err := w.history.CreateFile(request)
	if err != nil {
		w.logger.Error(ctx, "Failed to create history file", logging.NewError(err)...)
		return err
	}

	output := filepath.Join(dbFolder, db)
	err = w.override(file, output)
	if err != nil {
		w.logger.Error(ctx, "Failed to override database", append(logging.NewError(err), logging.NewKeyVal("file", file))...)
		return err
	}

	err = w.history.SetVersion(file)
	if err != nil {
		w.logger.Error(ctx, "Failed to set version", logging.NewError(err)...)
		return err
	}

	return nil
}

func (w *Writer) Lookup(ctx context.Context, ip string) (*geoip2.City, error) {
	network, info := w.writer.Get(net.ParseIP(ip))
	fmt.Println(network.String(), jsonhelper.ToString(info))

	return nil, nil
}

func NewWriter(logger logging.Logger) (*Writer, error) {
	if logger == nil {
		panic("logger is nil")
	}

	logger.Debug(context.Background(), "IPGeolite Writer is being initializing....")

	if len(db) == 0 {
		logger.Fatal(context.Background(), "Database is not configured")
		return nil, fmt.Errorf("database is not configured")
	}

	writer, err := mmdbwriter.Load(filepath.Join(dbFolder, db), mmdbwriter.Options{})
	if err != nil {
		logger.Fatal(context.Background(), "Failed to load writer", logging.NewError(err)...)
		return nil, err
	}

	logger.Info(context.Background(), "IPGeolite Writer is initialized", logging.NewKeyVal("ip_usecase", db))

	w := &Writer{
		writer:  writer,
		logger:  logger,
		history: newVersionHistoryManager(),
		once:    &sync.Once{},
	}

	w.loadToLatest()

	return w, nil
}

func (w *Writer) loadToLatest() {
	w.once.Do(func() {
		w.logger.Debug(context.Background(), "Database is being updated...")
		files, err := w.history.GetAllFiles()
		if err != nil {
			w.logger.Fatal(context.Background(), "Failed to get history files", logging.NewError(err)...)
			return
		}

		defer func() {
			if err != nil {
				w.logger.Fatal(context.Background(), "Failed to load history files", logging.NewError(err)...)
			}

			w.logger.Debug(context.Background(), "Database is up to date")
		}()

		if len(files) == 0 {
			return
		}

		updatedFiles, err := w.history.GetUpdateFilesFromVersion()
		if err != nil {
			w.logger.Fatal(context.Background(), "Failed to get update history files", logging.NewError(err)...)
			return
		}

		if len(updatedFiles) == 0 {
			return
		}

		w.logger.Debug(context.Background(), "Updating database with files", logging.NewKeyVal("number_files", len(updatedFiles)))

		mergeFile, err := mergeFiles(updatedFiles)
		if err != nil {
			w.logger.Fatal(context.Background(), "Failed to merge history files", logging.NewError(err)...)
			return
		}

		defer func() {
			w.history.RemoveFile(mergeFile)
		}()

		w.logger.Info(context.Background(), "Updating database with file", logging.NewKeyVal("file", mergeFile))

		var config History
		mergedFileBytes, err := os.ReadFile(filepath.Join(dbHistories, mergeFile))
		if err != nil {
			w.logger.Fatal(context.Background(), "Failed to read history file", logging.NewError(err)...)
			return
		}
		if err = json.Unmarshal(mergedFileBytes, &config); err != nil {
			w.logger.Fatal(context.Background(), "Failed to parse history file", logging.NewError(err)...)
			return
		}

		for _, overrideIP := range config.Overrides {
			_, network, err := net.ParseCIDR(fmt.Sprintf("%s/32", overrideIP.IP))
			if err != nil {
				w.logger.Fatal(context.Background(), "Failed to parse ip", append(logging.NewError(err), logging.NewKeyVal("ip", overrideIP.IP))...)
				return
			}

			err = w.writer.InsertFunc(network, func(value mmdbtype.DataType) (mmdbtype.DataType, error) {
				v, ok := value.(mmdbtype.Map)
				if !ok {
					return nil, fmt.Errorf("expected Map, got %T", value)
				}

				applyOverride(v, overrideIP)

				return v, nil
			})
			if err != nil {
				w.logger.Fatal(context.Background(), "Failed to insert override", append(logging.NewError(err), logging.NewKeyVal("ip", overrideIP.IP))...)
				return
			}
		}

		output := filepath.Join(dbFolder, db)
		err = w.override(mergeFile, output)
		if err != nil {
			w.logger.Fatal(context.Background(), "Failed to override database", logging.NewError(err)...)
			return
		}

		// update version file
		err = w.history.SetVersion(config.Name)
		if err != nil {
			w.logger.Fatal(context.Background(), "Failed to update version file", logging.NewError(err)...)
			return
		}

		w.logger.Info(context.Background(), "Database has an update to version", logging.KeyVal{Key: "version", Val: config.Name})
	})
}

func (w *Writer) override(mergedFile string, output string) error {
	// Read and parse the override file
	overrideData, err := os.ReadFile(filepath.Join(dbHistories, mergedFile))
	if err != nil {
		log.Fatalf("Error reading override file: %v", err)
		return err
	}

	var config History
	if err = json.Unmarshal(overrideData, &config); err != nil {
		log.Fatalf("Error parsing override file: %v", err)
		return err
	}

	// Process each override
	for _, override := range config.Overrides {
		_, network, err := net.ParseCIDR(fmt.Sprintf("%s/32", override.IP))
		if err != nil {
			log.Fatalf("Error parsing IP: %v", err)
			return err
		}

		err = w.writer.InsertFunc(network, func(value mmdbtype.DataType) (mmdbtype.DataType, error) {
			v, ok := value.(mmdbtype.Map)
			if !ok {
				return nil, fmt.Errorf("expected Map, got %T", value)
			}

			applyOverride(v, override)

			return v, nil
		})
		if err != nil {
			log.Fatalf("Error inserting override: %v", err)
			return err
		}
	}

	tmpOutput := fmt.Sprintf("%s__%d.tmp", output, time.Now().Unix())
	fh, err := os.Create(tmpOutput)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			w.logger.Error(context.Background(), "Failed to override database", logging.NewError(fmt.Errorf("%v", r))...)
			os.Remove(tmpOutput)
		}
		_ = fh.Close()
	}()

	_, err = w.writer.WriteTo(fh)
	if err != nil {
		log.Fatalf("Error writing to output file: %v", err)
		return err
	}

	// Close the file before moving
	if err = fh.Close(); err != nil {
		w.logger.Error(context.Background(), "Failed to close temporary file", logging.NewError(err)...)
		return err
	}

	// Move the temporary file to the final destination
	if err = os.Rename(tmpOutput, output); err != nil {
		w.logger.Error(context.Background(), "Failed to move temporary file to output",
			logging.KeyVal{Key: "tmp", Val: tmpOutput},
			logging.KeyVal{Key: "output", Val: output},
			logging.KeyVal{Key: "error", Val: err.Error()})
		return err
	}

	return nil
}

type History struct {
	ID        string                   `json:"id"`
	Name      string                   `json:"name"`
	Overrides []*model.IPUpdateRequest `json:"overrides"`
}

func mergeFiles(files []string) (mergeFile string, err error) {
	var (
		mergedOverrides = make([]*model.IPUpdateRequest, 0)
	)

	// Iterate over the files and merge them into a single file
	for i, file := range files {

		// Read the file
		data, err := os.ReadFile(file)
		if err != nil {
			return mergeFile, fmt.Errorf("error reading file %s: %v", file, err)
		}

		var config History
		if err := json.Unmarshal(data, &config); err != nil {
			return mergeFile, fmt.Errorf("error parsing file %s: %v", file, err)
		}

		// Merge the overrides
		mergedOverrides = append(mergedOverrides, config.Overrides...)

		if i == len(files)-1 {
			err := func() error {
				mergeFile = fmt.Sprintf("merged_%s", filepath.Base(file))
				f, err := os.Create(filepath.Join(dbHistories, mergeFile))
				if err != nil {
					return fmt.Errorf("error creating merged file: %v", err)
				}
				defer f.Close()

				// Write the merged overrides to the file
				mergedData, err := json.Marshal(History{ID: config.ID, Name: filepath.Base(file), Overrides: mergedOverrides})
				if err != nil {
					return fmt.Errorf("error marshaling merged data: %v", err)
				}

				_, err = f.Write(mergedData)
				if err != nil {
					return fmt.Errorf("error writing merged data to file: %v", err)
				}

				return nil
			}()
			if err != nil {
				return mergeFile, err
			}

			break
		}
	}

	return mergeFile, nil
}

func applyOverride(original mmdbtype.Map, overrideIP *model.IPUpdateRequest) {
	// Apply overrides
	if overrideIP.Continent != nil {
		original["continent"] = mmdbtype.Map{
			"code": mmdbtype.String(overrideIP.Continent.Code),
			"names": func() mmdbtype.Map {
				names := make(mmdbtype.Map)
				for lang, name := range overrideIP.Continent.Names {
					names[mmdbtype.String(lang)] = mmdbtype.String(name)
				}
				return names
			}(),
		}
	}

	if overrideIP.Country != nil {
		original["country"] = mmdbtype.Map{
			"iso_code": mmdbtype.String(overrideIP.Country.ISOCode),
			"names": func() mmdbtype.Map {
				names := make(mmdbtype.Map)
				for lang, name := range overrideIP.Country.Names {
					names[mmdbtype.String(lang)] = mmdbtype.String(name)
				}
				return names
			}(),
			"is_in_european_union": mmdbtype.Bool(overrideIP.Country.IsInEuropeanUnion),
		}
	}

	if overrideIP.Subdivisions != nil {
		original["subdivisions"] = func() mmdbtype.Slice {
			subdivisions := make(mmdbtype.Slice, len(overrideIP.Subdivisions))
			for i, subdivision := range overrideIP.Subdivisions {
				subdivisions[i] = mmdbtype.Map{
					"iso_code": mmdbtype.String(subdivision.ISOCode),
					"names": func() mmdbtype.Map {
						names := make(mmdbtype.Map)
						for lang, name := range subdivision.Names {
							names[mmdbtype.String(lang)] = mmdbtype.String(name)
						}
						return names
					}(),
				}
			}
			return subdivisions
		}()
	}

	if overrideIP.City != nil {
		original["city"] = mmdbtype.Map{
			"names": func() mmdbtype.Map {
				names := make(mmdbtype.Map)
				for lang, name := range overrideIP.City.Names {
					names[mmdbtype.String(lang)] = mmdbtype.String(name)
				}
				return names
			}(),
		}
	}
	if overrideIP.Location != nil {
		original["location"] = mmdbtype.Map{
			"latitude":        mmdbtype.Float64(overrideIP.Location.Latitude),
			"longitude":       mmdbtype.Float64(overrideIP.Location.Longitude),
			"accuracy_radius": mmdbtype.Uint16(overrideIP.Location.AccuracyRadius),
			"time_zone":       mmdbtype.String(overrideIP.Location.TimeZone),
		}
	}

	if overrideIP.Postal != nil {
		original["postal"] = mmdbtype.Map{
			"code": mmdbtype.String(overrideIP.Postal.Code),
		}
	}

	if overrideIP.RepresentedCountry != nil {
		original["represented_country"] = mmdbtype.Map{
			"iso_code": mmdbtype.String(overrideIP.RepresentedCountry.ISOCode),
			"names": func() mmdbtype.Map {
				names := make(mmdbtype.Map)
				for lang, name := range overrideIP.RepresentedCountry.Names {
					names[mmdbtype.String(lang)] = mmdbtype.String(name)
				}
				return names
			}(),
			"is_in_european_union": mmdbtype.Bool(overrideIP.RepresentedCountry.IsInEuropeanUnion),
		}
	}

	if overrideIP.RegisteredCountry != nil {
		original["registered_country"] = mmdbtype.Map{
			"iso_code": mmdbtype.String(overrideIP.RegisteredCountry.ISOCode),
			"names": func() mmdbtype.Map {
				names := make(mmdbtype.Map)
				for lang, name := range overrideIP.RegisteredCountry.Names {
					names[mmdbtype.String(lang)] = mmdbtype.String(name)
				}
				return names
			}(),
			"is_in_european_union": mmdbtype.Bool(overrideIP.RegisteredCountry.IsInEuropeanUnion),
		}
	}

	if overrideIP.Traits != nil {
		original["traits"] = mmdbtype.Map{
			"is_anonymous_proxy":    mmdbtype.Bool(overrideIP.Traits.IsAnonymousProxy),
			"is_satellite_provider": mmdbtype.Bool(overrideIP.Traits.IsSatelliteProvider),
			"is_anycast":            mmdbtype.Bool(overrideIP.Traits.IsAnycast),
		}
	}
}
