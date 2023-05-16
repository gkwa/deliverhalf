/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
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

	if err := viper.ReadInConfig(); err != nil {
		handleConfigReadError(err)
	}
}

func addConfigPaths() {
	viper.AddConfigPath(".")
}

func setConfigNameAndType() {
	viper.SetConfigName(".deliverhalf")
	viper.SetConfigType("yaml")
}

func handleConfigReadError(err error) {
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		createConfigFile()
	} else {
		log.Logger.Traceln("Error reading config file:", err)
		os.Exit(1)
	}
}

func createConfigFile() {
	log.Logger.Printf("Config file %s not found, creating it with default values...", viper.ConfigFileUsed())

	if err := viper.SafeWriteConfig(); err != nil {
		if os.IsNotExist(err) {
			if err := viper.WriteConfig(); err != nil {
				log.Logger.Traceln("Error writing config file:", err)
				os.Exit(1)
			}
		}
	}

	log.Logger.Printf("Config file %s created with default values.", viper.ConfigFileUsed())
}
