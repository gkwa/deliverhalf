/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	common "github.com/taylormonacelli/deliverhalf/cmd/common"

	"github.com/spf13/cobra"
)

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := common.SetupLogger()

		data := fetch(logger)
		writeDataWrapper(logger, data)
		writeBase64DataWrapper(logger, data)
	},
}

func init() {
	metaCmd.AddCommand(writeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// writeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// writeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func writeBase64DataWrapper(logger *log.Logger, data interface{}) error {
	wd, err := getWorkingDirectory()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
	}

	dataPath := filepath.Join(wd, "meta-b64.txt")
	deleteFile(logger, dataPath)

	// encode the map to a JSON-encoded byte array
	b, err := json.Marshal(data)
	if err != nil {
		logger.Println("Error encoding map to JSON:", err)
	}

	base64Str := base64.StdEncoding.EncodeToString(b)
	if err := ioutil.WriteFile(dataPath, []byte(base64Str), 0o644); err != nil {
		return fmt.Errorf("Error writing base64-encoded JSON to file: %s", err)
	}

	logger.Printf("Successfully wrote base64-encoded JSON data to file %s", dataPath)

	return nil
}

func getMapAsString(logger *log.Logger, data interface{}) string {
	jsBytes, _ := json.MarshalIndent(data, "", "    ")
	return string(jsBytes)
}

func writeData(logger *log.Logger, dataPath string, data interface{}) error {
	jsonStr := getMapAsString(logger, data)

	if err := ioutil.WriteFile(dataPath, []byte(jsonStr), 0o644); err != nil {
		return fmt.Errorf("Error writing JSON to file: %s", err)
	}

	logger.Printf("Successfully wrote JSON data to file %s", dataPath)
	return nil
}

func writeDataWrapper(logger *log.Logger, data interface{}) error {
	wd, err := getWorkingDirectory()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
	}
	dataPath := filepath.Join(wd, "meta.json")

	deleteFile(logger, dataPath)
	err = writeData(logger, dataPath, data)
	if err != nil {
		fmt.Println("Error writing data to file:", err)
	}
	return nil
}

func getWorkingDirectory() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return wd, nil
}

func deleteFile(logger *log.Logger, filePath string) error {
	if !common.FileExists(logger, filePath) {
		logger.Printf("%s doesn't exist, nothing to delete", filePath)
		return nil
	}

	logger.Printf("deleting %s", filePath)
	err := os.Remove(filePath)
	if err != nil {
		logger.Printf("%s couldn't be deleted", filePath)
		return err
	}
	return nil
}
