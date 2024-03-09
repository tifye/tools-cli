package amprod

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type validateTestIndexOptions struct {
	directory string
}

type testFileInfo struct {
	Name     string   `json:"name"`
	Filename string   `json:"file"`
	Methods  []string `json:"methods"`
}

type testsManifest struct {
	Name     string         `json:"name"`
	Sequence []testFileInfo `json:"sequence"`
}

func newValidateTestIndexCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &validateTestIndexOptions{}

	cmd := &cobra.Command{
		Use:   "validate-test-index",
		Short: "Validate test index",
		Run: func(cmd *cobra.Command, args []string) {
			defer timer("walkDir", tCli.Log)()

			err := filepath.WalkDir(opts.directory, walkDirFunc(tCli.Log))
			if err != nil {
				tCli.Log.Error("Error validating test index.json", "err", err)
				return
			}
		},
	}

	cmd.Flags().StringVarP(&opts.directory, "directory", "d", "", "Directory in which to validate test index.json and all subdirectories.")
	cmd.MarkFlagDirname("directory")
	cmd.MarkFlagRequired("directory")

	return cmd
}

func walkDirFunc(logger *log.Logger) fs.WalkDirFunc {
	return func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			logger.Debug("Walking", "path", path)
			return err
		}

		if entry.Name() != "index.json" {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		defer logger.Print("")

		decoder := json.NewDecoder(file)
		var manifest testsManifest
		err = decoder.Decode(&manifest)
		if err != nil {
			logger.Error("Error decoding index.json", "err", err, "path", path)
			return nil
		}

		err = validateTestManifest(manifest, filepath.Dir(path), logger)
		return err
	}
}

func validateTestManifest(manifest testsManifest, dir string, logger *log.Logger) error {
	logger.Info("Validating", "index", manifest.Name, "dir", dir)

	for _, test := range manifest.Sequence {
		if test.Filename == "" {
			logger.Warn("Test with missing filename", "test", test.Name, "methods", test.Methods)
			continue
		}

		err := func(test testFileInfo) error {
			fpath := filepath.Join(dir, test.Filename)
			testFile, err := os.Open(fpath)
			switch err.(type) {
			case *os.PathError:
				logger.Warn("Test file not found", "test", test.Name, "file", test.Filename, "path", fpath)
				return nil
			case nil:
				// continue
			default:
				return err
			}
			defer testFile.Close()

			return validateTestFile(testFile, test, logger)
		}(test)

		if err != nil {
			return err
		}
	}

	return nil
}

func validateTestFile(content io.Reader, testInfo testFileInfo, logger *log.Logger) error {
	methodsFound := make(map[string]bool)
	for _, method := range testInfo.Methods {
		methodsFound[method] = false
	}

	scanner := bufio.NewScanner(content)
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
			logger.Warn("Missing method", "test", testInfo.Name, "filename", testInfo.Filename, "method", method)
		}
	}

	return nil
}

func timer(name string, logger *log.Logger) func() {
	start := time.Now()
	return func() {
		logger.Debugf("%s: executed in %s", name, time.Since(start))
	}
}
