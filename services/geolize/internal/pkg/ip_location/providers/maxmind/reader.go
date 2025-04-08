package maxmind

import (
	"context"
	"fmt"
	"geolize/utilities/conf"
	"geolize/utilities/logging"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/oschwald/geoip2-golang"
)

type Reader struct {
	reader  *geoip2.Reader
	version string
	logger  logging.Logger
	watcher *fsnotify.Watcher
	history *versionHistoryManager
	done    chan struct{}
}

func newReader(logger logging.Logger, db string) (*Reader, error) {
	// Check if the MaxMind database file exists
	if _, err := os.Stat(filepath.Join(dbFolder, db)); os.IsNotExist(err) {
		err := fmt.Errorf("maxmind database file does not exist: %s", db)
		logger.Fatal(context.Background(), "MaxMind database file does not exist. Please download and install it.", logging.NewError(err)...)
		return nil, err
	}

	gReader, err := geoip2.Open(filepath.Join(dbFolder, db))
	if err != nil {
		return nil, err
	}

	vhm := newVersionHistoryManager()
	version, err := vhm.GetVersion()
	if err != nil {
		if os.IsNotExist(err) {
			if err = vhm.SetVersion(""); err != nil {
				return nil, err
			}
		} else {
			logger.Fatal(context.Background(), "Failed to get version file", logging.NewError(err)...)
			return nil, err
		}
	}

	reader := &Reader{
		reader:  gReader,
		version: version,
		history: vhm,
		logger:  logger,
	}

	// Start watching for version changes
	if err = reader.watch(); err != nil {
		reader.Close()
		return nil, err
	}

	return reader, nil
}

func (r *Reader) Close() {
	if r.watcher != nil {
		r.watcher.Close()
	}
	if r.done != nil {
		close(r.done)
	}
	if r.reader != nil {
		r.reader.Close()
	}
}

func (r *Reader) Lookup(ip string) (*geoip2.City, error) {
	ip = strings.TrimSpace(ip)
	return r.reader.City(net.ParseIP(ip))
}

func (r *Reader) reload() error {
	newVersion, err := r.history.GetVersion()
	if err != nil {
		return err
	}

	if newVersion == r.version {
		return nil // No version change
	}

	if r.version == newVersion {
		return nil
	}

	// Create a copy of the current reader
	oldReader := r.reader

	// Open new database file
	dbPath := filepath.Join(dbFolder, db)
	newReader, err := geoip2.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open new database: %w", err)
	}

	// Update reader and version
	r.reader = newReader
	r.version = newVersion

	// Close the old reader after the new one is successfully loaded
	if oldReader != nil {
		defer oldReader.Close()
	}

	return nil
}

func (r *Reader) Version() string {
	return r.version
}

func NewReader(logger logging.Logger) (*Reader, error) {
	if logger == nil {
		panic("Reader: logger is nil")
	}

	logger.Debug(context.Background(), "IPGeolite Reader is being initializing...")

	reader, err := newReader(logger, db)
	if err != nil {
		logger.Fatal(context.Background(), "Failed to create reader", logging.NewError(err)...)
		return nil, err
	}

	logger.Debug(context.Background(), "IPGeolite Reader is initialized")

	return reader, nil
}

func (r *Reader) watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	err = watcher.Add(versionFilePath)
	if err != nil {
		watcher.Close()
		return fmt.Errorf("failed to watch version file: %w", err)
	}

	r.watcher = watcher
	r.done = make(chan struct{})

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Remove) {
					r.logger.Error(context.Background(), "Version file removed. Recreating version file...")
					if err = r.history.SetVersion(r.Version()); err != nil {
						r.logger.Error(context.Background(), "Failed to create version file", logging.NewError(err)...)
					}
					r.logger.Info(context.Background(), "Version file recreated", logging.NewKeyVal("version", r.version))
				}
				if event.Has(fsnotify.Write) {
					r.logger.Info(context.Background(), "Version file changed. Reloading database...")
					if err = r.reload(); err != nil {
						r.logger.Error(context.Background(), "Failed to reload database", logging.NewError(err)...)
					}
					r.logger.Info(context.Background(), "Database is up to date with new version", logging.NewKeyVal("version", r.version))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				r.logger.Error(context.Background(), "Watcher error", logging.NewError(err)...)
			case <-r.done:
				return
			}
		}
	}()

	return nil
}

// downloadMaxMindDB downloads and extracts the MaxMind database
func downloadMaxMindDB(dbPath string, logger logging.Logger) error {
	// Get license key and download URL from configuration
	licenseKey, _ := conf.GetString("geolize", "maxmind_license_key", "")
	if licenseKey == "" {
		return fmt.Errorf("MaxMind license key not configured")
	}

	// Get download URL from config or use default
	downloadURLTemplate, _ := conf.GetString("geolize", "maxmind_download_url",
		"https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz")
	downloadURL := fmt.Sprintf(downloadURLTemplate, licenseKey)

	// Create a temporary directory for the download
	tempDir, err := os.MkdirTemp("", "maxmind")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download the database
	tarGzPath := filepath.Join(tempDir, "GeoLite2-City.tar.gz")

	logger.Info(context.Background(), "Downloading MaxMind database...")
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download database: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download database, status code: %d", resp.StatusCode)
	}

	out, err := os.Create(tarGzPath)
	if err != nil {
		return fmt.Errorf("failed to create download file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save downloaded file: %w", err)
	}

	// Extract the database
	logger.Info(context.Background(), "Extracting MaxMind database...")
	cmd := exec.Command("tar", "-xzf", tarGzPath, "-C", tempDir)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract database: %w", err)
	}

	// Find the extracted mmdb file
	var mmdbPath string
	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".mmdb") {
			mmdbPath = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to find extracted database: %w", err)
	}

	if mmdbPath == "" {
		return fmt.Errorf("could not find .mmdb file in extracted archive")
	}

	// Copy the database to the destination
	data, err := os.ReadFile(mmdbPath)
	if err != nil {
		return fmt.Errorf("failed to read extracted database: %w", err)
	}

	if err := os.WriteFile(dbPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write database to destination: %w", err)
	}

	logger.Info(context.Background(), "MaxMind database installed successfully")
	return nil
}
