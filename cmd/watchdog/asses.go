/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	// Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	mydb "github.com/taylormonacelli/deliverhalf/cmd/db"
	myimds "github.com/taylormonacelli/deliverhalf/cmd/ec2/imds"
	lt "github.com/taylormonacelli/deliverhalf/cmd/ec2/launchtemplate"
	volume "github.com/taylormonacelli/deliverhalf/cmd/ec2/volume"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// assesCmd represents the asses command
var assesCmd = &cobra.Command{
	Use:   "asses",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Trace("asses called")
		asses()
	},
}

func init() {
	WatchdogCmd.AddCommand(assesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// assesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// assesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func asses() {
	db, err := mydb.OpenDB("test.db")
	if err != nil {
		log.Logger.Fatalf("failed to connect to database: %v", err)
	}

	// Auto Migrate
	db.AutoMigrate(&myimds.IdentityBlob{})
	db.AutoMigrate(&lt.ExtendedGetLaunchTemplateDataOutput{})

	since := time.Now().Add(-time.Hour)

	var count int64
	if err := db.Model(&myimds.IdentityBlob{}).Count(&count).Error; err != nil {
		log.Logger.Fatalln(err)
	}

	var identityDocs []myimds.IdentityBlob
	query := "fetch_timestamp >= ?"
	db.Where(query, since).Group("instance_id").Find(&identityDocs)
	log.Logger.Debugf("found %d matching of %d identity documents", len(identityDocs), count)

	for _, value := range identityDocs {
		log.Logger.Trace(value)
	}
	if len(identityDocs) < 1 {
		log.Logger.Warnln("no identity documents found")
		return
	}

	doc := identityDocs[0].Doc
	instanceId := doc.InstanceId
	region := doc.Region
	tmpFname, _ := filepath.Abs(filepath.Join("data", fmt.Sprintf("%s-LaunchTemplate.json", instanceId)))
	var volumes []types.Volume

	template, err := lt.GenLaunchTemplateFromInstanceId(region, instanceId, tmpFname)
	if err != nil {
		log.Logger.Error(err)
		return
	}

	if template.LaunchTemplateData == nil {
		log.Logger.Warnf("instance %s in region %s no longer exists, can't fetch launch template from it", instanceId, region)
		return
	}

	blockDevices := template.LaunchTemplateData.BlockDeviceMappings
	pp.Print(blockDevices)
	volume.GetVolumesFromInstanceIdentityDoc(doc, &volumes)
}
