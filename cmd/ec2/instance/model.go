package cmd

import (
	"gorm.io/gorm"
)

type ExtendedInstanceDetail struct {
	gorm.Model
	InstanceId string
	JsonDef    string
	Region     string
	Name       string
}
