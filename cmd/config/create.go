/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		setupConfig()
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

func setupConfig() {
	addConfigPaths()
	setConfigNameAndType()
	setDefaultValues()

	if err := viper.ReadInConfig(); err != nil {
		handleConfigReadError(err)
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

func setDefaultValues() {
	viper.SetDefault("SNS", map[string]string{"topic-arn": "arn:aws:sns:us-west-2:123456789012:example-topic", "region": "us-west-2"})
}

func handleConfigReadError(err error) {
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		createConfigFile()
	} else {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}
}

func createConfigFile() {
	fmt.Println("Config file not found, creating it with default values...")

	if err := viper.SafeWriteConfig(); err != nil {
		if os.IsNotExist(err) {
			if err := viper.WriteConfig(); err != nil {
				fmt.Println("Error writing config file:", err)
				os.Exit(1)
			}
		}
	}

	fmt.Println("Config file created with default values.")
}
