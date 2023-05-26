//lint:file-ignore U1000 Return to this when i've pulled my head out of my ass
package cmd

import (
	"errors"

	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	mydb "github.com/taylormonacelli/deliverhalf/cmd/db"
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

func testCreateEc2InstanceFromLaunchTemplateFromDbFromLtName() {
	ltName := "taylor-test-deliverhalf"
	var x1 lt.ExtendedGetLaunchTemplateDataOutput

	dbResult := mydb.Db.First(&x1, lt.ExtendedGetLaunchTemplateDataOutput{InstanceName: ltName})
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		log.Logger.Warn("no record found")
		return
	}

	cltInput, err := lt.CreateLaunchTemplateInputFromString(x1.LaunchTemplateDataJsonStr)
	if err != nil {
		log.Logger.Fatalln(err)
	}
	pp.Print(cltInput.LaunchTemplateData.BlockDeviceMappings)
}
