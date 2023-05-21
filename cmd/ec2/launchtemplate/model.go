package cmd

import (
	"gorm.io/gorm"
)

type ExtendedGetLaunchTemplateDataOutput struct {
	gorm.Model
	InstanceId                string
	LaunchTemplateDataJsonStr string `gorm:"embedded"`
}
