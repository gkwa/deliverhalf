/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Trace("create called")
		// create()
		// createLaunchTemplateFromFile()
		// getLaunchDataFromAllTemplatesInDirectory()
		// getLaunchDataFromAllTemplates()
		// testCreate1()
		testCreateLaunchTemplateFromFile()
		// testCreate4()
	},
}

func init() {
	LaunchtemplateCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getConfig() (map[string]interface{}, error) {
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// Return a map of all the settings loaded by Viper
	return viper.AllSettings(), nil
}

func genRandName(prefix string) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	randomNumber := rng.Intn(1000000)
	ltName := fmt.Sprintf("%s-%06d", prefix, randomNumber)
	return ltName
}

func getLaunchTemplateFromName(ltName string) (map[string]interface{}, error) {
	config, err := getConfig()
	if err != nil {
		log.Logger.Fatal(err)
	}

	lt, ok := config["launch_templates"].(map[string]interface{})[ltName].(map[string]interface{})
	if !ok {
		msg := fmt.Sprintf("could not index %s.%s from %s",
			"launch_templates", ltName, viper.ConfigFileUsed())
		return nil, errors.New(msg)
	}
	return lt, nil
}

func create() {
	ltNameFromConfig := "test1"
	lt, err := getLaunchTemplateFromName(ltNameFromConfig)
	if err != nil {
		log.Logger.Fatalf("lookup template from '%s' failed", ltNameFromConfig)
	}

	region := lt["region"].(string)
	ltName := genRandName("deliverhalf")

	client, err := myec2.GetEc2Client(region)
	if err != nil {
		log.Logger.Errorln(err)
	}

	// Get the value from Viper and convert it to the custom type
	var instanceType types.InstanceType

	configIndex := fmt.Sprintf("%s.%s", "launch_templates", ltNameFromConfig)
	err = viper.UnmarshalKey(lt["instancetype"].(string), &instanceType)
	if err != nil {
		log.Logger.Fatalf("failed to unmarshal template %s", ltNameFromConfig)
	}

	// Specify the details of the launch template to create
	imageID := lt["imageid"].(string)
	keyName := "my-key-pair"

	sgIndex := fmt.Sprintf("%s.securitygroupids", configIndex)
	iamInstanceProfileIndex := fmt.Sprintf("%s.iaminstanceprofile", configIndex)
	log.Logger.Tracef("instance profile: %s", iamInstanceProfileIndex)

	instanceProfileName := viper.GetString(iamInstanceProfileIndex)
	securityGroupIDs := viper.GetStringSlice(sgIndex)
	if len(securityGroupIDs) == 0 {
		log.Logger.Warnf("no security groups matched %s from %s, using new security group",
			sgIndex, viper.ConfigFileUsed())
	}

	log.Logger.Tracef("security groups: %s", strings.Join(securityGroupIDs, ", "))

	script, err := createUserData()
	if err != nil {
		log.Logger.Fatalf("can't create userdata script, error: %s", err)
	}
	userDataEncoded := base64.StdEncoding.EncodeToString([]byte(script))

	// Create the launch template
	createLaunchTemplateInput := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateName: &ltName,
		LaunchTemplateData: &types.RequestLaunchTemplateData{
			ImageId:               &imageID,
			InstanceType:          instanceType,
			KeyName:               &keyName,
			SecurityGroupIds:      securityGroupIDs,
			UserData:              &userDataEncoded,
			IamInstanceProfile:    &types.LaunchTemplateIamInstanceProfileSpecificationRequest{Name: &instanceProfileName},
			InstanceMarketOptions: &types.LaunchTemplateInstanceMarketOptionsRequest{MarketType: "spot"},
		},
	}

	createLaunchTemplateOutput, err := client.CreateLaunchTemplate(context.Background(), createLaunchTemplateInput)
	if err != nil {
		log.Logger.Fatal("failed to create launch template, " + err.Error())
	}

	log.Logger.Trace("Launch template created with ID:", *createLaunchTemplateOutput.LaunchTemplate.LaunchTemplateId)
}

