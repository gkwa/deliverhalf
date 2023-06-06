package cmd

import (
	"gorm.io/gorm"
)

// Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details

type ExtendedSqsReceiveMessageOutput struct {
	gorm.Model
	JsonDef string
}

type MessageDetail struct {
	Version    string   `json:"version"`
	ID         string   `json:"id"`
	DetailType string   `json:"detail-type"`
	Source     string   `json:"source"`
	Account    string   `json:"account"`
	Time       string   `json:"time"`
	Region     string   `json:"region"`
	Resources  []string `json:"resources"`
	Detail     Detail   `json:"detail"`
}

type NotificationMessage struct {
	Type             string        `json:"Type"`
	MessageID        string        `json:"MessageId"`
	TopicArn         string        `json:"TopicArn"`
	Message          string        `json:"Message"`
	Timestamp        string        `json:"Timestamp"`
	SignatureVersion string        `json:"SignatureVersion"`
	Signature        string        `json:"Signature"`
	SigningCertURL   string        `json:"SigningCertURL"`
	UnsubscribeURL   string        `json:"UnsubscribeURL"`
	MessageDetail    MessageDetail `json:"-"`
}

type Detail struct {
	InstanceID string `json:"instance-id"`
	State      string `json:"state"`
}
