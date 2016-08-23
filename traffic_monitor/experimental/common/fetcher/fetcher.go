package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/handler"
	"github.com/davecheney/gmx"
)

type Fetcher interface {
	Fetch(string, string, uint64, chan<- uint64)
}

type HttpFetcher struct {
	Client  *http.Client
	Headers map[string]string
	Handler handler.Handler
	Counters
}

type Result struct {
	Source string
	Data   []byte
	Error  error
}

type Counters struct {
	Success *gmx.Counter
	Fail    *gmx.Counter
	Pending *gmx.Gauge
}

func (f HttpFetcher) Fetch(id string, url string, pollId uint64, pollFinishedChan chan<- uint64) {
	fmt.Printf("DEBUG poll %v %v fetch start\n", pollId, time.Now())
	req, err := http.NewRequest("GET", url, nil)
	// TODO: change this to use f.Headers. -jse
	req.Header.Set("User-Agent", "traffic_monitor/1.0")
	req.Header.Set("Connection", "keep-alive")
	if f.Pending != nil {
		f.Pending.Inc()
	}
	response, err := f.Client.Do(req)
	if f.Pending != nil {
		f.Pending.Dec()
	}
	defer func() {
		if response != nil && response.Body != nil {
			ioutil.ReadAll(response.Body)
			response.Body.Close()
		}
	}()

	if response != nil {
		if f.Success != nil {
			f.Success.Inc()
		}
		fmt.Printf("DEBUG poll %v %v fetch end\n", pollId, time.Now())
		f.Handler.Handle(id, response.Body, err, pollId, pollFinishedChan)
	} else {
		if f.Fail != nil {
			f.Fail.Inc()
		}
		f.Handler.Handle(id, nil, err, pollId, pollFinishedChan)
	}
}
