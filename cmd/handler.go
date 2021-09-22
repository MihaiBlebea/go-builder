package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var handlerName string
var packageName string

func init() {
	rootCmd.AddCommand(handlerCmd)

	handlerCmd.Flags().StringVarP(&handlerName, "name", "n", "", "Handler name")
	handlerCmd.Flags().StringVarP(&packageName, "package", "p", "", "Package name")
}

var handlerCmd = &cobra.Command{
	Use:   "handler",
	Short: "Generate a handler.",
	Long:  "Generate a handler.",
	Args: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			if !strings.Contains(arg, ":") {
				return fmt.Errorf("invalid argument format for %s. ex: user_name:string", arg)
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if handlerName == "" {
			return errors.New("handler name must be supplied")
		}

		if packageName == "" {
			return errors.New("package name must be supplied")
		}

		fileName := fmt.Sprintf("%s_handler.go", handlerName)

		reqArgs, err := parseArgs(args)
		if err != nil {
			return err
		}

		data := struct {
			HandlerName      string
			PackageName      string
			RequestArguments []ModelField
		}{
			PackageName:      strings.ToLower(packageName),
			HandlerName:      transformTitle(handlerName),
			RequestArguments: reqArgs,
		}

		if err := createFile("handler", fileName, &data); err != nil {
			return err
		}

		return nil
	},
}
