/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := common.SetupLogger()
		setupConfig(logger)
	},
}

func init() {
	configCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setupConfig(logger *log.Logger) {
	addConfigPaths()
	setConfigNameAndType()
	SetDefaultValues()

	if err := viper.ReadInConfig(); err != nil {
		handleConfigReadError(logger, err)
	}
}

func addConfigPaths() {
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
}

func setConfigNameAndType() {
	viper.SetConfigName(".deliverhalf")
	viper.SetConfigType("yaml")
}

func SetDefaultValues() {
	viper.SetDefault("SNS", map[string]string{
		"topic-arn": "arn:aws:sns:us-west-2:123456789012:example-topic",
		"region":    "us-west-2",
	})
	viper.SetDefault("SQS", map[string]string{
		"region":    "us-west-2",
		"queue-arn": "arn:aws:sqs:us-west-2:193048895737",
		"queue-url": "https://sqs.us-west-2.amazonaws.com/193048895737/somename",
	})
	configFname := filepath.Base(viper.ConfigFileUsed())
	viper.SetDefault("S3BUCKET", map[string]string{
		"region": "us-west-2",
		"name":   "mybucket",
		"s3path": configFname, // its in root of bucket
	})
	viper.SetDefault("CLIENT", map[string]string{
		"push_frequency": "1m",
	})
}

func handleConfigReadError(logger *log.Logger, err error) {
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		createConfigFile(logger)
	} else {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}
}

func createConfigFile(logger *log.Logger) {
	logger.Printf("Config file %s not found, creating it with default values...", viper.ConfigFileUsed())

	if err := viper.SafeWriteConfig(); err != nil {
		if os.IsNotExist(err) {
			if err := viper.WriteConfig(); err != nil {
				fmt.Println("Error writing config file:", err)
				os.Exit(1)
			}
		}
	}

	logger.Printf("Config file %s created with default values.", viper.ConfigFileUsed())
}
