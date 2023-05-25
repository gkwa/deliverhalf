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
	InstanceId   string `gorm:"primaryKey"`
	VolumeId     string `gorm:"primaryKey"`
	InstanceName string
	JsonDef      string
}
