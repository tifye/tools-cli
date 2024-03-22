package pkg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

type GSPacketRegistry struct {
	client   *http.Client
	cacheDir string
	baseUrl  string
	logger   *log.Logger
}

type GSPacketMetadata struct {
	Map        string
	TestBundle string
}

func NewGSPacketRegistry(cacheDir, baseUrl string, client *http.Client, logger *log.Logger) *GSPacketRegistry {
	return &GSPacketRegistry{
		cacheDir: cacheDir,
		baseUrl:  baseUrl,
		client:   client,
		logger:   logger,
	}
}

func (r *GSPacketRegistry) DownloadGSPacket(serialNumber uint, platform Platform, ctx context.Context) (*GSPacketMetadata, error) {
	gsp, err := r.GetGSPacketFromCache(serialNumber)
	if err != nil {
		return nil, err
	}
	if gsp != nil {
		r.logger.Debug("Using cached GSPacket", "serialNumber", serialNumber, "path", filepath.Dir(gsp.TestBundle))
		return gsp, nil
	}

	endpoint := fmt.Sprintf("%s/packet/%d/%s", r.baseUrl, serialNumber, platform)
	r.logger.Debug("Downloading GSPacket", "endpoint", endpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(r.cacheDir, fmt.Sprint(serialNumber))
	err = DownloadAndUnpack(req, r.client, dir)
	if err != nil {
		return nil, err
	}

	gsp, err = locateGSPPaths(dir, serialNumber)
	if err != nil {
		return nil, err
	}

	return gsp, nil
}

func (r *GSPacketRegistry) GetGSPacketFromCache(serialNumber uint) (*GSPacketMetadata, error) {
	dir := filepath.Join(r.cacheDir, fmt.Sprint(serialNumber))
	_, err := os.Stat(dir)
	if err == nil {
		return locateGSPPaths(dir, serialNumber)
	}
	if os.IsNotExist(err) {
		return nil, nil
	}
	return nil, err
}

func locateGSPPaths(dir string, serialNumber uint) (*GSPacketMetadata, error) {
	gspPaths := &GSPacketMetadata{}
	bundleSuffix := fmt.Sprintf("%d.zip", serialNumber)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		switch {
		case info.Name() == "map.json":
			gspPaths.Map = path
		case strings.HasSuffix(info.Name(), bundleSuffix):
			gspPaths.TestBundle = path
		default:
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if gspPaths.Map == "" || gspPaths.TestBundle == "" {
		return nil, errors.New("failed to locate GSP paths")
	}

	return gspPaths, nil
}
