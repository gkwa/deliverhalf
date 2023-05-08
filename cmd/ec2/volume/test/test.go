/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	volume "github.com/taylormonacelli/deliverhalf/cmd/ec2/volume"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"

	"gorm.io/gorm"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		test1()
	},
}

func init() {
	volume.VolumeCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func test1() {
	jsonStr := `{
	"Attachments": [
	  {
		"AttachTime": "2023-05-07T10:15:23Z",
		"DeleteOnTermination": true,
		"Device": "/dev/sdh",
		"InstanceId": "i-0a1b2c3d4e5f67890",
		"State": "attached",
		"VolumeId": "vol-0123456789abcdef"
	  }
	],
	"AvailabilityZone": "us-east-1a",
	"CreateTime": "2023-05-07T10:14:49.123Z",
	"Encrypted": true,
	"FastRestored": false,
	"Iops": 100,
	"KmsKeyId": "arn:aws:kms:us-east-1:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab",
	"MultiAttachEnabled": false,
	"OutpostArn": null,
	"Size": 100,
	"SnapshotId": "snap-0123456789abcdef",
	"State": "in-use",
	"Tags": [
	  {
		"Key": "Name",
		"Value": "test-volume"
	  }
	],
	"Throughput": 256,
	"VolumeId": "vol-0123456789abcdef",
	"VolumeType": "gp3"
  }
  `

	var myvol volume.ExtendedEc2Volume
	err := json.Unmarshal([]byte(jsonStr), &myvol)
	if err != nil {
		log.Logger.Fatalf("tried to unmarshal %s to a %T but got error %s",
			jsonStr, myvol, err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&volume.ExtendedEc2Volume{})
	db.Create(&volume.ExtendedEc2Volume{JsonDef: jsonStr})

	var extVol volume.ExtendedEc2Volume
	myvol2 := make(map[string]interface{})
	db.Last(&extVol)
	fmt.Print(extVol.JsonDef)
	err = json.Unmarshal([]byte(extVol.JsonDef), &myvol2)
	if err != nil {
		log.Logger.Fatalf("tried to unmarshal '%s' to a %T but got error %s",
			extVol.JsonDef, myvol2, err)
	}

	attachments := myvol2["Attachments"]
	pp.Print(attachments)
}
