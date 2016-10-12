package data

import "time"

type HttpResult struct {
	RequestTime time.Time `json:"requestTime"`
	Host        string    `json:"host"`
	LatencyUsec int64     `json:"latency"`
	Err         error     `json:"error"`
	Status      int       `json:"httpStatus"`
}

