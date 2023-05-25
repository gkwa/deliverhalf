package cmd

import (
	"gorm.io/gorm"
)

type ExtendedGetLaunchTemplateDataOutput struct {
	gorm.Model
	InstanceName              string // use instance name if we create lt from instance
	InstanceId                string
	Region                    string
	LaunchTemplateDataJsonStr string `gorm:"embedded"`
}
