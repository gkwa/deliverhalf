/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"

	mydb "github.com/taylormonacelli/deliverhalf/cmd/db"
	instance "github.com/taylormonacelli/deliverhalf/cmd/ec2/instance"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// learn2Cmd represents the learn2 command
var learn2Cmd = &cobra.Command{
	Use:   "learn2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Traceln("learn2 called")
		learn2()
		cleanOldRecords()
	},
}

func init() {
	fiddleCmd.AddCommand(learn2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// learn2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// learn2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func learn2() {
	var extendedInstances []instance.ExtendedInstanceDetail

	query := `
		JOIN (
			SELECT MAX(created_at) AS max_created_at
			FROM extended_instance_details
			GROUP BY name
		) AS subquery ON extended_instance_details.created_at = subquery.max_created_at
	`
	mydb.Db.Table("extended_instance_details").
		Select("extended_instance_details.*").
		Joins(query).
		Find(&extendedInstances)
	for _, value := range extendedInstances {
		fmt.Printf("%s %s\n", value.CreatedAt.Format("2006-01-02 15:04:05"), value.Name)
	}
	fmt.Printf("%d records\n", len(extendedInstances))
}

func cleanOldRecords() {
	// oneYearAgo := time.Now().AddDate(-1, 0, 0)
	// since := time.Now().AddDate(-1, 0, 0)
	// oneDayAgo := time.Now().AddDate(0, 0, -1)
	// since := time.Now().AddDate(0, 0, -1)
	tenMinutesAgo := time.Now().Add(-5 * time.Minute)
	// since := time.Now().Add(-10 * time.Minute)
	// ninetyDaysAgo := time.Now().AddDate(0, 0, -90)
	// since := time.Now().AddDate(0, 0, -90)
	since := tenMinutesAgo
	result := mydb.Db.Where("created_at < ?", since).Delete(&instance.ExtendedInstanceDetail{})

	if result.Error != nil {
		log.Logger.Errorln("Error occurred:", result.Error)
		return
	}

	recordsDeleted := decimal.NewFromInt(result.RowsAffected)
	log.Logger.Debugf("%s records deleted from ExtendedInstanceDetail because they are old", recordsDeleted.StringFixed(2))
}
