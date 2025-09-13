package domain

import "time"

type RegionAlarmInfo struct {
	RegionID      string                  `json:"regionId"`
	RegionType    string                  `json:"regionType"`
	RegionName    string                  `json:"regionName"`
	RegionEngName string                  `json:"regionEngName"`
	LastUpdate    time.Time               `json:"lastUpdate"`
	ActiveAlerts  []RegionActiveAlarmInfo `json:"activeAlerts"`
}

type RegionActiveAlarmInfo struct {
	RegionID   string    `json:"regionId"`
	RegionType string    `json:"regionType"`
	Type       string    `json:"type"`
	LastUpdate time.Time `json:"lastUpdate"`
}
