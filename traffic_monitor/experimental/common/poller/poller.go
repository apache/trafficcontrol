package poller

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/fsnotify.v1"

	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/fetcher"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/handler"
	instr "github.com/Comcast/traffic_control/traffic_monitor/experimental/common/instrumentation"
	traffic_ops "github.com/Comcast/traffic_control/traffic_ops/client"
)

type Poller interface {
	Poll()
}

type HttpPoller struct {
	Config        HttpPollerConfig
	ConfigChannel chan HttpPollerConfig
	Fetcher       fetcher.Fetcher
}

type HttpPollerConfig struct {
	Urls     map[string]string
	Interval time.Duration
}

type FilePoller struct {
	File                string
	ResultChannel       chan interface{}
	NotificationChannel chan int
}

type MonitorConfigPoller struct {
	Session          *traffic_ops.Session
	SessionChannel   chan *traffic_ops.Session
	ConfigChannel    chan traffic_ops.TrafficMonitorConfigMap
	OpsConfigChannel chan handler.OpsConfig
	Interval         time.Duration
	OpsConfig        handler.OpsConfig
}

func (p MonitorConfigPoller) Poll() {
	tick := time.NewTicker(p.Interval)

	for {
		select {
		case opsConfig := <-p.OpsConfigChannel:
			fmt.Println("MonitorConfigPoller: received new opsConfig: %v", opsConfig)
			p.OpsConfig = opsConfig
		case session := <-p.SessionChannel:
			fmt.Println("MonitorConfigPoller: received new session: %v", session)
			p.Session = session
		case <-tick.C:
			if p.Session != nil && p.OpsConfig.CdnName != "" {
				monitorConfig, err := p.Session.TrafficMonitorConfigMap(p.OpsConfig.CdnName)

				if err != nil {
					fmt.Printf("MonitorConfigPoller Error: %s\n %v", err, monitorConfig)
				} else {
					//fmt.Printf("MonitorConfigPoller: fetched monitorConfig\n")
					p.ConfigChannel <- *monitorConfig
				}
			} else {
				fmt.Println("MonitorConfigPoller: skipping this iteration, Session is nil")
			}
		}
	}
}

func (p HttpPoller) Poll() {
	tick := time.NewTicker(p.Config.Interval)
	last_time := time.Now()

	for {
		select {
		case config := <-p.ConfigChannel:
			p.Config = config // TODO: reset the ticker using the interval supplied in config. -jse
		case curr_time := <-tick.C:
			if int64(curr_time.Sub(last_time)) > int64(float64(p.Config.Interval)*1.01) {
				instr.TimerFail.Inc()
				fmt.Println("Intended Duration:", p.Config.Interval, "Actual Duration", curr_time.Sub(last_time))
			}
			last_time = curr_time
			if p.Config.Urls != nil {
				for id, url := range p.Config.Urls {
					go p.Fetcher.Fetch(id, url)
				}
			}
		}
	}
}

func (p FilePoller) Poll() {
	// initial read before watching for changes
	contents, err := ioutil.ReadFile(p.File)

	if err != nil {
		fmt.Printf("Error reading %s: %s\n", p.File, err)
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
					fmt.Printf("Error opening %s: %s\n", p.File, err)
				} else {
					p.ResultChannel <- contents
				}
			}
		case err := <-watcher.Errors:
			fmt.Println(time.Now(), "error:", err)
		}
	}
}
