/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"time"

	// Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"github.com/spf13/cobra"
	mydb "github.com/taylormonacelli/deliverhalf/cmd/db"
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
	db.AutoMigrate(&mydb.IdentityBlob{})

	// Calculate the time 1 hour ago from now
	oneHourAgo := time.Now().Add(-time.Hour)

	// Get the number of records in the "my_table" table
	var count int64
	if err := db.Model(&mydb.IdentityBlob{}).Count(&count).Error; err != nil {
		log.Logger.Fatalln(err)
	}

	// Find all recent records
	var results []mydb.IdentityBlob
	db.Where("fetchtimestamp >= ?", oneHourAgo).Group("instance_id").Find(&results)
	log.Logger.Debugf("Found %d of %d matching", len(results), count)

	for _, value := range results {
		log.Logger.Trace(value)
	}
}
