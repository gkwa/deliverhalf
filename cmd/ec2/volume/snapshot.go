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
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
)

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("snapshot called")
		logger := common.SetupLogger()
		createSnapshot(logger)

	},
}

func init() {
	volumeCmd.AddCommand(snapshotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// snapshotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// snapshotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Load the AWS SDK configuration

	// Load the AWS SDK configuration

}

func createSnapshot(logger *log.Logger) {
	volumeID := "vol-08f2578d51865489b"
	snapshotDesc := "created by deliverhalf"
	region := "us-west-2"

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
	ec2svc := ec2.NewFromConfig(cfg)

	input := &ec2.CreateSnapshotInput{
		VolumeId:    aws.String(volumeID),
		Description: aws.String(snapshotDesc),
	}

	resp, err := ec2svc.CreateSnapshot(context.Background(), input)
	if err != nil {
		panic("failed to create EBS snapshot")
	}

	snapshotID := *resp.SnapshotId
	fmt.Printf("Snapshot created with ID: %s\n", snapshotID)
	tagSnapshot(logger, snapshotID, region)
}

func tagSnapshot(logger *log.Logger, snapshotID string, region string) {
	// Add a tag to the snapshot
	tagInput := &ec2.CreateTagsInput{
		Resources: []string{snapshotID},
		Tags: []types.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("mytest"),
			},
		},
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
	ec2svc := ec2.NewFromConfig(cfg)

	_, err = ec2svc.CreateTags(context.Background(), tagInput)
	if err != nil {
		logger.Printf("failed to tag snapshot with ID %s: %v", snapshotID, err)
	} else {
		logger.Printf("successfully tagged snapshot with ID %s", snapshotID)
	}
}