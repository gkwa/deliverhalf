/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/taylormonacelli/deliverhalf/cmd"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	imds "github.com/taylormonacelli/deliverhalf/cmd/ec2/imds"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"

	"github.com/spf13/cobra"
)

// metaCmd represents the meta command
var metaCmd = &cobra.Command{
	Use:   "meta",
	Args:  cobra.OnlyValidArgs,
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Traceln("meta called")
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		os.Exit(1)
		return nil
	},
}

func init() {
	cmd.RootCmd.AddCommand(metaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func genPathToMetaJson() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "meta.json"), nil
}

func GetIdentityDocFromFile() (imds.ExtendedInstanceIdentityDocument, error) {
	metaPath, err := genPathToMetaJson()
	if err != nil {
		log.Logger.Fatalf("error generating path to file storing jsonblob: %s", err)
	}

	var doc imds.ExtendedInstanceIdentityDocument

	if !common.FileExists(metaPath) {
		log.Logger.Fatalf("file %s doesn't exist, but i expect to use it to unmarshal an %T", metaPath, doc)
	}

	jsonBlob, err := os.ReadFile(metaPath)
	if err != nil {
		log.Logger.Fatalf("failed to read from file %s: %s", metaPath, err)
	}
	jsonStr := string(jsonBlob)

	doc, err = GetIdentityDocFromStr(jsonStr)
	if err != nil {
		log.Logger.Fatalf("failed to unmarshal %s into an %T: %s", jsonStr, doc, err)
	}

	return doc, nil
}

func GetIdentityDocFromStr(str string) (imds.ExtendedInstanceIdentityDocument, error) {
	// Check if the JSON string is valid
	if json.Valid([]byte(str)) {
		log.Logger.Tracef("PASS: checking that json string is valid: %s", str)
	} else {
		log.Logger.Fatalf("JSON string is invalid: %s", str)
	}

	var doc imds.ExtendedInstanceIdentityDocument
	err := json.Unmarshal([]byte(str), &doc)
	if err != nil {
		log.Logger.Fatalf("can't unmarshal %s into %T: %s", str, doc, err)
	}

	return doc, nil
}
