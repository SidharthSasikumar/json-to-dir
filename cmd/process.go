package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process [input JSON file] [output directory]",
	Short: "Process a JSON file and create directories and files based on its structure",
	Args:  cobra.ExactArgs(2),
	RunE:  runProcess,
}

func init() {
	rootCmd.AddCommand(processCmd)
}

func runProcess(cmd *cobra.Command, args []string) error {
	inputFilePath := args[0]
	outputDir := args[1]

	// Load JSON data from the input file.
	data, err := loadJSON(inputFilePath)
	if err != nil {
		return fmt.Errorf("error loading JSON: %w", err)
	}

	// Create the directory structure and files.
	if err := createDirectoriesAndFiles(data, outputDir); err != nil {
		return fmt.Errorf("error creating directories and files: %w", err)
	}

	fmt.Println("Directories and files created successfully.")
	return nil
}

// createDirectoriesAndFiles processes the JSON and creates the required directories and files.
func createDirectoriesAndFiles(data map[string]interface{}, baseDir string) error {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return err
	}
	root, err := os.OpenRoot(baseDir)
	if err != nil {
		return err
	}
	defer root.Close()
	return writeTree(root, data, ".")
}

// Root operations confine filesystem access even when an existing symlink points outside.
func writeTree(root *os.Root, data map[string]interface{}, prefix string) error {
	for key, value := range data {
		if !filepath.IsLocal(key) || strings.Contains(key, "\\") {
			return fmt.Errorf("unsafe JSON key %q", key)
		}
		name := filepath.Join(prefix, key)
		if nested, ok := value.(map[string]interface{}); ok && !isSimpleMap(nested) {
			if err := root.MkdirAll(name, 0755); err != nil {
				return err
			}
			if err := writeTree(root, nested, name); err != nil {
				return err
			}
			continue
		}
		content, err := json.MarshalIndent(value, "", "    ")
		if err != nil {
			return err
		}
		if err := root.MkdirAll(filepath.Dir(name), 0755); err != nil {
			return err
		}
		if err := root.WriteFile(name+".json", content, 0644); err != nil {
			return err
		}
	}
	return nil
}

// isSimpleMap checks if a map is a simple map (i.e., it contains no nested maps).
func isSimpleMap(data map[string]interface{}) bool {
	for _, value := range data {
		if _, ok := value.(map[string]interface{}); ok {
			return false
		}
	}
	return true
}

// loadJSON loads the JSON content from a file.
func loadJSON(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var data map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON file: %w", err)
	}

	return data, nil
}
