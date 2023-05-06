package cmd

import (
	"database/sql/driver"
	"errors"
	"strings"

	"gorm.io/gorm"
	// "gorm.io/driver/sqlite" // Sqlite driver based on GGO
	// Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
)

type MultiString []string

// "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
type ExtendedInstanceIdentityDocument struct {
	gorm.Model
	AccountId               string
	Architecture            string
	AvailabilityZone        string
	Epochtime               int64
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
