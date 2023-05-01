/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
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
		test(logger)
	},
}

type AMI struct {
	Name       string
	ImageID    string
	SnapshotID string
	Region     string
}

func (ami AMI) String() string {
	return fmt.Sprintf("Name: %s, ImageID: %s, Region: %s, SnapshotID: %s",
		ami.Name,
		ami.ImageID,
		ami.Region,
		ami.SnapshotID,
	)
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
	amiName := fmt.Sprintf("my-image-%08d", time.Now().Unix())
	ami := AMI{
		Name:       amiName,
		SnapshotID: "snap-081fa6c4cec7f21f7",
		Region:     "us-west-2",
	}
	createAMIFromSnapshot(logger, &ami)
	logger.Printf("Created AMI with properties %s", ami)
}

func createAMIFromSnapshot(logger *log.Logger, ami *AMI) error {
	cfg, err := myec2.CreateConfig(logger, ami.Region)
	if err != nil {
		logger.Fatal(fmt.Errorf("error trying to create ami from snapshot id %s: %s", ami.SnapshotID, err))
	}

	ec2svc := ec2.NewFromConfig(cfg)

	// Call the RegisterImage function to register the AMI image
	input := &ec2.RegisterImageInput{
		Name:         aws.String(ami.Name),
		Architecture: types.ArchitectureValuesX8664,
		BlockDeviceMappings: []types.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				Ebs: &types.EbsBlockDevice{
					SnapshotId:          aws.String(ami.SnapshotID),
					VolumeSize:          aws.Int32(10),
					DeleteOnTermination: aws.Bool(true),
				},
			},
		},
		Description:    aws.String("created by deliverhalf"),
		RootDeviceName: aws.String("/dev/sda1"),
	}

	result, err := ec2svc.RegisterImage(context.Background(), input)
	if err != nil {
		logger.Fatalf("Error registering new AMI with snapshotID %s: %s", ami.SnapshotID, err)
	}
	ami.ImageID = *result.ImageId

	return err
}
