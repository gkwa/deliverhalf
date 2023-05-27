package cmd

import "gorm.io/gorm"

// Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details

type ExtendedSqsReceiveMessageOutput struct {
	gorm.Model
	JsonDef string
}
