/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	traffic_ops "github.com/Comcast/traffic_control/traffic_ops/client"
	log "github.com/cihub/seelog"
	influx "github.com/influxdb/influxdb/client/v2"
)

const (
	// FATAL will exit after printing error
	FATAL = iota
	// ERROR will just keep going, print error
	ERROR = iota
)

const (
	defaultPollingInterval             = 10
	defaultDailySummaryPollingInterval = 60
	defaultConfigInterval              = 300
	defaultPublishingInterval          = 10
)

// StartupConfig contains all fields necessary to create an InfluxDB session.
type StartupConfig struct {
	ToUser                      string                  `json:"toUser"`
	ToPasswd                    string                  `json:"toPasswd"`
	ToURL                       string                  `json:"toUrl"`
	InfluxUser                  string                  `json:"influxUser"`
	InfluxPassword              string                  `json:"influxPassword"`
	PollingInterval             int                     `json:"pollingInterval"`
	DailySummaryPollingInterval int                     `json:"dailySummaryPollingInterval"`
	PublishingInterval          int                     `json:"publishingInterval"`
	ConfigInterval              int                     `json:"configInterval"`
	StatusToMon                 string                  `json:"statusToMon"`
	SeelogConfig                string                  `json:"seelogConfig"`
	CacheRetentionPolicy        string                  `json:"cacheRetentionPolicy"`
	DsRetentionPolicy           string                  `json:"dsRetentionPolicy"`
	DailySummaryRetentionPolicy string                  `json:"dailySummaryRetentionPolicy"`
	BpsChan                     chan influx.BatchPoints `json:"-"`
}

// RunningConfig contains information about current InfluxDB connections.
type RunningConfig struct {
	HealthUrls    map[string]map[string]string // they 1st map key is CDN_name, the second is DsStats or CacheStats
	CacheGroupMap map[string]string            // map hostName to cacheGroup
	InfluxDBProps []struct {
		Fqdn string
		Port int64
	}
	LastSummaryTime time.Time
}

//Timers struct containts all the timers
type Timers struct {
	Poll         <-chan time.Time
	DailySummary <-chan time.Time
	Publish      <-chan time.Time
	Config       <-chan time.Time
}

func main() {
	var Bps map[string]*influx.BatchPoints
	var config StartupConfig
	var err error
	var tickers Timers

	configFile := flag.String("cfg", "", "The config file")
	flag.Parse()

	config, err = loadStartupConfig(*configFile, config)

	if err != nil {
		errHndlr(err, FATAL)
	}

	Bps = make(map[string]*influx.BatchPoints)
	config.BpsChan = make(chan influx.BatchPoints)

	defer log.Flush()

	configChan := make(chan RunningConfig)
	go getToData(config, true, configChan)
	runningConfig := <-configChan

	tickers = setTimers(config)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	hupChan := make(chan os.Signal, 1)
	signal.Notify(hupChan, syscall.SIGHUP)

	for {
		select {
		case <-hupChan:
			log.Info("HUP Received - reloading config")
			newConfig, err := loadStartupConfig(*configFile, config)

			if err != nil {
				errHndlr(err, ERROR)
			} else {
				config = newConfig
				tickers = setTimers(config)
			}
		case <-termChan:
			log.Info("Shutdown Request Received - Sending stored metrics then quitting")
			for _, val := range Bps {
				sendMetrics(config, runningConfig, *val, false)
			}
			os.Exit(0)
		case <-tickers.Publish:
			for key, val := range Bps {
				go sendMetrics(config, runningConfig, *val, true)
				delete(Bps, key)
			}
		case runningConfig = <-configChan:
		case <-tickers.Config:
			go getToData(config, false, configChan)
		case <-tickers.Poll:
			for cdnName, urls := range runningConfig.HealthUrls {
				for _, url := range urls {
					log.Debug(cdnName, " -> ", url)
					go calcMetrics(cdnName, url, runningConfig.CacheGroupMap, config, runningConfig)
				}
			}
		case now := <-tickers.DailySummary:
			go calcDailySummary(now, config, runningConfig)
		case batchPoints := <-config.BpsChan:
			log.Debug("Received ", len(batchPoints.Points()), " stats")
			key := fmt.Sprintf("%s%s", batchPoints.Database(), batchPoints.RetentionPolicy())
			bp, ok := Bps[key]
			if ok {
				b := *bp
				for _, p := range batchPoints.Points() {
					b.AddPoint(p)
				}
				log.Debug("Aggregating ", len(b.Points()), " stats to ", key)
			} else {
				Bps[key] = &batchPoints
				log.Debug("Created ", key)
			}
		}
	}
}

