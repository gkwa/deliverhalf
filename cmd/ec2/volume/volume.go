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
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	meta "github.com/taylormonacelli/deliverhalf/cmd/meta"
)

// volumeCmd represents the volume command
var VolumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("volume called")
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		os.Exit(1)
		return nil
	},
}

func init() {
	myec2.Ec2Cmd.AddCommand(VolumeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// volumeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// volumeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type VolumeTag struct {
	Tags []types.Tag
	Size int32
}

func getVolumes(logger *log.Logger) {
	blob := meta.ParseJsonFromFile(logger, "meta.json")
	instanceId := string(blob["instanceId"].(string))
	regionName := string(blob["region"].(string))
	logger.Printf("found instance id %s in region %s", instanceId, regionName)

	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(regionName))
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	// Create a new EC2 client
	svc := ec2.NewFromConfig(cfg)

	volumes, err := getVolumesForInstance(context.Background(), svc, instanceId)
	if err != nil {
		panic(err)
	}

	volumeTags := extractVolumeTags(volumes)

	printVolumeTags(volumeTags)
}

func extractVolumeTags(volumes []types.Volume) map[string]VolumeTag {
	volumeTags := make(map[string]VolumeTag)

	for _, volume := range volumes {
		volumeID := *volume.VolumeId
		tags := volume.Tags
		size := *volume.Size

		volumeTag := VolumeTag{
			Tags: tags,
			Size: size,
		}

		volumeTags[volumeID] = volumeTag
	}

	return volumeTags
}

func printVolumeTags(volumeTags map[string]VolumeTag) {
	for volumeID, volumeTag := range volumeTags {
		fmt.Printf("Volume ID: %s\n", volumeID)
		fmt.Printf("Size: %d\n", volumeTag.Size)
		fmt.Printf("Tags:\n")
		for _, tag := range volumeTag.Tags {
			fmt.Printf("  - %s: %s\n", *tag.Key, *tag.Value)
		}
		fmt.Println()
	}
}

func getVolumesForInstance(ctx context.Context, svc *ec2.Client, instanceID string) ([]types.Volume, error) {
	input := &ec2.DescribeVolumesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("attachment.instance-id"),
				Values: []string{instanceID},
			},
		},
	}
	resp, err := svc.DescribeVolumes(ctx, input)
	if err != nil {
		return nil, err
	}
	return resp.Volumes, nil
}
