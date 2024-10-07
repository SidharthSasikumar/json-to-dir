package cmd

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
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
    for key, value := range data {
        // Replace slashes in key with file path separators to handle nested directories.
        dirPath := filepath.Join(baseDir, strings.ReplaceAll(key, "/", string(os.PathSeparator)))
        if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
            return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
        }

        // Check if value is a nested map.
        if nestedMap, ok := value.(map[string]interface{}); ok {
            // Recursively create directories and files for nested maps.
            if err := createDirectoriesAndFiles(nestedMap, dirPath); err != nil {
                return err
            }
        } else {
            // Write the current value as a JSON file if it's not a map.
            filePath := filepath.Join(baseDir, fmt.Sprintf("%s.json", key))
            content, err := json.MarshalIndent(value, "", "    ")
            if err != nil {
                return fmt.Errorf("failed to marshal content for %s: %w", filePath, err)
            }
            if err := ioutil.WriteFile(filePath, content, 0644); err != nil {
                return fmt.Errorf("failed to write file %s: %w", filePath, err)
            }
        }
    }
    return nil
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
