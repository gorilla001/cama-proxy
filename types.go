package main

type AlertMessage struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
}

// Alert is a single alert.
type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    string            `json:"startsAt,omitempty"`
	EndsAt      string            `json:"EndsAt,omitempty"`
}

type Notification struct {
	ID        string `json:"IDENTIFIER"`
	Channel   string `json:"CHANNEL"`
	FirstTime string `json:"FIRST_TIME"`
	LastTime  string `json:"LAST_TIME"`
	Recover   int    `json:"RECOVER"`
	Merger    int    `json:"MERGER"`
	Node      string `json:"NODE"`
	NodeAlias string `json:"NODEALIAS"`
	ServerNo  string `json:"SERVER_NO"`
	EventDesc string `json:"EVENT_DESC"`
	Level     int    `json:"LEVEL"`
}

type Event struct {
	ID             string `json:"Event_ID"`
	Type           string `json:"Type"`
	Level          string `json:"Level"`
	SourceID       string `json:"Src_Sys_ID"`
	SourceType     string `json:"Src_Type"`
	NodeName       string `json:"Node_Name"`
	NodeIP         string `json:"Node_IP"`
	FirstTime      string `json:"First_Time"`
	LastTime       string `json:"Last_Time"`
	AlertTimes     int    `json:"Alert_Times"`
	AlertKeyType   string `json:"Alert_Key_Type"`
	AlertKey       string `json:"Alert_Key"`
	AlertValue     string `json:"Alert_Value"`
	AlertThreshold string `json:"Alert_Threshold"`
	AlertMsg       string `json:"Alert_Msg"`
}
