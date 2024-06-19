package tifdefinition

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/internal/tif"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type testOptions struct {
	filepath string
	command  string
}

func newTestCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &testOptions{}

	cmd := &cobra.Command{
		Use:   "test",
		Short: "test command",
		Run: func(cmd *cobra.Command, args []string) {
			err := runTest(tCli.Log, *opts)
			if err != nil {
				tCli.Log.Fatal(err)
			}
		},
	}

	//const testFile string = `D:\Projects\_work\_pocs\tools-cli\internal\tif\testdata\40.x_Main-App-P21-Win_47.35-build-75.json`
	const testFile string = `D:\Projects\_work\_pocs\tools-cli\internal\tif\testdata\test.json`
	cmd.Flags().StringVarP(&opts.filepath, "def", "d", testFile, "Path to the tif definition file")
	cmd.MarkFlagFilename("def", "json")

	cmd.Flags().StringVarP(&opts.command, "command", "c", "", "Command to parse")
	cmd.MarkFlagRequired("command")

	return cmd
}

func runTest(logger *log.Logger, opts testOptions) error {
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

	start := time.Now()
	decoder := json.NewDecoder(file)
	var def tif.TifDefinition
	err = decoder.Decode(&def)
	if err != nil {
		return err
	}
	logger.Debugf("tif-definition decode took %dms", time.Since(start).Milliseconds())

	// Dear lord, help the person who has to debug this in the future, Amen.
	re, _ := regexp.Compile(`([A-Za-z0-9]+\.[A-Za-z0-9]+\((\s*[A-Za-z0-9]+\s*\:\s*[A-Za-z0-9]+\s*,?\s*)*\))`)
	if !re.MatchString(opts.command) {
		return fmt.Errorf("malformed command; expected format Family.Command(param: value), got: %s", opts.command)
	}

	opts.command = strings.ReplaceAll(opts.command, " ", "")
	opts.command = strings.Trim(opts.command, ")")
	family, rest, found := strings.Cut(opts.command, ".")
	if !found {
		return fmt.Errorf("failed to decode family, malformed command; expected format Family.Command(param: name), got: %s", opts.command)
	}
	command, argumentsStr, found := strings.Cut(rest, "(")
	if !found {
		return fmt.Errorf("failed to decode family, malformed command; expected format Family.Command(param: name), got: %s", opts.command)
	}

	methods := make(map[string]tif.MethodDefinition, len(def.Methods))
	for _, method := range def.Methods {
		key := fmt.Sprintf("%s.%s", method.Family, method.Command)
		methods[key] = method
	}

	logger.Debug("decomposed command", "family", family, "command", command, "arguments", argumentsStr)

	key := fmt.Sprintf("%s.%s", family, command)
	method, ok := methods[key]
	if !ok {
		fmt.Printf("Could not find method for %s", key)
		return nil
	}

	printUsage(method)

	arguments, err := parseInputArguments(argumentsStr)
	if err != nil {
		return err
	}
	if len(arguments) != len(method.InParams) {
		return fmt.Errorf("mismatch on number of arguments provided and parameters in the command; expected %d but got %d", len(method.InParams), len(arguments))
	}

	commandArgs := make([]CommandArgument, len(method.InParams))
	for i, argument := range arguments {
		inParamIdx := slices.IndexFunc(method.InParams, func(inParam tif.InputParameter) bool {
			return strings.EqualFold(inParam.Name, argument.Name)
		})
		if inParamIdx < 0 {
			return fmt.Errorf("method %s has no parameter named %s", method.Name(), argument.Name)
		}

		param := method.InParams[inParamIdx]

		logger.Debug("matched input argument", "argument", argument.Name, "value", argument.Value, "parameter", param.Name, "type", param.Type)

		tifType, err := tif.ParseType(param.Type, argument.Value)
		if err != nil {
			return err
		}

		commandArgs[i] = CommandArgument{
			Name:  param.Name,
			Value: tifType,
		}
	}

	for _, cmdArg := range commandArgs {
		fmt.Printf("%v\n", cmdArg)
	}

	return nil
}

type CommandArgument struct {
	Name  string
	Value any
}

type userInputArgument struct {
	Name  string
	Value string
}

func parseInputArguments(arguments string) ([]userInputArgument, error) {
	if arguments == "" {
		return []userInputArgument{}, nil
	}
	nameValPairs := strings.Split(arguments, ",")
	params := make([]userInputArgument, len(nameValPairs))
	for i, nameValPair := range nameValPairs {
		name, value, found := strings.Cut(nameValPair, ":")
		if !found || value == "" {
			return nil, fmt.Errorf("malformed parameter; expected (name: value) but got %s ", nameValPair)
		}

		params[i] = userInputArgument{
			Name:  strings.ToLower(name),
			Value: value,
		}
	}
	return params, nil
}

func printUsage(method tif.MethodDefinition) {
	fmt.Println("Method found")
	fmt.Printf("Family: %s\n", method.Family)
	fmt.Printf("Command: %s\n", method.Command)
	fmt.Printf("Description: %s\n", method.Description)
	fmt.Printf("ElementType: %s\n", method.ElementType)
	if len(method.InParams) > 0 {
		fmt.Printf("InParams: %v\n", method.InParams)
	}
	if len(method.OutParams) > 0 {
		fmt.Printf("OutParams: %v\n", method.OutParams)
	}
	fmt.Printf("Tags: %v\n", method.Tags)
	fmt.Printf("LoginLevels: %v\n", method.LoginLevels)

	args := make([]string, len(method.InParams))
	for i, arg := range method.InParams {
		args[i] = fmt.Sprintf("%s: %s", arg.Name, arg.Type)
	}
	argsStr := strings.Join(args, ", ")

	fmt.Printf("Usage: %s.%s(%s)\n\n", method.Family, method.Command, argsStr)
}
