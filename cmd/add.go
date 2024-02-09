package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type ProjectInfo struct {
	path           string
	name           string
	runtime        string
	lang           string
	packageManager string
}

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "add project",
	Long:    "Add this project to gmp",
	Example: "gmp add",
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		pwd, pwdErr := os.Getwd()
		if pwdErr != nil {
			fmt.Println(pwdErr)
			panic("Invalid path")
		}

		projectInfo, projectInfoErr := getInfo(pwd)
		if projectInfoErr != nil {
			panic(projectInfoErr)
		}

		if projectInfoErr == nil && pwdErr == nil {
			fmt.Println("Adding project", projectInfo.name)
			fmt.Println(projectInfo)
		} else {
			fmt.Println("Project was Not added")
		}

	},
}

func getInfo(path string) (ProjectInfo, error) {

	name := getFolderName(path)

	runtime, err := getRuntime(path)
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		return ProjectInfo{}, err
	}
	switch runtime {
	case "node.js":
		packageManager, err := getPackageManager(path, runtime)
		if err != nil {
			fmt.Println(err)
			return ProjectInfo{}, err
		}
		return ProjectInfo{
				name:           name,
				path:           path,
				lang:           "typescript",
				runtime:        runtime,
				packageManager: packageManager},
			nil
	default:
		return ProjectInfo{
				name:           name,
				path:           path,
				lang:           runtime,
				runtime:        runtime,
				packageManager: runtime},
			nil
	}
}

func getRuntime(path string) (string, error) {
	runtimeFiles := map[string]string{
		"package.json": "node.js",
		"go.mod":       "go",
		"cargo.toml":   "rust",
		// TODO: Add more runtimes
	}

	return checkRuntimeFiles(path, runtimeFiles)
}

func checkRuntimeFiles(path string, files map[string]string) (string, error) {
	for filename, runtime := range files {
		if exists, err := checkFileExists(filepath.Join(path, filename)); err != nil {
			return "", err
		} else if exists {
			return runtime, nil
		}
	}

	// TODO: Handle the error
	return "", fmt.Errorf("unknown runtime")
}

func getPackageManager(path string, runtime string) (string, error) {
	switch runtime {
	case "node.js":
		lockFiles := map[string]string{
			"package-lock.json": "npm",
			"yarn.lock":         "yarn",
			"pnpm-lock.yml":     "pnpm",
			"bun.lock":          "bun",
			"wrangler.toml":     "cloudflare",
		}
		for filename, manager := range lockFiles {
			if exists, err := checkFileExists(filepath.Join(path, filename)); err != nil {
				return "", err
			} else if exists {
				return manager, nil
			}
		}
	default:
		err := fmt.Errorf("unknown package manager")
		return "", err
	}

	fmt.Println("Unknown package manager")
	return "", nil
}
