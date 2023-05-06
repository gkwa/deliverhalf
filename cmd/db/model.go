package cmd

import (
	imds "github.com/taylormonacelli/deliverhalf/cmd/ec2/imds"
	"gorm.io/gorm"
)

type IdentityBlob struct {
	gorm.Model
	Doc                     imds.ExtendedInstanceIdentityDocument `gorm:"embedded"`
	B64SNSMessage           string
	B64SNSMessageCompressed string
}
