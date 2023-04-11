package util

type SystemStatus struct {
	Device string  `json:"device"`
	CPU    float64 `json:"cpu" unit:"%"`
	Memory float64 `json:"memory" unit:"%"`
	Disk   float64 `json:"disk" unit:"%"`
	NetIn  float64 `json:"net_in" unit:"kbps"`
	NetOut float64 `json:"net_out" unit:"kbps"`
	Time   string  `json:"time"`
	CpuId  string  `json:"cpu_id"`
}
