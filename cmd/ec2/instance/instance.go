/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	mydb "github.com/taylormonacelli/deliverhalf/cmd/db"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// InstanceCmd represents the instance command
var InstanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
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
	myec2.Ec2Cmd.AddCommand(InstanceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instanceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instanceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	mydb.Db.AutoMigrate(&ExtendedInstance{})
}

func getJsonDescriptionOfAllInstancesInAllRegions() {
	// Get a list of all AWS regions
	regions := myec2.GetAllAwsRegions()

	// Create a buffered channel to limit the number of simultaneous goroutines
	ch := make(chan types.Region, 3)

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
			describeEc2InstancesInRegionToJson(*region.RegionName)
		}(region)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func describeEc2InstancesInRegionToJson(region string) {
	client, err := myec2.GetEc2Client(region)
	if err != nil {
		log.Logger.Errorln(err)
	}

	input := &ec2.DescribeInstancesInput{
		// Filters: []types.Filter{
		// 	{
		// 		Name:   aws.String("instance-state-name"),
		// 		Values: []string{"running"},
		// 	},
		// },
	}

	resp, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		fmt.Println("Failed to describe instances:", err)
		return
	}

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			fileName := getInstanceFileName(instance)
			path, err := filepath.Abs(filepath.Join("data", "types.instance", fileName))
			if err != nil {
				log.Logger.Errorln(err)
			}
			if common.FileExists(path) {
				log.Logger.Warnf("skipping %s because it already exists", path)
				continue
			}
			if err := writeInstanceDetailsToFile(path, instance); err != nil {
				fmt.Printf("Failed to write instance details to file: %v\n", err)
				continue
			}
			fmt.Printf("Successfully wrote instance details to file: %s\n", path)
		}
	}
}

func testReadInstanceFromJsonFile() (types.Instance, error) {
	fname := "data/types.instance/i-02a6935a1e72ebffd.json"

	inst, err := readInstanceFromJsonFile(fname)
	if err != nil {
		panic(err)
	}
	pp.Print(inst)
	return inst, nil
}

func readInstanceFromJsonFile(pathToJsonFile string) (types.Instance, error) {
	var inst types.Instance

	jsonBlob, err := os.ReadFile(pathToJsonFile)
	if err != nil {
		log.Logger.Fatalf("Reading JSON into byte slice failed with error: %s", err)
	}

	// unmarshal the JSON data into the map
	err = json.Unmarshal(jsonBlob, &inst)
	if err != nil {
		return types.Instance{}, err
	}
	return inst, nil
}

func queryExtendedInstancesFromDb() {
	var extInst ExtendedInstance
	WriteInstancesJsonToFiles()
	testWriteExtendedInstanceJsonDb()
	mydb.Db.First(&extInst, 1)

	var instance types.Instance
	json.Unmarshal([]byte(extInst.JsonDef), &instance)

	fmt.Println(extInst.JsonDef)
	// fmt.Println(instance.InstanceId)

	// Create an empty array to store the strings
	var volumeIds []string

	// Loop over the list and append each item to the array
	for _, item := range instance.BlockDeviceMappings {
		volumeIds = append(volumeIds, *item.Ebs.VolumeId)
	}

	// Print the array volume ids
	for _, volId := range volumeIds {
		fmt.Println(volId)
	}
}

func testWriteExtendedInstanceJsonDb() {
	writeExtendedInstanceJsonDb()
}

func writeExtendedInstanceJsonDb() {
	inst, err := testReadInstanceFromJsonFile()
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(inst, "", "  ")
	if err != nil {
		panic(err)
	}

	var instExt ExtendedInstance
	instExt.JsonDef = string(jsonData)
	instExt.InstanceId = *inst.InstanceId
	mydb.Db.Create(&instExt)
}

func WriteInstancesJsonToFiles() {
	client, err := myec2.GetEc2Client("us-west-2")
	if err != nil {
		log.Logger.Errorln(err)
	}

	input := &ec2.DescribeInstancesInput{
		// Filters: []types.Filter{
		// 	{
		// 		Name:   aws.String("instance-state-name"),
		// 		Values: []string{"running"},
		// 	},
		// },
	}

	resp, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		fmt.Println("Failed to describe instances:", err)
		return
	}

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			fileName := getInstanceFileName(instance)
			path, err := filepath.Abs(filepath.Join("data", "types.instance", fileName))
			if err != nil {
				log.Logger.Errorln(err)
			}
			if common.FileExists(path) {
				log.Logger.Warnf("skipping %s because it already exists", path)
				continue
			}
			if err := writeInstanceDetailsToFile(path, instance); err != nil {
				fmt.Printf("Failed to write instance details to file: %v\n", err)
				continue
			}
			fmt.Printf("Successfully wrote instance details to file: %s\n", path)
		}
	}
}

func getInstanceFileName(instance types.Instance) string {
	return *instance.InstanceId + ".json"
}

func writeInstanceDetailsToFile(fileName string, instance types.Instance) error {
	common.EnsureParentDirectoryExists(fileName)

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file for instance %s: %v", *instance.InstanceId, err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(instance, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for instance %s: %v", *instance.InstanceId, err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write JSON data to file for instance %s: %v", *instance.InstanceId, err)
	}

	return nil
}
