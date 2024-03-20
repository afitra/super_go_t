package model

import "time"

type LogModel struct {
	Mode   string
	Detail string
	Header string
	Body   string
}

type RequestDetail struct {
	Request_date string `json:"request_date"`
	Request_id   string `json:"request_id"`
	Method       string `json:"method"`
	Url          string `json:"url"`
	Ip           string `json:"ip"`
	Api_type     string `json:"api_type"`
	Level_name   string `json:"level_name"`
	Channel      string `json:"channel"`
	Timezone     string `json:"timezone"`
	Bytes_in     string `json:"bytes_in"`
}

type ResponseDetail struct {
	RequestDetail
	Bytes_out     string        `json:"bytes_out"`
	Error         string        `json:"error"`
	Latency       time.Duration `json:"latency"`
	Latency_human string        `json:"latency_human"`
	Status        int           `json:"status"`
}

type RequestLogger struct {
	RequestID string      `json:"requestID,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
}
