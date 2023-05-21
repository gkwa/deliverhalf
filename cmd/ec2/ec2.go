/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
	"github.com/taylormonacelli/deliverhalf/cmd"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// ec2Cmd represents the ec2 command
var Ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Args:  cobra.OnlyValidArgs,
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Traceln("ec2 called")
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
	cmd.RootCmd.AddCommand(Ec2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ec2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ec2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func CreateConfig(region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, err
}

func GetEc2Client(region string) (*ec2.Client, error) {
	config, err := CreateConfig(region)
	if err != nil {
		return nil, err
	}
	// Create an EC2 client
	return ec2.NewFromConfig(config), nil
}

func GetAllAwsRegions() []types.Region {
	// Load the AWS SDK configuration
	client, err := GetEc2Client("us-west-2")
	if err != nil {
		log.Logger.Traceln("failed to load AWS SDK config:", err)
		return []types.Region{}
	}

	// Get a list of all AWS regions
	resp, err := client.DescribeRegions(context.Background(), nil)
	if err != nil {
		log.Logger.Traceln("failed to describe AWS regions:", err)
		return []types.Region{}
	}

	// Create an empty slice of types.Region
	regions := make([]types.Region, 0, len(resp.Regions))
	regions = append(regions, resp.Regions...)

	// Print the region names
	for _, region := range regions {
		log.Logger.Traceln(*region.RegionName)
	}
	return regions
}
