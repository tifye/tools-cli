package pkg

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
)

var (
	ProductCatalogV2 string = "https://hqvrobotics.azure-api.net/productcatalog/v2"
	ProductCatalogQA string = "https://hqvrobotics.azure-api.net/productcatalog/qa"
)

type ProductCatalogService struct {
	base   string
	client *http.Client
	logger *log.Logger
}

func NewProductCatalogService(logger *log.Logger, api string, client *http.Client) *ProductCatalogService {
	return &ProductCatalogService{
		base:   api,
		client: client,
		logger: logger,
	}
}

func (pc *ProductCatalogService) DownloadDocument(ctx context.Context) (io.ReadCloser, error) {
	res, err := pc.client.Get(pc.base)
	if err != nil {
		return nil, fmt.Errorf("request to download product catalog failed with: %w", err)
	}

	if res.StatusCode > 299 {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("response failed with %s, failed to read body: %s", res.Status, err)
		}
		return nil, fmt.Errorf("response failed with %s, %s", res.Status, string(b))
	}

	return res.Body, nil
}