func setTimers(config StartupConfig) Timers {
	var timers Timers

	<-time.NewTimer(time.Now().Truncate(time.Duration(config.PollingInterval) * time.Second).Add(time.Duration(config.PollingInterval) * time.Second).Sub(time.Now())).C
	timers.Poll = time.Tick(time.Duration(config.PollingInterval) * time.Second)
	timers.DailySummary = time.Tick(time.Duration(config.DailySummaryPollingInterval) * time.Second)
	timers.Publish = time.Tick(time.Duration(config.PublishingInterval) * time.Second)
	timers.Config = time.Tick(time.Duration(config.ConfigInterval) * time.Second)

	return timers
}

func loadStartupConfig(configFile string, oldConfig StartupConfig) (StartupConfig, error) {
	var config StartupConfig

	file, err := os.Open(configFile)

	if err != nil {
		return config, err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&config)

	if err != nil {
		return config, err
	}

	config.BpsChan = oldConfig.BpsChan

	if config.PollingInterval == 0 {
		config.PollingInterval = defaultPollingInterval
	}
	if config.DailySummaryPollingInterval == 0 {
		config.DailySummaryPollingInterval = defaultDailySummaryPollingInterval
	}
	if config.PublishingInterval == 0 {
		config.PublishingInterval = defaultPublishingInterval
	}
	if config.ConfigInterval == 0 {
		config.ConfigInterval = defaultConfigInterval
	}

	logger, err := log.LoggerFromConfigAsFile(config.SeelogConfig)
	if err != nil {
		errHndlr(fmt.Errorf("error reading Seelog config %s", config.SeelogConfig), ERROR)
	} else {
		log.ReplaceLogger(logger)
		log.Info("Replaced logger, see log file according to", config.SeelogConfig)
	}

	return config, nil
}

func calcDailySummary(now time.Time, config StartupConfig, runningConfig RunningConfig) {
	log.Infof("lastSummaryTime is %v", runningConfig.LastSummaryTime)
	if runningConfig.LastSummaryTime.Day() != now.Day() {
		startTime := now.Truncate(24 * time.Hour).Add(-24 * time.Hour)
		endTime := startTime.Add(24 * time.Hour)
		log.Info("Summarizing from ", startTime, " (", startTime.Unix(), ") to ", endTime, " (", endTime.Unix(), ")")

		// influx connection
		influxClient, err := influxConnect(config, runningConfig)
		if err != nil {
			log.Error("Could not connect to InfluxDb to get daily summary stats!!")
			errHndlr(err, ERROR)
			return
		}

		//create influxdb query
		q := fmt.Sprintf("SELECT sum(value)/6 FROM bandwidth where time > '%s' and time < '%s' group by time(60s), cdn fill(0)", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
		log.Infof(q)
		res, err := queryDB(influxClient, q, "cache_stats")
		if err != nil {
			errHndlr(err, ERROR)
			return
		}

		bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
			Database:        "daily_stats",
			Precision:       "s",
			RetentionPolicy: config.DailySummaryRetentionPolicy,
		})
		for _, row := range res[0].Series {
			prevtime := startTime
			max := float64(0)
			bytesServed := float64(0)
			cdn := row.Tags["cdn"]
			for _, record := range row.Values {
				kbps, err := record[1].(json.Number).Float64()
				if err != nil {
					errHndlr(err, ERROR)
					continue
				}
				sampleTime, err := time.Parse(time.RFC3339, record[0].(string))
				if err != nil {
					errHndlr(err, ERROR)
					continue
				}
				max = FloatMax(max, kbps)
				duration := sampleTime.Unix() - prevtime.Unix()
				bytesServed += float64(duration) * kbps / 8
				prevtime = sampleTime
			}
			maxGbps := max / 1000000
			bytesServedTb := bytesServed / 1000000000
			log.Infof("max gbps for cdn %v = %v", cdn, maxGbps)
			log.Infof("Tbytes served for cdn %v = %v", cdn, bytesServedTb)

			//write daily_maxgbps in traffic_ops
			var statsSummary traffic_ops.StatsSummary
			statsSummary.CdnName = cdn
			statsSummary.DeliveryService = "all"
			statsSummary.StatName = "daily_maxgbps"
			statsSummary.StatValue = strconv.FormatFloat(maxGbps, 'f', 2, 64)
			statsSummary.SummaryTime = now.Format(time.RFC3339)
			statsSummary.StatDate = startTime.Format("2006-01-02")
			go writeSummaryStats(config, statsSummary)

			tags := map[string]string{
				"deliveryservice": statsSummary.DeliveryService,
				"cdn":             statsSummary.CdnName,
			}

			fields := map[string]interface{}{
				"value": maxGbps,
			}
			pt, err := influx.NewPoint(
				statsSummary.StatName,
				tags,
				fields,
				startTime,
			)
			if err != nil {
				errHndlr(err, ERROR)
				continue
			}
			bp.AddPoint(pt)

			// write bytes served data to traffic_ops
			statsSummary.StatName = "daily_bytesserved"
			statsSummary.StatValue = strconv.FormatFloat(bytesServedTb, 'f', 2, 64)
			go writeSummaryStats(config, statsSummary)

			pt, err = influx.NewPoint(
				statsSummary.StatName,
				tags,
				fields,
				startTime,
			)
			if err != nil {
				errHndlr(err, ERROR)
				continue
			}
			bp.AddPoint(pt)
		}
		config.BpsChan <- bp
		log.Info("Collected daily stats @ ", now)
	}
}

