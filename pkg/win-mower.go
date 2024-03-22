package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Tifufu/tools-cli/internal/winmower"
	"github.com/charmbracelet/log"
)

type WinMowerRegistry struct {
	CacheDir       string
	bundleRegistry *BundleRegistry
	client         http.Client
	logger         *log.Logger
}

type WinMower struct {
	Path string
}

func NewWinMowerRegistry(cacheDir string, bregsitry *BundleRegistry, logger *log.Logger) *WinMowerRegistry {
	return &WinMowerRegistry{
		bundleRegistry: bregsitry,
		CacheDir:       cacheDir,
		client:         *http.DefaultClient,
		logger:         logger,
	}
}

func (w *WinMowerRegistry) WithClient(client http.Client) {
	w.client = client
}

func (w *WinMowerRegistry) DownloadWinMower(platform Platform, ctx context.Context) (*WinMower, error) {
	wm, err := w.GetCachedWinMower(platform, ctx)
	if err != nil {
		return nil, err
	}
	if wm != nil {
		w.logger.Debug("Using cached winmower")
		return wm, nil
	}

	btypes, err := w.bundleRegistry.FetchBundleTypes(ctx)
	if err != nil {
		return nil, err
	}
	w.logger.Debugf("Found %d bundle types", len(btypes))

	btypes = FilterBundleTypes(btypes, platform)
	if len(btypes) == 0 {
		return nil, fmt.Errorf("no bundle types found for platform %s", platform)
	}
	w.logger.Debugf("Found %d bundle types for platform %s", len(btypes), platform)

	// Endpoint returns them sorted by date (i think)
	latestType := btypes[0]
	w.logger.Debugf("Latest bundle type: %s", latestType.Name)

	latestBuild, err := w.bundleRegistry.FetchLatestRelease(ctx, latestType.Name)
	if err != nil {
		return nil, err
	}
	w.logger.Debugf("Latest build: %s", latestBuild.BlobUrl)

	dir := filepath.Join(w.CacheDir, platform.String())
	req, err := http.NewRequestWithContext(ctx, "GET", latestBuild.BlobUrl, nil)
	w.logger.Debug(latestBuild.BlobUrl)
	if err != nil {
		return nil, err
	}
	w.logger.Debug("Downloading and unpacking winmower...")
	err = DownloadAndUnpack(req, &w.client, dir)
	if err != nil {
		return nil, err
	}

	index := filepath.Join(dir, "index.json")
	indexFile, err := os.OpenFile(index, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("error opening index.json: %v", err)
	}
	defer func() {
		closeErr := indexFile.Close()
		if closeErr != nil {
			w.logger.Error("error closing index.json", "err", err)
		}
	}()

	decoder := json.NewDecoder(indexFile)
	var manifest winmower.Manifest
	err = decoder.Decode(&manifest)
	if err != nil {
		return nil, fmt.Errorf("error decoding index.json: %v", err)
	}

	return &WinMower{
		Path: filepath.Join(dir, manifest.Metadata.UniqueDescriptiveName+".exe"),
	}, nil
}

func (w *WinMowerRegistry) GetCachedWinMower(platform Platform, ctx context.Context) (*WinMower, error) {
	var wmDir string
	err := filepath.WalkDir(w.CacheDir, func(path string, d fs.DirEntry, err error) error {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		if d.IsDir() && d.Name() == platform.String() {
			wmDir = path
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	if wmDir == "" {
		return nil, nil
	}

	index := filepath.Join(wmDir, "index.json")
	indexFile, err := os.OpenFile(index, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("error opening index.json: %v", err)
	}
	defer func() {
		closeErr := indexFile.Close()
		if closeErr != nil {
			w.logger.Error("error closing index.json", "err", err)
		}
	}()

	decoder := json.NewDecoder(indexFile)
	var manifest winmower.Manifest
	err = decoder.Decode(&manifest)
	if err != nil {
		return nil, fmt.Errorf("error decoding index.json: %v", err)
	}

	return &WinMower{
		Path: filepath.Join(wmDir, manifest.Metadata.UniqueDescriptiveName+".exe"),
	}, nil
}
