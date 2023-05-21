package cmd

import (
	"gorm.io/gorm"
)

type ExtendedInstance struct {
	gorm.Model
	InstanceId string
	JsonDef    string
}
