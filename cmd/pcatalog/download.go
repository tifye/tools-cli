package pcatalog

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/spf13/cobra"
)

type downloadOptions struct {
	output string
}

func newDownloadCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := downloadOptions{}

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Dowload product catalog",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.output == "" {
				opts.output = fmt.Sprintf("%s/product-catalog.json", cli.ConfigDir())
			}
			pcs := pkg.NewProductCatalogService(tCli.Log, pkg.ProductCatalogV2, tCli.Client)
			err := runDownload(cmd.Context(), tCli, pcs, opts)
			if err != nil {
				tCli.Log.Fatalf("failed to run download, got: %s", err)
			}
		},
	}

	cmd.Flags().StringVarP(&opts.output, "output", "o", "", "File to output document to")
	cmd.MarkFlagFilename("output", "json")

	return cmd
}

type productCatalogDownloader interface {
	DownloadDocument(context.Context) (io.ReadCloser, error)
}

func runDownload(
	ctx context.Context,
	tCli *cli.ToolsCli,
	pcs productCatalogDownloader,
	opts downloadOptions,
) error {
	rc, err := pcs.DownloadDocument(ctx)
	if err != nil {
		return fmt.Errorf("failed to download document, got: %w", err)
	}
	defer func() {
		err := rc.Close()
		if err != nil {
			tCli.Log.Warn("failed to close response body")
		}
	}()

	f, err := os.OpenFile(opts.output, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			tCli.Log.Warn("failed to close product catalog file")
		}
	}()

	_, err = io.Copy(f, rc)
	if err != nil {
		return fmt.Errorf("failed writing response to file, got: %s", err)
	}

	return nil
}
