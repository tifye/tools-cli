package tifdefinition

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/internal/tif"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type listOptions struct {
	filepath     string
	familyFilter string
}

func newListCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tif definitions",
		Run: func(cmd *cobra.Command, args []string) {
			err := runList(tCli.Log, *opts)
			if err != nil {
				tCli.Log.Fatal("failed to run list command", "error", err)
			}
		},
	}

	cmd.Flags().StringVarP(&opts.filepath, "def", "d", "", "Path to the tif definition file")
	cmd.MarkFlagRequired("file")
	cmd.MarkFlagFilename("file", "json")

	cmd.Flags().StringVarP(&opts.familyFilter, "family", "f", "", "Family to filter output by")

	return cmd
}

func runList(logger *log.Logger, opts listOptions) error {
	if filepath.IsLocal(opts.filepath) {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		opts.filepath = filepath.Join(wd, opts.filepath)
		logger.Debug("local path provided, combined into absolute", "result", opts.filepath)
	}

	file, err := os.Open(opts.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var def tif.TifDefinition
	err = decoder.Decode(&def)
	if err != nil {
		return err
	}

	methodIdmap := make(map[string]tif.MethodDefinition, len(def.Methods))
	for _, method := range def.Methods {
		key := fmt.Sprintf("%s.%s", method.Family, method.Command)
		methodIdmap[key] = method
	}

	for _, attr := range def.AttributesV2 {
		if opts.familyFilter != "" && attr.Family != opts.familyFilter {
			continue
		}

		fmt.Printf("%s.%s, %s\n", attr.Family, attr.Name, attr.Description)
	}
	return nil
}
