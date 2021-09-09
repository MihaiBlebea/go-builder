package cmd

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

type ModelField struct {
	Title      string
	FieldType  string
	DBTitle    string
	ParamTitle string
}

type Data struct {
	PackageName  string
	ModelName    string
	RepoName     string
	ModelVarName string
	RepoVarName  string
	Fields       []ModelField
	ExtraFields  []ModelField
}

var fieldTypes = []string{"string", "int", "bool", "time.Time"}

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

func createFile(tmplName, outputFile string, data Data) error {
	funcMap := template.FuncMap{
		"ToLower":          strings.ToLower,
		"ToTitle":          strings.Title,
		"ToUpper":          strings.ToUpper,
		"RenderFuncParams": RenderFuncParams,
	}

	tmpl, err := template.New(fmt.Sprintf("%s.tmpl", tmplName)).Funcs(funcMap).ParseFS(templates, fmt.Sprintf("templates/%s.tmpl", tmplName))
	if err != nil {
		return err
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}

func parseArgs(args []string) (fields []ModelField, _ error) {
	for _, arg := range args {
		params := strings.Split(arg, ":")

		if !contains(fieldTypes, params[1]) {
			return fields, fmt.Errorf("unsupported field type: %s", params[1])
		}

		fields = append(fields, ModelField{
			Title:      transformTitle(params[0]),
			FieldType:  params[1],
			DBTitle:    params[0],
			ParamTitle: transformParamTitle(params[0]),
		})
	}

	return fields, nil
}

func transformTitle(title string) (updated string) {
	params := strings.Split(strings.Title(strings.ReplaceAll(title, "_", " ")), " ")

	for _, param := range params {
		if len(param) < 3 {
			param = strings.ToUpper(param)
		}

		updated += param
	}

	return updated
}

func transformParamTitle(title string) (updated string) {
	params := strings.Split(strings.Title(strings.ReplaceAll(title, "_", " ")), " ")

	for index, param := range params {
		if index == 0 {
			updated += strings.ToLower(param)
			continue
		}

		if len(param) < 3 {
			param = strings.ToUpper(param)
		}

		updated += param
	}

	return updated
}

func RenderFuncParams(params []ModelField) (render string) {
	for index, param := range params {
		if index < len(params)-1 {
			render += fmt.Sprintf("%s %s, ", param.ParamTitle, param.FieldType)
			continue
		}

		render += fmt.Sprintf("%s %s", param.ParamTitle, param.FieldType)
	}

	return render
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
