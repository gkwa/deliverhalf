package cmd

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type MultiString []string

// "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
type ExtendedInstanceIdentityDocument struct {
	gorm.Model
	AccountId               string
	Architecture            string
	AvailabilityZone        string
	FetchTimestamp          time.Time
	ImageId                 string
	InstanceId              string
	InstanceType            string
	KernelId                string
	PendingTime             string
	PrivateIp               string
	RamdiskId               string
	Region                  string
	Version                 string
	BillingProducts         MultiString `gorm:"type:text"` // []string
	DevpayProductCodes      MultiString `gorm:"type:text"` // []string
	MarketplaceProductCodes MultiString `gorm:"type:text"` // []string
}

type IdentityBlob struct {
	gorm.Model
	Doc                     ExtendedInstanceIdentityDocument `gorm:"embedded"`
	B64SNSMessage           string
	B64SNSMessageCompressed string
}

func (s *MultiString) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New("failed to scan multistring field - source is not a string")
	}
	*s = strings.Split(str, ",")
	return nil
}

func (s MultiString) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	return strings.Join(s, ","), nil
}