func getLaunchDataFromAllTemplates() {
	launchTemplates, err := getLaunchDataFromAllTemplatesInDirectory()
	if err != nil {
		log.Logger.Fatalln("failed to get all launch templates")
	}

	for _, lt := range launchTemplates {
		if lt.LaunchTemplateData.UserData == nil {
			log.Logger.Tracef("launch template doesn't have userdata: %s", *lt.LaunchTemplateData.ImageId)
		} else {
			decoded, err := base64.StdEncoding.DecodeString(*lt.LaunchTemplateData.UserData)
			if err != nil {
				log.Logger.Fatalf("Failed to decode base64 string")
			}
			log.Logger.Tracef("instance %s launch template user data: %s",
				*lt.LaunchTemplateData.ImageId, string(decoded))
		}
	}
}

func getLaunchDataFromAllTemplatesInDirectory() ([]ec2.GetLaunchTemplateDataOutput, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return []ec2.GetLaunchTemplateDataOutput{}, err
	}
	dir := filepath.Join(cwd, "data")

	paths, err := getPathsToMarshalledLaunchTemplates(dir)
	if err != nil {
		log.Logger.Fatalf("could not get list of launch templates in directory %s", dir)
	}

	var launchTemplates []ec2.GetLaunchTemplateDataOutput
	for _, path := range paths {
		lt, err := getLaunchTemplateDataOutputFromFile(path)
		if err != nil {
			log.Logger.Warnf("could not read path %s", path)
			return nil, err
		}
		launchTemplates = append(launchTemplates, *lt)
	}
	return launchTemplates, nil
}

func getLaunchTemplateDataOutputFromFile(path string) (*ec2.GetLaunchTemplateDataOutput, error) {
	fileContents, err := os.ReadFile(path)
	if err != nil {
		log.Logger.Warnf("could not read path %s", path)
		return &ec2.GetLaunchTemplateDataOutput{}, err
	}

	// Unmarshal the JSON into a ResponseLaunchTemplateData struct
	var ltData *ec2.GetLaunchTemplateDataOutput
	err = json.Unmarshal(fileContents, &ltData)
	if err != nil {
		log.Logger.Warnf("couldn't unmarshal launchtemplate from file %s", path)
		return &ec2.GetLaunchTemplateDataOutput{}, err
	}

	return ltData, nil
}

func getPathsToMarshalledLaunchTemplates(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Logger.Fatal(err)
	}

	var matchingFiles []string

	pattern := regexp.MustCompile(`^lt-i-[a-zA-Z0-9]{,16}\.json$`)

	for _, file := range files {
		if !file.IsDir() && pattern.MatchString(file.Name()) {
			matchingFiles = append(matchingFiles, filepath.Join(dir, file.Name()))
		}
	}

	log.Logger.Traceln(matchingFiles)
	return matchingFiles, nil
}

func testCreateLaunchTemplateFromFile() (*types.LaunchTemplate, error) {
	fname := "data/GetLaunchTemplateDataOutput/lt-i-0a026f9c40b0337ca.json"
	getLaunchTemplateDataOutputFile, err := filepath.Abs(fname)
	if err != nil {
		log.Logger.Errorln(err)
		return &types.LaunchTemplate{}, err
	}

	template, err := CreateLaunchTemplateFromFile(getLaunchTemplateDataOutputFile)
	if err != nil {
		return &types.LaunchTemplate{}, err
	}
	log.Logger.Debugf("created launch template %s with id %s from file %s",
		*template.LaunchTemplateName, *template.LaunchTemplateId, getLaunchTemplateDataOutputFile)
	return template, nil
}

func CreateLaunchTemplateFromFile(path string) (*types.LaunchTemplate, error) {
	createLaunchTemplateOutput, err := CreateLaunchTemplateOutputFromFile(path)
	if err != nil {
		return &types.LaunchTemplate{}, err
	}

	return createLaunchTemplateOutput.LaunchTemplate, nil
}

