//lint:file-ignore U1000 Return to this when i've pulled my head out of my ass
package cmd

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	mydb "github.com/taylormonacelli/deliverhalf/cmd/db"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	lt "github.com/taylormonacelli/deliverhalf/cmd/ec2/launchtemplate"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// learn4Cmd represents the learn4 command
var learn4Cmd = &cobra.Command{
	Use:   "learn4",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Trace("learn4 called")
		testCreateEc2InstanceFromLaunchTemplateFromDbFromLtName()
	},
}

func init() {
	fiddleCmd.AddCommand(learn4Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// learn4Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// learn4Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func testCreateEc2InstanceFromLaunchTemplateFromDbFromLtName() (*ec2.CreateLaunchTemplateOutput, error) {
	ltName := "taylor-test-deliverhalf"
	var x1 lt.ExtendedGetLaunchTemplateDataOutput

	dbResult := mydb.Db.First(&x1, lt.ExtendedGetLaunchTemplateDataOutput{InstanceName: ltName})
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		log.Logger.Warn("no record found")
		return &ec2.CreateLaunchTemplateOutput{}, nil
	}

	cltInput, err := lt.CreateLaunchTemplateInputFromString(x1.LaunchTemplateDataJsonStr)
	if err != nil {
		log.Logger.Fatalln(err)
	}

	ltName2 := lt.GenRandName(ltName)
	cltInput.LaunchTemplateName = &ltName2
	svc, err := myec2.GetEc2Client(x1.Region)
	if err != nil {
		log.Logger.Fatal(err)
	}

	jsBytes, err := json.MarshalIndent(cltInput, "", "  ")
	if err != nil {
		log.Logger.Warnf("failed to unmarshal template output: %s", err)
	}
	log.Logger.Tracef("launch template %s", string(jsBytes))

	ctOutput, err := svc.CreateLaunchTemplate(context.Background(), cltInput)
	if err != nil {
		log.Logger.Warnf("failed to create launch template: %s", err)
		return &ec2.CreateLaunchTemplateOutput{}, err
	}

	jsBytes, err = json.MarshalIndent(ctOutput, "", "  ")
	if err != nil {
		log.Logger.Warnf("failed to unmarshal template output: %s", err)
	}
	log.Logger.Tracef("launch template created successfully: %s", string(jsBytes))

	return ctOutput, nil
}
