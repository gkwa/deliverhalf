/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
)

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Create ec2 volume snapshot and tag it",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("snapshot called")
		logger := common.SetupLogger()
		createVolumeSnapshot(logger)
		queryRegionForSnapshotsWithTag(logger, "us-west-2")
	},
}

func init() {
	VolumeCmd.AddCommand(snapshotCmd)

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

func genSnapDesc() string {
	snapshotDesc := "created by deliverhalf"
	return snapshotDesc
}

func genSnapTags() []types.Tag {
	tags := []types.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("mytest"),
		},
		{
			Key:   aws.String("Other TagName"),
			Value: aws.String("Other Tag Value"),
		},
		{
			Key:   aws.String("CreatedBy"),
			Value: aws.String("deliverhalf"),
		},
	}
	return tags
}

func createVolumeSnapshot(logger *log.Logger) (string, error) {
	volumeID := "vol-08f2578d51865489b"
	region := "us-west-2"

	snapshotID, err := snapAndTagVolume(logger, volumeID, region)
	if err != nil {
		return "", err
	}

	return snapshotID, err
}

func snapAndTagVolume(logger *log.Logger, volumeID string, region string) (string, error) {
	tags := genSnapTags()
	description := genSnapDesc()

	tagsStr := joinTagsToStr(logger, tags)
	logger.Printf("Creating snapshot with description '%s' for "+
		"volumeID: %s in region: %s and tagging with: '%s'",
		description, volumeID, region, tagsStr)

	snapshotID, err := snapVolume(logger, volumeID, region, description)
	if err != nil {
		logger.Printf("Error snapshotting volume: %s", err)
		return "", err
	}

	logger.Printf("Snapshot created with ID: %s\n", snapshotID)
	err = tagSnapshot(logger, snapshotID, region, tags)
	return "", err
}

func snapVolume(logger *log.Logger, volumeID string, region string, snapshotDesc string) (string, error) {
	cfg, err := myec2.CreateConfig(logger, region)
	if err != nil {
		logger.Fatalf("Could not create config %s", err)
	}

	ec2svc := ec2.NewFromConfig(cfg)

	input := &ec2.CreateSnapshotInput{
		VolumeId:    aws.String(volumeID),
		Description: aws.String(snapshotDesc),
	}

	// Create a snapshot
	resp, err := ec2svc.CreateSnapshot(context.Background(), &ec2.CreateSnapshotInput{
		VolumeId: aws.String(volumeID),
	})
	if err != nil {
		logger.Fatalf("tried to create snapshot for volumeID %s, but got error %s",
			*input.VolumeId, err)
	}

	snapshotID := *resp.SnapshotId
	return snapshotID, err
}

func queryRegionForSnapshotsWithTag(logger *log.Logger, region string) {
	cfg, err := myec2.CreateConfig(logger, region)
	ec2svc := ec2.NewFromConfig(cfg)

	// Get all snapshots with tag key "CreatedBy" and value "deliverhalf"
	input1 := &ec2.DescribeSnapshotsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:CreatedBy"),
				Values: []string{"deliverhalf"},
			},
		},
	}
	output, err := ec2svc.DescribeSnapshots(context.Background(), input1)
	if err != nil {
		fmt.Println("Error listing snapshots:", err)
	}

	// Print out the snapshot IDs
	for _, snapshot := range output.Snapshots {
		fmt.Println(*snapshot.SnapshotId)
	}
}

func tagSnapshot(logger *log.Logger, snapshotID string, region string, tags []types.Tag) error {
	// Add a tag to the snapshot

	tagInput := &ec2.CreateTagsInput{
		Resources: []string{snapshotID},
		Tags:      tags,
	}

	cfg, err := myec2.CreateConfig(logger, region)
	if err != nil {
		logger.Fatalf("Could not create config %s", err)
	}
	ec2svc := ec2.NewFromConfig(cfg)

	_, err = ec2svc.CreateTags(context.Background(), tagInput)
	if err != nil {
		logger.Fatalf("Failed to tag snapshot with ID %s: %v", snapshotID, err)
	} else {
		tagsStr := joinTagsToStr(logger, tags)
		logger.Printf("Successfully tagged snapshot %s with tags %s", snapshotID, tagsStr)
	}
	return err
}

func joinTagsToStr(logger *log.Logger, tags []types.Tag) string {
	var sb strings.Builder
	for _, s := range tags {
		sb.WriteString(*s.Key + "=" + *s.Value + ";")
	}

	result := sb.String()

	return result
}
