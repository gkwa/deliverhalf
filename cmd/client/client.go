/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/taylormonacelli/deliverhalf/cmd"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	meta "github.com/taylormonacelli/deliverhalf/cmd/meta"
	sns "github.com/taylormonacelli/deliverhalf/cmd/sns"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client called")
		logger := common.SetupLogger()
		client(logger)
	},
}

func init() {
	cmd.RootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func client(logger *log.Logger) {
	data := meta.Fetch(logger)
	jsBytes, _ := json.MarshalIndent(data, "", "    ")
	jsonStr := string(jsBytes)
	sns.SendJsonStr(logger, jsonStr)
}
