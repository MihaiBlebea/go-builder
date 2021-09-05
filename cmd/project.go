package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var projectName string

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name")
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Generate a new project from template.",
	Long:  "Generate a new project from template.",
	Args: func(cmd *cobra.Command, args []string) error {

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateProjectName(projectName)
		if err != nil {
			return err
		}

		folderPath := fmt.Sprintf("./%s", projectName)

		res := exec.Command(
			"git",
			"clone",
			"https://github.com/MihaiBlebea/go-template",
			folderPath,
		)
		err = res.Run()
		if err != nil {
			return err
		}

		// Go over each file and update go-template into go-project name
		if err := filepath.Walk(folderPath, visit); err != nil {
			return err
		}

		return nil
	},
}

func validateProjectName(name string) error {
	if name == "" {
		return errors.New("project name cannot be empty")
	}

	if !strings.Contains(name, "go-") {
		return errors.New("project name must start with go-. ex: go-casino ")
	}

	if _, err := os.Stat(fmt.Sprintf("./%s", name)); !os.IsNotExist(err) {
		return fmt.Errorf("folder %s already exists", name)
	}

	return nil
}

func visit(path string, fi os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	matched, err := matchFilePatterns(fi.Name(), []string{"*.go", "*.mod"})
	if err != nil {
		return err
	}

	if matched {
		read, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Print updated path
		printUpdatedFile(path, projectName)

		newContents := strings.Replace(string(read), "go-template", projectName, -1)

		err = ioutil.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func matchFilePatterns(fileName string, patterns []string) (bool, error) {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, fileName)
		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}

func printUpdatedFile(path, projectName string) {
	fmt.Printf("updated path: %s - change go-template with %s\n", path, projectName)
}
