//lint:file-ignore U1000 Return to this when i've pulled my head out of my ass

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
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
		// testCreateLaunchTemplateFromFile()
		testCreateLaunchTemplateOutputFromFile()
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

func getLaunchTemplateDataOutputFromString(ltString string) (*ec2.GetLaunchTemplateDataOutput, error) {
	// Unmarshal the JSON into a ResponseLaunchTemplateData struct
	var ltData *ec2.GetLaunchTemplateDataOutput
	err := json.Unmarshal([]byte(ltString), &ltData)
	if err != nil {
		log.Logger.Warnln("couldn't unmarshal launchtemplate from string")
		return &ec2.GetLaunchTemplateDataOutput{}, err
	}

	return ltData, nil
}

func getLaunchTemplateDataOutputFromFile(path string) (*ec2.GetLaunchTemplateDataOutput, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Logger.Errorln(err)
	}

	fileContents, err := os.ReadFile(absPath)
	if err != nil {
		log.Logger.Warnf("could not read path %s", absPath)
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
			x1, err := filepath.Abs(filepath.Join(dir, file.Name()))
			if err != nil {
				log.Logger.Error(err)
			}
			matchingFiles = append(matchingFiles, x1)
		}
	}

	log.Logger.Traceln(matchingFiles)
	return matchingFiles, nil
}

func testCreateLaunchTemplateFromFile() error {
	fname := "data/GetLaunchTemplateDataOutput/lt-i-07e53ad5c1747dd52.json"
	path, err := filepath.Abs(fname)
	if err != nil {
		log.Logger.Errorln(err)
		return err
	}

	cltInput, err := CreateLaunchTemplateInputFromFile(path)
	if err != nil {
		log.Logger.Fatalln(err)
		return err
	}

	jsonData, err := json.MarshalIndent(cltInput, "", "  ")
	if err != nil {
		fmt.Println("error marshaling struct to JSON:", err)
		return err
	}

	log.Logger.WithField("data", string(jsonData)).Trace("log indented JSON")

	return nil
}

func CreateLaunchTemplateFromFile(path string) (*types.LaunchTemplate, error) {
	output, err := CreateLaunchTemplateOutputFromFile(path)
	if err != nil {
		return &types.LaunchTemplate{}, err
	}

	return output.LaunchTemplate, nil
}

func readFileToString(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	content := ""
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return content, nil
}

func CreateLaunchTemplateInputFromFile(ltPath string) (*ec2.CreateLaunchTemplateInput, error) {
	content, err := readFileToString(ltPath)
	if err != nil {
		log.Logger.Fatalf("can't read file %s", ltPath)
	}
	cltInput, err := CreateLaunchTemplateInputFromString(content)
	if err != nil {
		log.Logger.Fatalln("can't create launch template input", err)
	}
	return cltInput, nil
}

