/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	mydb "github.com/taylormonacelli/deliverhalf/cmd/db"
	lt "github.com/taylormonacelli/deliverhalf/cmd/ec2/launchtemplate"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
	"gorm.io/gorm"
)

// learn3Cmd represents the learn3 command
var learn3Cmd = &cobra.Command{
	Use:   "learn3",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("learn3 called")
		// testGormCheckWhetherRecordExists()
		testGormCheckWhetherRecordExistsShort()
	},
}

func init() {
	lt.LaunchtemplateCmd.AddCommand(learn3Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// learn3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// learn3Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func testGormCheckWhetherRecordExists() {
	instanceID := "i-0d681411429a045cc"
	var result lt.ExtendedGetLaunchTemplateDataOutput
	err := mydb.Db.First(&result, lt.ExtendedGetLaunchTemplateDataOutput{InstanceId: instanceID}).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Handle the case when the record is not found
			// Your code here
			fmt.Printf("%s not found instance\n", instanceID)
		} else {
			// Handle other types of errors
			// Your code here
			fmt.Println(fmt.Errorf("unexpected error"))
		}
	} else {
		// Record found, use the result variable
		// Your code here
		fmt.Printf("found instance %s\n", instanceID)
	}
}

func testGormCheckWhetherRecordExistsShort() {
	instanceID := "i-0d681411429a045cc1"
	var record lt.ExtendedGetLaunchTemplateDataOutput
	dbRresult := mydb.Db.First(&record, lt.ExtendedGetLaunchTemplateDataOutput{InstanceId: instanceID})

	if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		log.Logger.Warnf("%s not found instance\n", instanceID)
		return
	}
	log.Logger.Printf("found instance %s\n", instanceID)
}
