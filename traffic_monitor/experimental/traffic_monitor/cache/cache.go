package cache

import (
	"encoding/json"
	"io"
	"time"
)

type Handler struct {
	ResultChannel chan Result
	Notify        int
}

func NewHandler() Handler {
	return Handler{ResultChannel: make(chan Result)}
}

type Result struct {
	Id        string
	Available bool
	Errors    []error
	Astats    Astats
	Time      time.Time
	Vitals    Vitals
}

type Vitals struct {
	LoadAvg    float64
	BytesOut   int64
	BytesIn    int64
	KbpsOut    int64
	MaxKbpsOut int64
}

type Stat struct {
	Time  int64       `json:"time"`
	Value interface{} `json:"value"`
}

type Stats struct {
	Caches map[string]map[string][]Stat `json:"caches"`
}

const (
	NOTIFY_NEVER = iota
	NOTIFY_CHANGE
	NOTIFY_ALWAYS
)

func StatsMarshall(statHistory map[string][]Result, historyCount int) ([]byte, error) {
	var stats Stats

	stats.Caches = map[string]map[string][]Stat{}

	count := 1

	for id, history := range statHistory {
		for _, result := range history {
			for stat, value := range result.Astats.Ats {
				s := Stat{
					Time:  result.Time.UnixNano() / 1000000,
					Value: value,
				}

				_, exists := stats.Caches[id]

				if !exists {
					stats.Caches[id] = map[string][]Stat{}
				}

				stats.Caches[id][stat] = append(stats.Caches[id][stat], s)
			}

			if historyCount > 0 && count == historyCount {
				break
			}

			count++
		}
	}

	return json.Marshal(stats)
}

func (handler Handler) Handle(id string, r io.Reader, err error) {
	result := Result{
		Id:        id,
		Available: false,
		Errors:    []error{},
		Time:      time.Now(),
	}

	if err != nil {
		result.Errors = append(result.Errors, err)
	}

	if r != nil {
		dec := json.NewDecoder(r)
		err := dec.Decode(&result.Astats)

		if err != nil {
			result.Errors = append(result.Errors, err)
		} else {
			result.Available = true
		}
	}

	handler.ResultChannel <- result
}
