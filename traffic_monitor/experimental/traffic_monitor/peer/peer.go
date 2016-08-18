package peer

import (
	"encoding/json"
	"io"
)

type Handler struct {
	ResultChannel chan Result
	Notify        int
}

func NewHandler() Handler {
	return Handler{ResultChannel: make(chan Result)}
}

type Result struct {
	Id           string
	Available    bool
	Errors       []error
	PeerStats    Crstates
	PollID       uint64
	PollFinished chan<- uint64
}

const (
	NOTIFY_NEVER = iota
	NOTIFY_CHANGE
	NOTIFY_ALWAYS
)

func (handler Handler) Handle(id string, r io.Reader, err error, pollId uint64, pollFinished chan<- uint64) {
	result := Result{
		Id:           id,
		Available:    false,
		Errors:       []error{},
		PollID:       pollId,
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
