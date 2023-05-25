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
	VolumeId     string
	InstanceId   string
	InstanceName string
	JsonDef      string
}
