package cmd

import (
	"gorm.io/gorm"
)

type ExtendedEc2Volume struct {
	gorm.Model
	JsonDef string
}
