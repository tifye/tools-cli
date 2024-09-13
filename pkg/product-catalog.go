package pkg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

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

type ProductCatalogModel struct {
	Model      string `json:"model"`
	Type       int    `json:"type"`
	Variant    int    `json:"variant"`
	Brand      string `json:"brand"`
	Platform   string `json:"platform"`
	Generation string `json:"generation"`
}

type ProductCatalog struct {
	Models []ProductCatalogModel `json:"models"`
}

func (pc *ProductCatalog) ListPlatformsForBrand(brand string) []string {
	if pc == nil {
		panic("about to dereference nil pointer on ProductCatalogModel in call to ListPlatformsForBrand")
	}

	platforms := make(map[string]struct{}, 0)
	for _, model := range pc.Models {
		if strings.EqualFold(brand, model.Brand) {
			platforms[model.Platform] = struct{}{}
		}
	}

	a := make([]string, 0, len(platforms))
	for k, _ := range platforms {
		a = append(a, k)
	}

	return a
}
