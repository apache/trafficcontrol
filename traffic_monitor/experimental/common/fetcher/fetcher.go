package fetcher

import (
	"io/ioutil"
	"net/http"

	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/handler"
	"github.com/davecheney/gmx"
)

type Fetcher interface {
	Fetch(string, string)
}

type HttpFetcher struct {
	Client  http.Client
	Headers map[string]string
	Handler handler.Handler
	Success *gmx.Counter
	Fail    *gmx.Counter
	Pending *gmx.Gauge
}

type Result struct {
	Source string
	Data   []byte
	Error  error
}

func (f HttpFetcher) Fetch(id string, url string) {
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
		f.Handler.Handle(id, response.Body, err)
	} else {
		if f.Fail != nil {
			f.Fail.Inc()
		}
		f.Handler.Handle(id, nil, err)
	}
}
