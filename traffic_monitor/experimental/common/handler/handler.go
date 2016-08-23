package handler

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	NOTIFY_NEVER = iota
	NOTIFY_CHANGE
	NOTIFY_ALWAYS
)

type Handler interface {
	Handle(string, io.Reader, error, uint64, chan<- uint64)
}

type OpsConfigFileHandler struct {
	Content          interface{}
	ResultChannel    chan interface{}
	OpsConfigChannel chan OpsConfig
}

type OpsConfig struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Url          string `json:"url"`
	Insecure     bool   `json:"insecure"`
	CdnName      string `json:"cdnName"`
	HttpListener string `json:"httpListener"`
}

func (handler OpsConfigFileHandler) Listen() {
	for {
		result := <-handler.ResultChannel
		var toc OpsConfig

		err := json.Unmarshal(result.([]byte), &toc)

		if err != nil {
			fmt.Printf("Error unmarshalling JSON: %s\n", err)
		} else {
			handler.OpsConfigChannel <- toc
		}
	}
}
