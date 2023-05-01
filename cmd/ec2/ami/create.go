/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
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
		fmt.Println("create called")
		logger := common.SetupLogger()
		test(logger)
	},
}

func init() {
	amiCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func test(logger *log.Logger) {
	region := "us-west-2"
	snapID := "snap-081fa6c4cec7f21f7"
	createFromSnapshot(logger, snapID, region)
}

func createFromSnapshot(logger *log.Logger, snapshotID string, region string) error {
	cfg := createConfig(logger, region)
	ec2svc := ec2.NewFromConfig(cfg)

	logger.Println(ec2svc)
	return nil
}

func createConfig(logger *log.Logger, region string) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
	return cfg
}
