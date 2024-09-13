package list

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/spf13/cobra"
)

type platformsOptions struct {
	brand string
}

func newPlatformsCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := platformsOptions{}

	cmd := &cobra.Command{
		Use:   "platforms",
		Short: "Output list of platforms from the product catalog for the given brand",
		Run: func(cmd *cobra.Command, args []string) {
			fpath := fmt.Sprintf("%s/product-catalog.json", cli.ConfigDir())
			file, err := os.Open(fpath)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println("Could not locate product catalog document. Did you forget to first download it? To download run `tools-cli catalog download`")
					os.Exit(1)
				} else {
					tCli.Log.Fatal(err)
				}
			}

			decoder := json.NewDecoder(file)
			var catalog pkg.ProductCatalog
			err = decoder.Decode(&catalog)
			if err != nil {
				tCli.Log.Fatal(err)
			}

			platforms := catalog.ListPlatformsForBrand(opts.brand)
			for _, p := range platforms {
				fmt.Println(p)
			}
		},
	}

	cmd.Flags().StringVarP(&opts.brand, "brand", "b", "", "Brand to list platforms for")

	return cmd
}
