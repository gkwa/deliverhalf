package cmd

import (
	"gorm.io/gorm"
)

type ExtendedEc2Volume struct {
	gorm.Model
	JsonDef string
}

type ExtendedEc2BlockDeviceMapping struct {
	gorm.Model

	InstanceId   string `gorm:"primaryKey;uniqueIndex:idx_unique_mapping"`
	VolumeId     string `gorm:"primaryKey;uniqueIndex:idx_unique_mapping"`
	InstanceName string
	JsonDef      string
}
