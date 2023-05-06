/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "deliverhalf",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.deliverhalf.yaml)")
	RootCmd.PersistentFlags().String("log-level",
		"info", "Log level (trace, debug, info, warn, error, fatal, panic)",
	)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".deliverhalf" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".deliverhalf")
	}

	SetDefaultValues()

	viper.BindPFlag("log-level", RootCmd.Flags().Lookup("log-level"))
	logLevel := logging.ParseLogLevel(viper.GetString("log-level"))
	logging.Logger.SetLevel(logLevel)

	viper.BindPFlag("config", RootCmd.Flags().Lookup("config"))

	viper.AutomaticEnv() // read in environment variables that match
	viper.ReadInConfig()
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
		"push-frequency": "1m",
	})
}
