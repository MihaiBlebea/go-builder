package cmd

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

type Data struct {
	PackageName  string
	ModelName    string
	RepoName     string
	ModelVarName string
	RepoVarName  string
	Fields       []ModelField
	ExtraFields  []ModelField
}

//go:embed templates/*
var templates embed.FS

var modelName string

func init() {
	rootCmd.AddCommand(modelCmd)

	modelCmd.Flags().StringVarP(&modelName, "name", "n", "", "Model name")
	modelCmd.Flags().BoolVarP(&mute, "mute", "m", false, "Mute command output")
}

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generate a model & repository.",
	Long:  "Generate a model & repository.",
	Args: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			if !strings.Contains(arg, ":") {
				return fmt.Errorf("invalid argument format for %s. ex: user_name:string", arg)
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateCurrentFolder()
		if err != nil {
			return err
		}

		fields, err := parseArgs(args)
		if err != nil {
			return err
		}

		if modelName == "" {
			return errors.New("model name must be supplied")
		}

		if strings.Title(modelName) == modelName {
			return errors.New("model name must be lowercase")
		}

		if strings.Contains(modelName, " ") {
			return errors.New("model name should not contain empty spaces")
		}

		data := Data{
			PackageName:  strings.ToLower(modelName),
			ModelName:    strings.Title(modelName),
			RepoName:     fmt.Sprintf("%sRepo", strings.Title(modelName)),
			ModelVarName: strings.ToLower(modelName),
			RepoVarName:  fmt.Sprintf("%sRepo", strings.ToLower(modelName)),
			Fields:       fields,
			ExtraFields:  addExtraFields(fields),
		}

		if !mute {
			// Print output to confirm
			fmt.Printf("1. Create model %s\n", data.ModelName)
			fmt.Printf("2. Create repo %s\n", data.RepoName)
			fmt.Println("")

			// Print the model fields
			headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
			columnFmt := color.New(color.FgYellow).SprintfFunc()

			tbl := table.New("#", "Name", "Type", "JSON")
			tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			for i, f := range data.ExtraFields {
				tbl.AddRow(i+1, f.Title, f.FieldType, f.DBTitle)
			}

			tbl.Print()
			fmt.Println("")

			confirm, err := askConfirm()
			if err != nil {
				return err
			}

			if !confirm {
				fmt.Println("Terminating...")

				return nil
			}
		}

		if err := createOrSkipFolder(fmt.Sprintf("./%s", modelName)); err != nil {
			return err
		}

		err = createFile("model", fmt.Sprintf("./%s/%s.go", modelName, modelName), data)
		if err != nil {
			return err
		}

		err = createFile("repo", fmt.Sprintf("./%s/repo.go", modelName), data)
		if err != nil {
			return err
		}

		return nil
	},
}

func createOrSkipFolder(folderPath string) error {
	if _, err := os.Stat(folderPath); !os.IsNotExist(err) {
		return nil
	}

	err := os.Mkdir(folderPath, 0755)
	if err != nil {
		return err
	}

	return nil
}

func addExtraFields(fields []ModelField) []ModelField {
	var extra = []string{"id:int", "created:time.Time", "updated:time.Time"}
	result, _ := parseArgs(extra)

	return append(fields, result...)
}

func validateCurrentFolder() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	if parts := strings.Split(dir, "/"); len(parts) > 0 {
		if strings.Contains(parts[len(parts)-1], "go") {
			return nil
		}
	}

	return errors.New("run this command only in a go project folder. ex: go-casino")
}
