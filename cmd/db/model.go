package cmd

import (
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"gorm.io/gorm"
)

type IdentityBlob struct {
	gorm.Model
	Doc           imds.InstanceIdentityDocument `gorm:"embedded"`
	B64SNSMessage string
}
