package peer

import (
	"encoding/json"
	"io"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
)

type Handler struct {
	ResultChannel chan Result
	Notify        int
}

func NewHandler() Handler {
	return Handler{ResultChannel: make(chan Result)}
}

type Result struct {
	ID           enum.TrafficMonitorName
	Available    bool
	Errors       []error
	PeerStats    Crstates
	PollID       uint64
	PollFinished chan<- uint64
}

func (handler Handler) Handle(id string, r io.Reader, err error, pollID uint64, pollFinished chan<- uint64) {
	result := Result{
		ID:           enum.TrafficMonitorName(id),
		Available:    false,
		Errors:       []error{},
		PollID:       pollID,
		PollFinished: pollFinished,
	}

	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	if r != nil {
		dec := json.NewDecoder(r)

		if err := dec.Decode(&result.PeerStats); err == io.EOF {
			result.Available = true
		} else if err != nil {
			result.Errors = append(result.Errors, err)
		}
	}

	handler.ResultChannel <- result
}
