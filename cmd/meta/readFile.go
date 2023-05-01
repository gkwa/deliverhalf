/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	common "github.com/taylormonacelli/deliverhalf/cmd/common"

	"github.com/spf13/cobra"
)

// readFileCmd represents the readFile command
var readFileCmd = &cobra.Command{
	Use:   "readFile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("readFile called")
	},
}

func init() {
	metaCmd.AddCommand(readFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readFileCmd.PersistentFlags().String("file", "", "A help for file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ParseJsonFromFile(logger *log.Logger, filename string) map[string]interface{} {
	if !common.FileExists(logger, filename) {
		logger.Fatalf("Can't find file %s, qutting", filename)
		os.Exit(1)
	}

	// read the JSON file into a byte slice
	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// create a map to hold the decoded JSON data
	data := make(map[string]interface{})

	// unmarshal the JSON data into the map
	err = json.Unmarshal(jsonBlob, &data)
	if err != nil {
		panic(err)
	}
	return data
}