func queryDB(con influx.Client, cmd string, database string) (res []influx.Result, err error) {
	q := influx.Query{
		Command:  cmd,
		Database: database,
	}
	if response, err := con.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return
}

func writeSummaryStats(config StartupConfig, statsSummary traffic_ops.StatsSummary) {
	to, err := traffic_ops.Login(config.ToURL, config.ToUser, config.ToPasswd, true)
	if err != nil {
		newErr := fmt.Errorf("Could not store summary stats! Error logging in to %v: %v", config.ToURL, err)
		log.Error(newErr)
		return
	}
	err = to.AddSummaryStats(statsSummary)
	if err != nil {
		log.Error(err)
	}
}

func getToData(config StartupConfig, init bool, configChan chan RunningConfig) {
	var runningConfig RunningConfig
	to, err := traffic_ops.Login(config.ToURL, config.ToUser, config.ToPasswd, true)
	if err != nil {
		msg := fmt.Sprintf("Error logging in to %v: %v", config.ToURL, err)
		if init {
			panic(msg)
		}
		log.Error(msg)
		return
	}

	servers, err := to.Servers()
	if err != nil {
		msg := fmt.Sprintf("Error getting server list from %v: %v ", config.ToURL, err)
		if init {
			panic(msg)
		}
		log.Error(msg)
		return
	}

	runningConfig.CacheGroupMap = make(map[string]string)
	for _, server := range servers {
		runningConfig.CacheGroupMap[server.HostName] = server.Location
		if server.Type == "INFLUXDB" && server.Status == "ONLINE" {
			fqdn := server.HostName + "." + server.DomainName
			port, err := strconv.ParseInt(server.TcpPort, 10, 32)
			if err != nil {
				port = 8086 //default port
			}
			runningConfig.InfluxDBProps = append(runningConfig.InfluxDBProps, struct {
				Fqdn string
				Port int64
			}{fqdn, port})
		}
	}

	cacheStatPath := "/publish/CacheStats?hc=1&stats="
	dsStatPath := "/publish/DsStats?hc=1&wildcard=1&stats="
	parameters, err := to.Parameters("TRAFFIC_STATS")
	if err != nil {
		msg := fmt.Sprintf("Error getting parameter list from %v: %v", config.ToURL, err)
		if init {
			panic(msg)
		}
		log.Error(msg)
		return
	}

	for _, param := range parameters {
		if param.Name == "DsStats" {
			statName := param.Value
			dsStatPath += "," + statName
		} else if param.Name == "CacheStats" {
			cacheStatPath += "," + param.Value
		}
	}
	cacheStatPath = strings.Replace(cacheStatPath, "=,", "=", 1)
	dsStatPath = strings.Replace(dsStatPath, "=,", "=", 1)

	runningConfig.HealthUrls = make(map[string]map[string]string)
	for _, server := range servers {
		if server.Type == "RASCAL" && server.Status == config.StatusToMon {
			cdnName := server.CdnName
			if cdnName == "" {
				log.Error("Unable to find CDN name for " + server.HostName + ".. skipping")
				continue
			}

			if runningConfig.HealthUrls[cdnName] == nil {
				runningConfig.HealthUrls[cdnName] = make(map[string]string)
			}
			url := "http://" + server.IpAddress + cacheStatPath
			runningConfig.HealthUrls[cdnName]["CacheStats"] = url
			url = "http://" + server.IpAddress + dsStatPath
			runningConfig.HealthUrls[cdnName]["DsStats"] = url
		}
	}

	lastSummaryTimeStr, err := to.SummaryStatsLastUpdated("daily_maxgbps")
	if err != nil {
		errHndlr(err, ERROR)
	} else {
		lastSummaryTime, err := time.Parse("2006-01-02 15:04:05", lastSummaryTimeStr)
		if err != nil {
			errHndlr(err, ERROR)
		} else {
			runningConfig.LastSummaryTime = lastSummaryTime
		}
	}

	configChan <- runningConfig
}

