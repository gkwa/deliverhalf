/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
	common "github.com/taylormonacelli/deliverhalf/cmd/common"
	myec2 "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Traceln("fetch called")
		region := "us-west-2"
		// test( region)
		fetchConfigFromS3(region)
	},
}

func test(region string) {
	fetchConfigFromS3(region)
	reloadConfig()
	showSettings()
}

func init() {
	configCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createAWSConfig(region string) aws.Config {
	cfg, err := myec2.CreateConfig(region)
	if err != nil {
		log.Logger.Fatalf("Could not create config %s", err)
	}
	return cfg
}

func createS3Downloader(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg)
}

func logDecodedS3Path(s3path string) {
	log.Logger.Tracef("decoded s3path: %s", s3path)
}

func getUserDirectory() string {
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("USERPROFILE")
	} else {
		dir = os.Getenv("HOME")
	}
	return dir
}

func checkFileExistence(filePath string) {
	isEmpty, err := common.IsFileEmpty(filePath)
	if err != nil {
		log.Logger.Errorln("error:", err)
	} else if isEmpty {
		log.Logger.Tracef("%s exists, but is empty", filePath)
	} else {
		log.Logger.Tracef("%s exists", filePath)
	}
}

func logDownloadedPath(s3path, localTargetPath string) {
	log.Logger.Tracef("%s downloaded to %s", s3path, localTargetPath)
}

func fetchConfigFromS3(region string) {
	cfg := createAWSConfig(region)
	downloader := createS3Downloader(cfg)
	bucketName := common.DecodeBase64String("c3RyZWFtYm94LWRlbGl2ZXJoYWxm")
	config := common.DecodeBase64String("LmRlbGl2ZXJoYWxmLnlhbWw=")
	s3path := fmt.Sprintf("s3://%s/%s", bucketName, config)
	logDecodedS3Path(s3path)

	dir := getUserDirectory()
	filename := config
	localTargetPath := filepath.Join(dir, filename)

	checkFileExistence(localTargetPath)
	err := common.DownloadFileFromS3(downloader, bucketName, filename, localTargetPath)
	common.HandleDownloadError(err)

	logDownloadedPath(s3path, localTargetPath)
}
