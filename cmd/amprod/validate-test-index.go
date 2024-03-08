package amprod

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type validateTestIndexOptions struct {
	indexPath string
}

type testFileInfo struct {
	Filename string   `json:"file"`
	Methods  []string `json:"methods"`
}

type testsManifest struct {
	Name     string         `json:"name"`
	Sequence []testFileInfo `json:"sequence"`
}

func newValidateTestIndexCommand(toolsCli *cli.ToolsCli) *cobra.Command {
	opts := &validateTestIndexOptions{}

	cmd := &cobra.Command{
		Use:   "validate-test-index",
		Short: "Validate test index",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runValidateTestIndex(toolsCli, opts); err != nil {
				log.Error("Error validating test index.json", "err", err)
				return
			}
		},
	}

	cmd.Flags().StringVarP(&opts.indexPath, "index", "i", "", "Path to the index.json file")
	cmd.MarkFlagFilename("index.json")
	cmd.MarkFlagRequired("index")

	return cmd
}

func runValidateTestIndex(toolsCli *cli.ToolsCli, opts *validateTestIndexOptions) error {
	file, err := os.Open(opts.indexPath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var manifest testsManifest
	err = decoder.Decode(&manifest)
	if err != nil {
		return fmt.Errorf("error decoding index.json: %w", err)
	}

	logger := log.NewWithOptions(os.Stderr, log.Options{
		Level:           log.InfoLevel,
		ReportTimestamp: false,
		ReportCaller:    false,
	})
	testsDir := filepath.Dir(opts.indexPath)
	for _, test := range manifest.Sequence {
		err := validateTestFile(testsDir, test.Filename, test, logger)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateTestFile(dir, filename string, testInfo testFileInfo, logger *log.Logger) error {
	logger = logger.WithPrefix(filename)

	testFilePath := filepath.Join(dir, filename)
	logger.Info("Validating", "methods", testInfo.Methods)

	testFile, err := os.Open(testFilePath)
	switch err.(type) {
	case *os.PathError:
		logger.Error("Test file not found", "file", testFilePath)
		return nil
	case nil:
		// continue
	default:
		return err
	}
	defer testFile.Close()

	methodsFound := make(map[string]bool)
	for _, method := range testInfo.Methods {
		methodsFound[method] = false
	}

	scanner := bufio.NewScanner(testFile)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		line, hadPrefix := bytes.CutPrefix(line, []byte("def"))
		if !hadPrefix {
			continue
		}

		line, hadSuffix := bytes.CutSuffix(line, []byte("():"))
		if !hadSuffix {
			continue
		}

		line = bytes.TrimSpace(line)
		methodsFound[string(line)] = true
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for method, found := range methodsFound {
		if !found {
			logger.Error("Missing", "method", method)
		}
	}

	return nil
}
