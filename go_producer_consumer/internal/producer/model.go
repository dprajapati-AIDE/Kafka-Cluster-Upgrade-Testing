package producer

type DeviceMessageModel struct {
	Timestamp   string `json:"timestamp"`
	DeviceID    string `json:"device_id"`
	DeviceIP    string `json:"device_ip"`
	DeviceType  string `json:"device_type"`
	DeviceModel string `json:"device_model"`
	Vendor      string `json:"vendor"`
	Message     string `json:"message"`
}
