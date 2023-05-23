package cmd

import (
	"gorm.io/gorm"
)

type ExtendedGetLaunchTemplateDataOutput struct {
	gorm.Model
	InstanceId                string
	Region                    string
	LaunchTemplateDataJsonStr string `gorm:"embedded"`
}
