package peer

import (
	"encoding/json"
	"io"
)

type Handler struct {
	ResultChannel chan Result
	Notify        int
}

type Result struct {
	Id        string
	Available bool
	Errors    []error
	PeerStats Crstates
}

const (
	NOTIFY_NEVER = iota
	NOTIFY_CHANGE
	NOTIFY_ALWAYS
)

func (handler Handler) Handle(id string, r io.Reader, err error) {
	result := Result{
		Id:        id,
		Available: false,
		Errors:    make([]error, 0, 0),
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
