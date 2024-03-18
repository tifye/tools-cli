package pkg

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/charmbracelet/log"
)

type SimulatorRegistry struct {
	cacheDir       string
	bundleRegistry *BundleRegistry
	logger         *log.Logger
	client         *http.Client
}

type Simulator struct {
	Path string
}

func NewSimulatorRegistry(cacheDir string, bundleRegistry *BundleRegistry, client *http.Client, logger *log.Logger) *SimulatorRegistry {
	return &SimulatorRegistry{
		cacheDir:       cacheDir,
		bundleRegistry: bundleRegistry,
		client:         client,
		logger:         logger,
	}
}

func (s *SimulatorRegistry) DownloadSimulator(ctx context.Context) (*Simulator, error) {
	sim, err := s.GetCachedSimulator(ctx)
	if err != nil {
		return nil, err
	}
	if sim != nil {
		s.logger.Debug("Using cached simulator")
		return sim, nil
	}

	s.logger.Debug("Fetching simulator...")
	latestBuild, err := s.bundleRegistry.FetchLatestRelease(ctx, "GardenSimulator")
	if err != nil {
		return nil, err
	}
	s.logger.Debug("Latest Simulator build", "url", latestBuild.BlobUrl)

	req, err := http.NewRequestWithContext(ctx, "GET", latestBuild.BlobUrl, nil)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Downloading and unpacking simulator...")
	err = DownloadAndUnpack(req, s.client, s.cacheDir)
	if err != nil {
		return nil, err
	}

	return s.GetCachedSimulator(ctx)
}

func (s *SimulatorRegistry) GetCachedSimulator(ctx context.Context) (*Simulator, error) {
	var exePath string
	err := filepath.Walk(s.cacheDir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Base(path) == "GardenSimulator.exe" {
			exePath = path
			return nil
		}
		return nil
	})

	if errors.Is(err, fs.ErrNotExist) || exePath == "" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &Simulator{
		Path: exePath,
	}, nil
}