func calcMetrics(cdnName string, url string, cacheGroupMap map[string]string, config StartupConfig, runningConfig RunningConfig) {
	sampleTime := int64(time.Now().Unix())
	// get the data from trafficMonitor
	trafMonData, err := getURL(url)
	if err != nil {
		log.Error("Unable to connect to Traffic Monitor @ ", url, " - skipping timeslot")
		return
	}

	if strings.Contains(url, "CacheStats") {
		err = calcCacheValues(trafMonData, cdnName, sampleTime, cacheGroupMap, config)
	} else if strings.Contains(url, "DsStats") {
		err = calcDsValues(trafMonData, cdnName, sampleTime, config)
	} else {
		log.Warn("Don't know what to do with ", url)
	}
}

func errHndlr(err error, severity int) {
	if err != nil {
		switch {
		case severity == ERROR:
			log.Error(err)
		case severity == FATAL:
			log.Error(err)
			panic(err)
		}
	}
}

/* the ds json looks like:
{
  "deliveryService": {
    "linear-gbr-hls-sbr": {
      "	.us-ma-woburn.kbps": [{
        "index": 520281,
        "time": 1398893383605,
        "value": "0",
        "span": 520024
      }],
      "location.us-de-newcastle.kbps": [{
        "index": 520281,
        "time": 1398893383605,
        "value": "0",
        "span": 517707
      }],
    }
 }
*/
func calcDsValues(rascalData []byte, cdnName string, sampleTime int64, config StartupConfig) error {
	type DsStatsJSON struct {
		Pp              string `json:"pp"`
		Date            string `json:"date"`
		DeliveryService map[string]map[string][]struct {
			Index uint64 `json:"index"`
			Time  int    `json:"time"`
			Value string `json:"value"`
			Span  uint64 `json:"span"`
		} `json:"deliveryService"`
	}

	var jData DsStatsJSON
	err := json.Unmarshal(rascalData, &jData)
	errHndlr(err, ERROR)

	statCount := 0
	bps, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:        "deliveryservice_stats",
		Precision:       "ms",
		RetentionPolicy: config.DsRetentionPolicy,
	})
	for dsName, dsData := range jData.DeliveryService {
		for dsMetric, dsMetricData := range dsData {
			//create dataKey (influxDb series)
			var cachegroup, statName string
			if strings.Contains(dsMetric, "total.") {
				s := strings.Split(dsMetric, ".")
				cachegroup, statName = s[0], s[1]
			} else {
				s := strings.Split(dsMetric, ".")
				cachegroup, statName = s[1], s[2]
			}

			//convert stat time to epoch
			statTime := strconv.Itoa(dsMetricData[0].Time)
			msInt, err := strconv.ParseInt(statTime, 10, 64)
			if err != nil {
				errHndlr(err, ERROR)
			}

			newTime := time.Unix(0, msInt*int64(time.Millisecond))
			//convert stat value to float
			statValue := dsMetricData[0].Value
			statFloatValue, err := strconv.ParseFloat(statValue, 64)
			if err != nil {
				statFloatValue = 0.0
			}
			tags := map[string]string{
				"deliveryservice": dsName,
				"cdn":             cdnName,
				"cachegroup":      cachegroup,
			}

			fields := map[string]interface{}{
				"value": statFloatValue,
			}
			pt, err := influx.NewPoint(
				statName,
				tags,
				fields,
				newTime,
			)
			if err != nil {
				errHndlr(err, ERROR)
				continue
			}
			bps.AddPoint(pt)
			statCount++
		}
	}
	config.BpsChan <- bps
	log.Info("Collected ", statCount, " deliveryservice stats values for ", cdnName, " @ ", sampleTime)
	return nil
}

