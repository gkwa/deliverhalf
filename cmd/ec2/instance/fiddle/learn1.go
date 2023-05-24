package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
	mydb "github.com/taylormonacelli/deliverhalf/cmd/db"
	instance "github.com/taylormonacelli/deliverhalf/cmd/ec2/instance"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
	"gorm.io/gorm"
)

// learn1Cmd represents the learn1 command
var learn1Cmd = &cobra.Command{
	Use:   "learn1",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Trace("learn1 called")
		learn1()
	},
}

func init() {
	fiddleCmd.AddCommand(learn1Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// learn1Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// learn1Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func learn1() {
	var myExtendedInstance instance.ExtendedInstanceDetail
	queryResult := mydb.Db.Last(&myExtendedInstance)
	if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
		log.Logger.Warn("no instances found")
		return
	}
	var inst types.Instance
	err := json.Unmarshal([]byte(myExtendedInstance.JsonDef), &inst)
	if err != nil {
		log.Logger.Fatal(err)
	}
	name := instance.GetTagValue(&inst.Tags, "Name")
	fmt.Printf("instance name: %s\n", name)
}
