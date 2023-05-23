//lint:file-ignore U1000 Return to this when i've pulled my head out of my ass
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// fetchallCmd represents the fetchall command
var fetchallCmd = &cobra.Command{
	Use:   "fetchall",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Trace("fetchall called")
		getAllEc2InstanceLaunchTemplates()
	},
}

type EC2Instance struct {
	InstanceId   string
	InstanceName string
}

type Instance struct {
	ID   string
	Name string
}

func init() {
	LaunchtemplateCmd.AddCommand(fetchallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getInstanceName(ctx context.Context, client *ec2.Client, instanceID string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := client.DescribeInstances(ctx, input)
	if err != nil {
		return "", err
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("no instances found with ID %s", instanceID)
	}

	return aws.ToString(result.Reservations[0].Instances[0].Tags[0].Value), nil
}

// printLaunchTemplateData prints the LaunchTemplateData to stdout
func printLaunchTemplateData(data *types.ResponseLaunchTemplateData) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Logger.Fatalf("failed to marshal LaunchTemplateData to JSON: %v", err)
	}
	log.Logger.Traceln(string(jsonData))
}

func checkInstanceIdExists(instanceID string, client *ec2.Client) (bool, error) {
	// Create the DescribeInstancesInput with the instance ID as a filter
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	// Call the DescribeInstances API
	output, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return false, fmt.Errorf("error describing instances: %w", err)
	}

	// Check if any reservations and instances are returned
	if len(output.Reservations) > 0 && len(output.Reservations[0].Instances) > 0 {
		// Instance exists
		return true, nil
	} else {
		// Instance does not exist
		return false, nil
	}
}

// getLaunchTemplateDataFromInstanceId retrieves the LaunchTemplateData for the specified instance ID
func getLaunchTemplateDataFromInstanceId(ctx context.Context, client *ec2.Client, instanceID string) (*ec2.GetLaunchTemplateDataOutput, error) {
	input := &ec2.GetLaunchTemplateDataInput{InstanceId: aws.String(instanceID)}
	resp, err := client.GetLaunchTemplateData(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get LaunchTemplateData: %v", err)
	}
	return resp, nil
}

func getInstanceIdsForRegion(region string) ([]string, error) {
	svc, err := myec2.GetEc2Client(region)
	if err != nil {
		log.Logger.Errorln(err)
	}

	input := &ec2.DescribeInstancesInput{}

	output, err := svc.DescribeInstances(context.Background(), input)
	if err != nil {
		log.Logger.Traceln("failed to describe EC2 instances:", err)
		return []string{}, err
	}

	var instanceIDs []string

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, *instance.InstanceId)
		}
	}

	jsonData, err := json.MarshalIndent(instanceIDs, "", "  ")
	if err != nil {
		log.Logger.Traceln("failed to marshal list to JSON:", err)
		return []string{}, err
	}

	log.Logger.Traceln(jsonData)
	return instanceIDs, err
}

func getInstanceNameOrExit(ctx context.Context, client *ec2.Client, instanceID string) string {
	instanceName, err := getInstanceName(ctx, client, instanceID)
	if err != nil {
		log.Logger.Traceln("failed to get instance name:", err)
		os.Exit(1)
	}
	return instanceName
}

func getAllEc2InstanceLaunchTemplates() {
	// Get a list of all AWS regions
	regions := myec2.GetAllAwsRegions()

	// Create a buffered channel to limit the number of simultaneous goroutines
	ch := make(chan types.Region, 6)

	// Create a wait group to wait for all goroutines to finish
	wg := sync.WaitGroup{}

	// Iterate over the regions and start a goroutine for each one
	for _, region := range regions {
		// Add the region to the channel
		ch <- region

		// Start a new goroutine
		wg.Add(1)
		go func(region types.Region) {
			// Remove the region from the channel when the goroutine completes
			defer func() {
				<-ch
				wg.Done()
			}()

			// write templates to data/lt-*.json
			err := genLaunchTemplatesForAllEc2InstancesInregion(*region.RegionName)
			if err != nil {
				log.Logger.Fatalln(err)
			}
		}(region)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

// writeLaunchTemplateDataToFile writes the LaunchTemplateData to a JSON file with the specified name
func writeLaunchTemplateDataToFile(data *ec2.GetLaunchTemplateDataOutput, fileName string) error {
	err := common.EnsureParentDirectoryExists(fileName)
	if err != nil {
		log.Logger.Fatalln(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal LaunchTemplateData to JSON: %w", err)
	}

	_, err = file.Write(jsonBytes)
	if err != nil {
		return fmt.Errorf("failed to write LaunchTemplateData to file: %w", err)
	}

	return nil
}

func writeRequestResponseToFile(data *types.ResponseLaunchTemplateData, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal LaunchTemplateData to JSON: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write request response to log file: %w", err)
	}

	return nil
}