/* The caches json looks like:
{
	caches: {
		odol-atsmid-est-01: { },
		odol-atsec-sfb-05: {
			ats.proxy.process.net.read_bytes: [
				{
					index: 332545,
					time: 1396711793883,
					value: "5547585527895",
					span: 1
				}
			],
			ats.proxy.process.http.transaction_counts.hit_fresh.process: [
				{
					index: 332545,
					time: 1396711793883,
					value: "2109053611",
					span: 1
				}
			],
		}
	}
}
*/

func calcCacheValues(trafmonData []byte, cdnName string, sampleTime int64, cacheGroupMap map[string]string, config StartupConfig) error {

	type CacheStatsJSON struct {
		Pp     string `json:"pp"`
		Date   string `json:"date"`
		Caches map[string]map[string][]struct {
			Index uint64 `json:"index"`
			Time  int    `json:"time"`
			Value string `json:"value"`
			Span  uint64 `json:"span"`
		} `json:"caches"`
	}
	var jData CacheStatsJSON
	err := json.Unmarshal(trafmonData, &jData)
	errHndlr(err, ERROR)

	statCount := 0
	bps, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:        "cache_stats",
		Precision:       "ms",
		RetentionPolicy: config.CacheRetentionPolicy,
	})
	if err != nil {
		errHndlr(err, ERROR)
	}
	for cacheName, cacheData := range jData.Caches {
		for statName, statData := range cacheData {
			dataKey := statName
			dataKey = strings.Replace(dataKey, ".bandwidth", ".kbps", 1)
			dataKey = strings.Replace(dataKey, "-", "_", -1)

			//Get the stat time and convert to epoch
			statTime := strconv.Itoa(statData[0].Time)
			msInt, err := strconv.ParseInt(statTime, 10, 64)
			if err != nil {
				errHndlr(err, ERROR)
			}

			newTime := time.Unix(0, msInt*int64(time.Millisecond))
			//Get the stat value and convert to float
			statValue := statData[0].Value
			statFloatValue, err := strconv.ParseFloat(statValue, 64)
			if err != nil {
				statFloatValue = 0.00
			}
			tags := map[string]string{
				"cachegroup": cacheGroupMap[cacheName],
				"hostname":   cacheName,
				"cdn":        cdnName,
			}

			fields := map[string]interface{}{
				"value": statFloatValue,
			}
			pt, err := influx.NewPoint(
				dataKey,
				tags,
				fields,
				newTime,
			)
			if err != nil {
				errHndlr(err, ERROR)
				continue
			}
			bps.AddPoint(pt)
			statCount++
		}
	}
	config.BpsChan <- bps
	log.Debug("Collected ", statCount, " cache stats values for ", cdnName, " @ ", sampleTime)
	return nil
}

func getURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func influxConnect(config StartupConfig, runningConfig RunningConfig) (influx.Client, error) {
	// Connect to InfluxDb
	var urls []string

	for _, InfluxHost := range runningConfig.InfluxDBProps {
		u := fmt.Sprintf("http://%s:%d", InfluxHost.Fqdn, InfluxHost.Port)
		urls = append(urls, u)
	}

	for len(urls) > 0 {
		n := rand.Intn(len(urls))
		url := urls[n]
		urls = append(urls[:n], urls[n+1:]...)

		conf := influx.HTTPConfig{
			Addr:     url,
			Username: config.InfluxUser,
			Password: config.InfluxPassword,
		}

		con, err := influx.NewHTTPClient(conf)
		if err != nil {
			errHndlr(err, ERROR)
			continue
		}

		return con, nil
	}

	err := errors.New("Could not connect to any of the InfluxDb servers that are ONLINE in traffic ops.")
	return nil, err
}

func sendMetrics(config StartupConfig, runningConfig RunningConfig, bps influx.BatchPoints, retry bool) {
	//influx connection
	influxClient, err := influxConnect(config, runningConfig)
	if err != nil {
		if retry {
			config.BpsChan <- bps
		}
		errHndlr(err, ERROR)
		return
	}
	influxClient.Write(bps)

	log.Info(fmt.Sprintf("Sent %v stats for %v", len(bps.Points()), bps.Database()))
}

//IntMin returns the lesser of two ints
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//FloatMax returns the greater of two float64 values
func FloatMax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