func CreateLaunchTemplateInputFromString(ltOutput string) (*ec2.CreateLaunchTemplateInput, error) {
	myI, err := getLaunchTemplateDataOutputFromString(ltOutput)
	if err != nil {
		log.Logger.Fatalf("failed to create launch template from string")
	}

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
		n := len(myI.LaunchTemplateData.NetworkInterfaces[i].Ipv6Addresses)
		niSpecs[i] = types.LaunchTemplateInstanceNetworkInterfaceSpecificationRequest{
			AssociatePublicIpAddress:       ni.AssociatePublicIpAddress,
			DeleteOnTermination:            ni.DeleteOnTermination,
			Description:                    ni.Description,
			DeviceIndex:                    ni.DeviceIndex,
			Groups:                         ni.Groups,
			InterfaceType:                  ni.InterfaceType,
			Ipv6AddressCount:               ni.Ipv6AddressCount,
			Ipv6Addresses:                  make([]types.InstanceIpv6AddressRequest, n),
			NetworkCardIndex:               ni.DeviceIndex,
			NetworkInterfaceId:             ni.NetworkInterfaceId,
			PrivateIpAddress:               ni.PrivateIpAddress,
			PrivateIpAddresses:             ni.PrivateIpAddresses,
			SecondaryPrivateIpAddressCount: ni.SecondaryPrivateIpAddressCount,
			SubnetId:                       ni.SubnetId,
		}
	}

	// Convert the block device mappings to the correct type
	var bdMappings []types.LaunchTemplateBlockDeviceMappingRequest
	for _, bdMapping := range myI.LaunchTemplateData.BlockDeviceMappings {
		bdMappings = append(bdMappings, types.LaunchTemplateBlockDeviceMappingRequest{
			DeviceName: bdMapping.DeviceName,
			Ebs: &types.LaunchTemplateEbsBlockDeviceRequest{
				VolumeSize:          bdMapping.Ebs.VolumeSize,
				VolumeType:          bdMapping.Ebs.VolumeType,
				DeleteOnTermination: bdMapping.Ebs.DeleteOnTermination,
				Encrypted:           bdMapping.Ebs.Encrypted,
				SnapshotId:          bdMapping.Ebs.SnapshotId,
			},
		})
	}

	// Convert CapacityReservationSpecification from response to request type
	capacityReservationSpec := &types.LaunchTemplateCapacityReservationSpecificationRequest{
		CapacityReservationPreference: myI.LaunchTemplateData.CapacityReservationSpecification.CapacityReservationPreference,
	}

	metadataOptions := &types.LaunchTemplateInstanceMetadataOptionsRequest{
		HttpEndpoint:            myI.LaunchTemplateData.MetadataOptions.HttpEndpoint,
		HttpPutResponseHopLimit: myI.LaunchTemplateData.MetadataOptions.HttpPutResponseHopLimit,
		HttpTokens:              myI.LaunchTemplateData.MetadataOptions.HttpTokens,
		InstanceMetadataTags:    myI.LaunchTemplateData.MetadataOptions.InstanceMetadataTags,
		HttpProtocolIpv6:        myI.LaunchTemplateData.MetadataOptions.HttpProtocolIpv6,
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

	instanceMarketOptions := &types.LaunchTemplateInstanceMarketOptionsRequest{
		MarketType: types.MarketTypeSpot,
	}

	if instanceMarketOptions.MarketType == types.MarketTypeSpot {
		*myI.LaunchTemplateData.DisableApiStop = false
		*myI.LaunchTemplateData.DisableApiTermination = false
	}

	// Create a new RequestLaunchTemplateData object based on the retrieved launch template data
	requestData := &types.RequestLaunchTemplateData{
		BlockDeviceMappings:               bdMappings,
		CapacityReservationSpecification:  capacityReservationSpec,
		CpuOptions:                        (*types.LaunchTemplateCpuOptionsRequest)(myI.LaunchTemplateData.CpuOptions),
		CreditSpecification:               (*types.CreditSpecificationRequest)(myI.LaunchTemplateData.CreditSpecification),
		DisableApiStop:                    myI.LaunchTemplateData.DisableApiStop,
		DisableApiTermination:             myI.LaunchTemplateData.DisableApiTermination,
		EbsOptimized:                      myI.LaunchTemplateData.EbsOptimized,
		ElasticGpuSpecifications:          elasticGpuSpecs,
		EnclaveOptions:                    (*types.LaunchTemplateEnclaveOptionsRequest)(myI.LaunchTemplateData.EnclaveOptions),
		HibernationOptions:                (*types.LaunchTemplateHibernationOptionsRequest)(myI.LaunchTemplateData.HibernationOptions),
		IamInstanceProfile:                (*types.LaunchTemplateIamInstanceProfileSpecificationRequest)(myI.LaunchTemplateData.IamInstanceProfile),
		ImageId:                           myI.LaunchTemplateData.ImageId,
		InstanceInitiatedShutdownBehavior: myI.LaunchTemplateData.InstanceInitiatedShutdownBehavior,
		InstanceMarketOptions:             instanceMarketOptions,
		InstanceRequirements:              instanceRequirementsRequest,
		InstanceType:                      myI.LaunchTemplateData.InstanceType,
		KernelId:                          myI.LaunchTemplateData.KernelId,
		KeyName:                           myI.LaunchTemplateData.KeyName,
		MaintenanceOptions:                (*types.LaunchTemplateInstanceMaintenanceOptionsRequest)(myI.LaunchTemplateData.MaintenanceOptions),
		MetadataOptions:                   metadataOptions,
		Monitoring:                        (*types.LaunchTemplatesMonitoringRequest)(myI.LaunchTemplateData.Monitoring),
		NetworkInterfaces:                 niSpecs,
		Placement:                         (*types.LaunchTemplatePlacementRequest)(myI.LaunchTemplateData.Placement),
		PrivateDnsNameOptions:             (*types.LaunchTemplatePrivateDnsNameOptionsRequest)(myI.LaunchTemplateData.PrivateDnsNameOptions),
		RamDiskId:                         myI.LaunchTemplateData.RamDiskId,
		SecurityGroupIds:                  myI.LaunchTemplateData.SecurityGroupIds,
		SecurityGroups:                    myI.LaunchTemplateData.SecurityGroups,
		TagSpecifications:                 tagSpecs,
		UserData:                          myI.LaunchTemplateData.UserData,
	}

	// Create the launch template
	ltName := genRandName("deliverhalf")
	cltInput := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateName: &ltName,
		LaunchTemplateData: requestData,
	}

	return cltInput, nil
}

func CreateLaunchTemplateOutputFromFile(ltPath string) (*ec2.CreateLaunchTemplateOutput, error) {
	cltInput, err := CreateLaunchTemplateInputFromFile(ltPath)
	if err != nil {
		log.Logger.Fatalln(err)
	}

	svc, err := myec2.GetEc2Client("us-west-2")
	if err != nil {
		log.Logger.Fatal(err)
	}

	ctOutput, err := svc.CreateLaunchTemplate(context.Background(), cltInput)
	if err != nil {
		log.Logger.Warnf("failed to create launch template: %s", err)
		return &ec2.CreateLaunchTemplateOutput{}, err
	}

	jsBytes, err := json.MarshalIndent(ctOutput, "", "  ")
	if err != nil {
		log.Logger.Warnf("failed to unmarshal template output: %s", err)
	}
	log.Logger.Tracef("launch template created successfully: %s", string(jsBytes))

	return ctOutput, nil
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

	jsBytes, err := json.MarshalIndent(createTemplateOutput, "", "  ")
	if err != nil {
		log.Logger.Warnf("failed to unmarshal tempalte output: %s", err)
	}
	log.Logger.Tracef("Launch template created successfully: %s", string(jsBytes))
}

func testCreateLaunchTemplateOutputFromFile() error {
	fname := "data/GetLaunchTemplateDataOutput/lt-i-07e53ad5c1747dd52.json"
	path, err := filepath.Abs(fname)
	if err != nil {
		log.Logger.Errorln(err)
		return err
	}

	cltOutput, err := CreateLaunchTemplateOutputFromFile(path)
	if err != nil {
		log.Logger.Fatalln(err)
		return err
	}

	jsonData, err := json.MarshalIndent(cltOutput, "", "  ")
	if err != nil {
		fmt.Println("error marshaling struct to JSON:", err)
		return err
	}

	log.Logger.WithField("data", string(jsonData)).Trace("log indented JSON")

	return nil
}
