package cmd

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

// Used by the mute flag in commands
var mute bool

var fieldTypes = []string{"string", "int", "bool", "time.Time"}

type ModelField struct {
	Title      string
	FieldType  string
	DBTitle    string
	ParamTitle string
}

func askConfirm() (bool, error) {
	fmt.Println("Please confirm [yes/no]:")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, err
	}

	okResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokResponses := []string{"n", "N", "no", "No", "NO"}

	if contains(okResponses, response) {
		return true, nil
	} else if contains(nokResponses, response) {
		return false, nil
	} else {
		return false, nil
	}
}

func contains(list []string, needle string) bool {
	for _, el := range list {
		if el == needle {
			return true
		}
	}

	return false
}

func createFile(tmplName, outputFile string, data interface{}) error {
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
