package poller

import (
	"io/ioutil"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"gopkg.in/fsnotify.v1"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/fetcher"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/handler"
	instr "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/instrumentation"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopswrapper" // TODO move to common
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

type Poller interface {
	Poll()
}

type HttpPoller struct {
	Config        HttpPollerConfig
	ConfigChannel chan HttpPollerConfig
	Fetcher       fetcher.Fetcher
	TickChan      chan uint64
}

type HttpPollerConfig struct {
	Urls     map[string]string
	Interval time.Duration
}

// Creates and returns a new HttpPoller.
// If tick is false, HttpPoller.TickChan() will return nil
func NewHTTP(interval time.Duration, tick bool, httpClient *http.Client, counters fetcher.Counters, fetchHandler handler.Handler) HttpPoller {
	var tickChan chan uint64
	if tick {
		tickChan = make(chan uint64)
	}
	return HttpPoller{
		TickChan:      tickChan,
		ConfigChannel: make(chan HttpPollerConfig),
		Config: HttpPollerConfig{
			Interval: interval,
		},
		Fetcher: fetcher.HttpFetcher{
			Handler:  fetchHandler,
			Client:   httpClient,
			Counters: counters,
		},
	}
}

type FilePoller struct {
	File                string
	ResultChannel       chan interface{}
	NotificationChannel chan int
}

type MonitorConfigPoller struct {
	Session          towrap.ITrafficOpsSession
	SessionChannel   chan towrap.ITrafficOpsSession
	ConfigChannel    chan to.TrafficMonitorConfigMap
	OpsConfigChannel chan handler.OpsConfig
	Interval         time.Duration
	OpsConfig        handler.OpsConfig
}

// Creates and returns a new HttpPoller.
// If tick is false, HttpPoller.TickChan() will return nil
func NewMonitorConfig(interval time.Duration) MonitorConfigPoller {
	return MonitorConfigPoller{
		Interval:         interval,
		SessionChannel:   make(chan towrap.ITrafficOpsSession),
		ConfigChannel:    make(chan to.TrafficMonitorConfigMap),
		OpsConfigChannel: make(chan handler.OpsConfig),
	}
}

func (p MonitorConfigPoller) Poll() {
	tick := time.NewTicker(p.Interval)
	for {
		select {
		case opsConfig := <-p.OpsConfigChannel:
			log.Infof("MonitorConfigPoller: received new opsConfig: %v\n", opsConfig)
			p.OpsConfig = opsConfig
		case session := <-p.SessionChannel:
			log.Infof("MonitorConfigPoller: received new session: %v\n", session)
			p.Session = session
		case <-tick.C:
			if p.Session != nil && p.OpsConfig.CdnName != "" {
				monitorConfig, err := p.Session.TrafficMonitorConfigMap(p.OpsConfig.CdnName)

				if err != nil {
					log.Errorf("MonitorConfigPoller: %s\n %v\n", err, monitorConfig)
				} else {
					log.Infoln("MonitorConfigPoller: fetched monitorConfig")
					p.ConfigChannel <- *monitorConfig
				}
			} else {
				log.Warnln("MonitorConfigPoller: skipping this iteration, Session is nil")
			}
		}
	}
}

var debugPollNum uint64

func (p HttpPoller) Poll() {
	// iterationCount := uint64(0)
	// iterationCount++ // on tick<:
	// case p.TickChan <- iterationCount:
	killChans := map[string]chan<- struct{}{}
	for {
		select {
		case newConfig := <-p.ConfigChannel:
			deletions, additions := diffConfigs(p.Config, newConfig)
			for _, id := range deletions {
				killChan := killChans[id]
				go func() { killChan <- struct{}{} }() // go - we don't want to wait for old polls to die.
				delete(killChans, id)
			}
			for _, info := range additions {
				kill := make(chan struct{})
				killChans[info.ID] = kill
				go pollHttp(info.Interval, info.ID, info.URL, p.Fetcher, kill)
			}
			p.Config = newConfig
		}
	}
}

type HTTPPollInfo struct {
	Interval time.Duration
	ID       string
	URL      string
}

// diffConfigs takes the old and new configs, and returns a list of deleted IDs, and a list of new polls to do
func diffConfigs(old HttpPollerConfig, new HttpPollerConfig) ([]string, []HTTPPollInfo) {
	deletions := []string{}
	additions := []HTTPPollInfo{}

	if old.Interval != new.Interval {
		for id, _ := range old.Urls {
			deletions = append(deletions, id)
		}
		for id, url := range new.Urls {
			additions = append(additions, HTTPPollInfo{Interval: new.Interval, ID: id, URL: url})
		}
		return deletions, additions
	}

	for id, oldUrl := range old.Urls {
		newUrl, newIdExists := new.Urls[id]
		if !newIdExists {
			deletions = append(deletions, id)
		} else if newUrl != oldUrl {
			deletions = append(deletions, id)
			additions = append(additions, HTTPPollInfo{Interval: new.Interval, ID: id, URL: newUrl})
		}
	}

	for id, newUrl := range new.Urls {
		_, oldIdExists := old.Urls[id]
		if !oldIdExists {
			additions = append(additions, HTTPPollInfo{Interval: new.Interval, ID: id, URL: newUrl})
		}
	}

	return deletions, additions
}

func (p FilePoller) Poll() {
	// initial read before watching for changes
	contents, err := ioutil.ReadFile(p.File)

	if err != nil {
		log.Errorf("reading %s: %s\n", p.File, err)
		os.Exit(1) // TODO: this is a little drastic -jse
	} else {
		p.ResultChannel <- contents
	}

	watcher, _ := fsnotify.NewWatcher()
	watcher.Add(p.File)

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				contents, err := ioutil.ReadFile(p.File)

				if err != nil {
					log.Errorf("opening %s: %s\n", p.File, err)
				} else {
					p.ResultChannel <- contents
				}
			}
		case err := <-watcher.Errors:
			log.Errorln(time.Now(), "error:", err)
		}
	}
}

// TODO iterationCount and/or p.TickChan?
func pollHttp(interval time.Duration, id string, url string, fetcher fetcher.Fetcher, die <-chan struct{}) {
	tick := time.NewTicker(interval)
	lastTime := time.Now()
	for {
		select {
		case now := <-tick.C:
			realInterval := now.Sub(lastTime)
			if realInterval > interval+(time.Millisecond*100) {
				instr.TimerFail.Inc()
				log.Infof("Intended Duration: %v Actual Duration: %v\n", interval, realInterval)
			}
			lastTime = time.Now()

			pollId := atomic.AddUint64(&debugPollNum, 1)
			pollFinishedChan := make(chan uint64)
			log.Debugf("poll %v %v start\n", pollId, time.Now())
			go fetcher.Fetch(id, url, pollId, pollFinishedChan) // TODO persist fetcher, with its own die chan?
			<-pollFinishedChan
		case <-die:
			return
		}
	}
}
