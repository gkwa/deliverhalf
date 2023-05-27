package cmd

type Detail struct {
	InstanceID string `json:"instance-id"`
	State      string `json:"state"`
}

type DetailType string

const (
	EC2StateChangeNotification DetailType = "EC2 Instance State-change Notification"
	OtherDetailType            DetailType = "Other Detail Type"
)

type EC2StateChangeEvent struct {
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
