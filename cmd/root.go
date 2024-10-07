package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "json-cli",
    Short: "A CLI tool to process JSON files and create directories and files",
    Long: `json-cli is a command line tool that reads a JSON file and creates
directories and files based on the structure of the JSON content.`,
}

// Execute runs the root command and all child commands.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