func CreateLaunchTemplateOutputFromFile(ltPath string) (*ec2.CreateLaunchTemplateOutput, error) {
	myI, err := getLaunchTemplateDataOutputFromFile(ltPath)
	if err != nil {
		log.Logger.Fatalf("failed to create launch tempalte from file %s", ltPath)
	}

	pp.Print(myI)
	pp.Print(myI.LaunchTemplateData)
	pp.Print(myI.LaunchTemplateData.TagSpecifications)

	// Convert TagSpecifications to LaunchTemplateTagSpecificationRequest
	tagSpecs := make([]types.LaunchTemplateTagSpecificationRequest, len(myI.LaunchTemplateData.TagSpecifications))
	for i, ts := range myI.LaunchTemplateData.TagSpecifications {
		tagSpecs[i] = types.LaunchTemplateTagSpecificationRequest{
			ResourceType: ts.ResourceType,
			Tags:         ts.Tags,
		}
	}

	// Convert NetworkInterfaces to LaunchTemplateInstanceNetworkInterfaceSpecificationRequest
	niSpecs := make([]types.LaunchTemplateInstanceNetworkInterfaceSpecificationRequest, len(myI.LaunchTemplateData.NetworkInterfaces))
	for i, ni := range myI.LaunchTemplateData.NetworkInterfaces {
		niSpecs[i] = types.LaunchTemplateInstanceNetworkInterfaceSpecificationRequest{
			AssociatePublicIpAddress: ni.AssociatePublicIpAddress,
			DeleteOnTermination:      ni.DeleteOnTermination,
			Description:              ni.Description,
			DeviceIndex:              ni.DeviceIndex,
			Groups:                   ni.Groups,
			InterfaceType:            ni.InterfaceType,
			Ipv6AddressCount:         ni.Ipv6AddressCount,
			NetworkInterfaceId:       ni.NetworkInterfaceId,
			PrivateIpAddress:         ni.PrivateIpAddress,
			// PrivateIpAddresses:             ni.PrivateIpAddresses,
			SecondaryPrivateIpAddressCount: ni.SecondaryPrivateIpAddressCount,
			SubnetId:                       ni.SubnetId,
			// Ipv6Addresses:                  ni.Ipv6Addresses,
		}
	}

	// Convert the block device mappings to the correct type
	var bdMappings []types.LaunchTemplateBlockDeviceMappingRequest
	for _, bdMapping := range myI.LaunchTemplateData.BlockDeviceMappings {
		bdMappings = append(bdMappings, types.LaunchTemplateBlockDeviceMappingRequest{
			DeviceName: bdMapping.DeviceName,
			Ebs: &types.LaunchTemplateEbsBlockDeviceRequest{
				VolumeSize: bdMapping.Ebs.VolumeSize,
				VolumeType: bdMapping.Ebs.VolumeType,
			},
		})
	}

	svc, err := myec2.GetEc2Client("us-west-2")
	if err != nil {
		log.Logger.Fatal(err)
	}
	ltName := genRandName("deliverhalf")

	// Convert CapacityReservationSpecification from response to request type
	capacityReservationSpec := &types.LaunchTemplateCapacityReservationSpecificationRequest{
		CapacityReservationPreference: myI.LaunchTemplateData.CapacityReservationSpecification.CapacityReservationPreference,
	}

	metadataOptions := &types.LaunchTemplateInstanceMetadataOptionsRequest{
		HttpEndpoint:            myI.LaunchTemplateData.MetadataOptions.HttpEndpoint,
		HttpPutResponseHopLimit: myI.LaunchTemplateData.MetadataOptions.HttpPutResponseHopLimit,
		HttpTokens:              myI.LaunchTemplateData.MetadataOptions.HttpTokens,
	}

	var elasticGpuSpecs []types.ElasticGpuSpecification

	for _, gpuSpec := range myI.LaunchTemplateData.ElasticGpuSpecifications {
		elasticGpuSpec := types.ElasticGpuSpecification{
			Type: gpuSpec.Type,
			// Assign other fields as needed
		}

		elasticGpuSpecs = append(elasticGpuSpecs, elasticGpuSpec)
	}

	instanceRequirements := myI.LaunchTemplateData.InstanceRequirements
	var instanceRequirementsRequest *types.InstanceRequirementsRequest
	if instanceRequirements != nil {
		instanceRequirementsRequest = &types.InstanceRequirementsRequest{
			// Assign fields from instanceRequirements to instanceRequirementsRequest
			// Assign other fields as needed
		}
	}

	// Create a new RequestLaunchTemplateData object based on the retrieved launch template data
	requestData := &types.RequestLaunchTemplateData{
		// DisableApiStop:                    myI.LaunchTemplateData.DisableApiStop,  // not with spot instance
		// DisableApiTermination:             myI.LaunchTemplateData.DisableApiStop,  // not with spot instance
		BlockDeviceMappings:               bdMappings,
		CapacityReservationSpecification:  capacityReservationSpec,
		CpuOptions:                        (*types.LaunchTemplateCpuOptionsRequest)(myI.LaunchTemplateData.CpuOptions),
		CreditSpecification:               (*types.CreditSpecificationRequest)(myI.LaunchTemplateData.CreditSpecification),
		EbsOptimized:                      myI.LaunchTemplateData.EbsOptimized,
		ElasticGpuSpecifications:          elasticGpuSpecs,
		EnclaveOptions:                    (*types.LaunchTemplateEnclaveOptionsRequest)(myI.LaunchTemplateData.EnclaveOptions),
		HibernationOptions:                (*types.LaunchTemplateHibernationOptionsRequest)(myI.LaunchTemplateData.HibernationOptions),
		IamInstanceProfile:                (*types.LaunchTemplateIamInstanceProfileSpecificationRequest)(myI.LaunchTemplateData.IamInstanceProfile),
		ImageId:                           myI.LaunchTemplateData.ImageId,
		InstanceInitiatedShutdownBehavior: myI.LaunchTemplateData.InstanceInitiatedShutdownBehavior,
		// InstanceMarketOptions:             &types.LaunchTemplateInstanceMarketOptionsRequest{MarketType: "spot"},
		InstanceRequirements:  instanceRequirementsRequest,
		InstanceType:          myI.LaunchTemplateData.InstanceType,
		KernelId:              myI.LaunchTemplateData.KernelId,
		KeyName:               myI.LaunchTemplateData.KeyName,
		MaintenanceOptions:    (*types.LaunchTemplateInstanceMaintenanceOptionsRequest)(myI.LaunchTemplateData.MaintenanceOptions),
		MetadataOptions:       metadataOptions,
		Monitoring:            (*types.LaunchTemplatesMonitoringRequest)(myI.LaunchTemplateData.Monitoring),
		NetworkInterfaces:     niSpecs,
		Placement:             (*types.LaunchTemplatePlacementRequest)(myI.LaunchTemplateData.Placement),
		PrivateDnsNameOptions: (*types.LaunchTemplatePrivateDnsNameOptionsRequest)(myI.LaunchTemplateData.PrivateDnsNameOptions),
		RamDiskId:             myI.LaunchTemplateData.RamDiskId,
		SecurityGroupIds:      myI.LaunchTemplateData.SecurityGroupIds,
		SecurityGroups:        myI.LaunchTemplateData.SecurityGroups,
		TagSpecifications:     tagSpecs,
		UserData:              myI.LaunchTemplateData.UserData,
	}
	requestJson, err := json.Marshal(requestData)
	if err != nil {
		log.Logger.Warnf("could not marshal requestData: %s", err)
	}

	log.Logger.Tracef("RequestLaunchTemplateData object created successfully: %s", string(requestJson))

	// Create the launch template
	createTemplateInput := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateName: &ltName,
		LaunchTemplateData: requestData,
	}

	createTemplateOutput, err := svc.CreateLaunchTemplate(context.Background(), createTemplateInput)
	if err != nil {
		log.Logger.Warnf("failed to create launch template: %s", err)
		return &ec2.CreateLaunchTemplateOutput{}, err
	}

	jsBytes, err := json.MarshalIndent(createTemplateOutput, "", "  ")
	if err != nil {
		log.Logger.Warnf("failed to unmarshal tempalte output: %s", err)
	}
	log.Logger.Tracef("Launch template created successfully: %s", string(jsBytes))
	fmt.Println("Launch template created successfully: " + string(jsBytes))
	return createTemplateOutput, nil
}

func testCreate4() {
	region := "us-west-2"
	svc, err := myec2.GetEc2Client(region)
	if err != nil {
		panic(err)
	}

	path := "data/i-0476d67631ffc9996-LaunchTemplate.json"
	path = "/Users/mtm/pdev/taylormonacelli/deliverhalf/data/i-0476d67631ffc9996-LaunchTemplate.json"
	fileContents, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var ltData types.RequestLaunchTemplateData
	err = json.Unmarshal(fileContents, &ltData)
	if err != nil {
		panic(err)
	}
	pp.Print(ltData.Placement)

	// instanceId := "i-0476d67631ffc9996"
	myname := "mytest"
	input1 := ec2.CreateLaunchTemplateInput{
		LaunchTemplateName: &myname,
		LaunchTemplateData: &ltData,
	}
	createTemplateOutput, err := svc.CreateLaunchTemplate(context.Background(), &input1)
	if err != nil {
		log.Logger.Warnf("failed to create launch template: %s", err)
		return
	}

	jsBytes, err := json.Marshal(createTemplateOutput)
	if err != nil {
		log.Logger.Warnf("failed to unmarshal tempalte output: %s", err)
	}
	log.Logger.Tracef("Launch template created successfully: %s", string(jsBytes))
}
