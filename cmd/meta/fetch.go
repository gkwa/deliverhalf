/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	common "github.com/taylormonacelli/deliverhalf/cmd/common"
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
		logger := common.SetupLogger()
		data := fetch(logger)
		fmt.Println(getMapAsString(logger, data))
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

func mergeData(logger *log.Logger, data []byte) map[string]interface{} {
	parsedData, err := parseData(data)
	if err != nil {
		logger.Fatalf("Error parsing JSON data:%s", err)
		panic(err)
	}

	// add epochtime timestamp blob
	newData := addEpochTimestamp(parsedData)
	return newData
}

func toJsonStr(logger *log.Logger, data map[string]interface{}) string {
	// Convert the map to a flat JSON string
	jsonStr, err := json.Marshal(data)
	if err != nil {
		logger.Println("Error parsing JSON data:", err)
		panic(err)
	}
	logger.Printf("json: %s", jsonStr)
	return string(jsonStr)
}

func toJsonPrettyStr(logger *log.Logger, data map[string]interface{}) string {
	// Convert the map to a pretty JSON string
	jsonStrPretty, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		logger.Println("Error marshaling data:", err)
		panic(err)
	}
	logger.Printf("json: %s", jsonStrPretty)
	return string(jsonStrPretty)
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

func fetch(logger *log.Logger) map[string]interface{} {
	body, err := fetchData()
	if err != nil {
		logger.Fatalf("Error fetching data: %s", err)
	}

	mergedData := mergeData(logger, body)
	return mergedData
}
