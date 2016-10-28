package peer

import (
	"encoding/json"
	"io"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
)

// Handler handles peer Traffic Monitor data, taking a raw reader, parsing the data, and passing a result object to the ResultChannel. This fulfills the common `Handler` interface.
type Handler struct {
	ResultChannel chan Result
	Notify        int
}

// NewHandler returns a new peer Handler.
func NewHandler() Handler {
	return Handler{ResultChannel: make(chan Result)}
}

// Result contains the data parsed from polling a peer Traffic Monitor.
type Result struct {
	ID           enum.TrafficMonitorName
	Available    bool
	Errors       []error
	PeerStats    Crstates
	PollID       uint64
	PollFinished chan<- uint64
}

// Handle handles a response from a polled Traffic Monitor peer, parsing the data and forwarding it to the ResultChannel.
func (handler Handler) Handle(id string, r io.Reader, reqTime time.Duration, err error, pollID uint64, pollFinished chan<- uint64) {
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
