/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/natefinch/lumberjack.v2"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fetch()
	},
}

func init() {
	metaCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func writeBase64Data(logger *log.Logger, dataPath string, data []byte) error {
	base64Str := base64.StdEncoding.EncodeToString(data)
	if err := ioutil.WriteFile(dataPath, []byte(base64Str), 0o644); err != nil {
		return fmt.Errorf("Error writing base64-encoded JSON to file: %s", err)
	}

	logger.Printf("Successfully wrote base64-encoded JSON data to file %s", dataPath)
	return nil
}

func setupLogger() *log.Logger {
	logFile := &lumberjack.Logger{
		Filename:   "fetchmeta.log",
		MaxSize:    1, // In megabytes
		MaxBackups: 0,
		MaxAge:     365, // In days
	}
	defer logFile.Close()
	logWriter := io.MultiWriter(logFile, os.Stderr)
	return log.New(logWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)
}

func fileExists(logger *log.Logger, filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func deleteFile(logger *log.Logger, filePath string) error {
	if !fileExists(logger, filePath) {
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

func parseData(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("Error parsing JSON data: %s", err)
	}
	return result, nil
}

func addEpochTimestamp(data map[string]interface{}) map[string]interface{} {
	timestamp := time.Now().Unix()
	newData := map[string]interface{}{
		"epochtime": timestamp,
	}
	for k, v := range newData {
		data[k] = v
	}
	return data
}

func writeData(logger *log.Logger, dataPath string, data interface{}) error {
	jsonStr, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("Error pretty-printing JSON: %s", err)
	}

	if err := ioutil.WriteFile(dataPath, jsonStr, 0o644); err != nil {
		return fmt.Errorf("Error writing JSON to file: %s", err)
	}

	logger.Printf("Successfully wrote JSON data to file %s", dataPath)
	return nil
}

func fetchData() ([]byte, error) {
	url := "http://169.254.169.254/latest/dynamic/instance-identity/document"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating HTTP request: %s", err)
	}

	client := &http.Client{
		Timeout: time.Second * 2,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making HTTP request: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %s", err)
	}

	return body, nil
}

func fetch() {
	logger := setupLogger()

	logger.Println("fetch called")

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return
	}

	dataPath := filepath.Join(wd, "meta.json")
	dataPath2 := filepath.Join(wd, "meta-b64.txt")

	deleteFile(logger, dataPath)
	deleteFile(logger, dataPath2)

	body, err := fetchData()
	if err != nil {
		logger.Fatalf("Error fetching data: %s", err)
	}

	parsedData, err := parseData(body)
	if err != nil {
		logger.Fatalf("Error parsing JSON data:%s", err)
		panic(err)
	}

	// add epochtime timestamp blob
	newData := addEpochTimestamp(parsedData)

	// Convert the map to a flat JSON string
	jsonStr, err := json.Marshal(newData)
	if err != nil {
		logger.Println("Error parsing JSON data:", err)
		panic(err)
	}
	logger.Printf("json: %s", jsonStr)

	// Convert the map to a pretty JSON string
	jsonStrPretty, _ := json.MarshalIndent(newData, "", "    ")
	logger.Printf("json: %s", jsonStrPretty)

	writeData(logger, dataPath, newData)
	writeBase64Data(logger, dataPath2, jsonStrPretty)
}
