/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	"github.com/taylormonacelli/deliverhalf/cmd/logging"

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
		data := Fetch()
		writeDataWrapper(data)
		writeBase64DataWrapper(data)
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

func writeBase64DataWrapper(data interface{}) error {
	wd, err := getWorkingDirectory()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
	}

	dataPath := filepath.Join(wd, "meta-b64.txt")
	deleteFile(dataPath)

	// encode the map to a JSON-encoded byte array
	b, err := json.Marshal(data)
	if err != nil {
		logging.Logger.Println("Error encoding map to JSON:", err)
	}

	base64Str := base64.StdEncoding.EncodeToString(b)
	if err := ioutil.WriteFile(dataPath, []byte(base64Str), 0o644); err != nil {
		return fmt.Errorf("Error writing base64-encoded JSON to file: %s", err)
	}

	logging.Logger.Printf("Successfully wrote base64-encoded JSON data to file %s", dataPath)

	return nil
}

func getMapAsString(data interface{}) string {
	jsBytes, _ := json.MarshalIndent(data, "", "    ")
	return string(jsBytes)
}

func writeData(dataPath string, data interface{}) error {
	jsonStr := getMapAsString(data)

	if err := ioutil.WriteFile(dataPath, []byte(jsonStr), 0o644); err != nil {
		return fmt.Errorf("Error writing JSON to file: %s", err)
	}

	logging.Logger.Printf("Successfully wrote JSON data to file %s", dataPath)
	return nil
}

func writeDataWrapper(data interface{}) error {
	wd, err := getWorkingDirectory()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
	}
	dataPath := filepath.Join(wd, "meta.json")

	deleteFile(dataPath)
	err = writeData(dataPath, data)
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

func deleteFile(filePath string) error {
	if !common.FileExists(filePath) {
		logging.Logger.Printf("%s doesn't exist, nothing to delete", filePath)
		return nil
	}

	logging.Logger.Printf("deleting %s", filePath)
	err := os.Remove(filePath)
	if err != nil {
		logging.Logger.Printf("%s couldn't be deleted", filePath)
		return err
	}
	return nil
}
