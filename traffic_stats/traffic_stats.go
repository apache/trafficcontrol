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
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	log "github.com/cihub/seelog"
	influx "github.com/influxdata/influxdb/client/v2"
)

const UserAgent = "traffic-stats"
const TrafficOpsRequestTimeout = time.Second * time.Duration(10)

const (
	// FATAL will exit after printing error
	FATAL = iota
	// ERROR will just keep going, print error
	ERROR
	// WARN will keep going and print a warning
	WARN
)

const (
	defaultPollingInterval             = 10
	defaultDailySummaryPollingInterval = 60
	defaultConfigInterval              = 300
	defaultPublishingInterval          = 30
	defaultMaxPublishSize              = 10000
)

// StartupConfig contains all fields necessary to create a traffic stats session.
type StartupConfig struct {
	ToUser                      string   `json:"toUser"`
	ToPasswd                    string   `json:"toPasswd"`
	ToURL                       string   `json:"toUrl"`
	InfluxUser                  string   `json:"influxUser"`
	InfluxPassword              string   `json:"influxPassword"`
	InfluxURLs                  []string `json:"influxUrls"`
	PollingInterval             int      `json:"pollingInterval"`
	DailySummaryPollingInterval int      `json:"dailySummaryPollingInterval"`
	PublishingInterval          int      `json:"publishingInterval"`
	ConfigInterval              int      `json:"configInterval"`
	MaxPublishSize              int      `json:"maxPublishSize"`
	StatusToMon                 string   `json:"statusToMon"`
	SeelogConfig                string   `json:"seelogConfig"`
	CacheRetentionPolicy        string   `json:"cacheRetentionPolicy"`
	DsRetentionPolicy           string   `json:"dsRetentionPolicy"`
	DailySummaryRetentionPolicy string   `json:"dailySummaryRetentionPolicy"`
	BpsChan                     chan influx.BatchPoints
	InfluxDBs                   []*InfluxDBProps
}

// RunningConfig is used to store runtime configuration for Traffic Stats.  This includes information
// about caches, cachegroups, and health urls
type RunningConfig struct {
	HealthUrls      map[string]map[string]string  // the 1st map key is CDN_name, the second is DsStats or CacheStats
	CacheMap        map[string]traffic_ops.Server // map hostName to cache
	LastSummaryTime time.Time
}

//InfluxDBProps contains URL and connection information for InfluxDB servers
type InfluxDBProps struct {
	URL          string
	InfluxClient influx.Client
}

//Timers struct contains all the timers
type Timers struct {
	Poll         <-chan time.Time
	DailySummary <-chan time.Time
	Publish      <-chan time.Time
	Config       <-chan time.Time
}

func main() {
	var Bps map[string]influx.BatchPoints
	var config StartupConfig
	var err error
	var tickers Timers

	configFile := flag.String("cfg", "", "The config file")
	flag.Parse()

	config, err = loadStartupConfig(*configFile, config)

	if err != nil {
		errHndlr(err, FATAL)
	}

	Bps = make(map[string]influx.BatchPoints)
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
				sendMetrics(config, runningConfig, val, false)
			}
			os.Exit(0)
		case <-tickers.Publish:
			for key, val := range Bps {
				go sendMetrics(config, runningConfig, val, true)
				delete(Bps, key)
			}
		case runningConfig = <-configChan:
		case <-tickers.Config:
			go getToData(config, false, configChan)
		case <-tickers.Poll:
			for cdnName, urls := range runningConfig.HealthUrls {
				for _, url := range urls {
					log.Debug(cdnName, " -> ", url)
					go calcMetrics(cdnName, url, runningConfig.CacheMap, config, runningConfig)
				}
			}
		case now := <-tickers.DailySummary:
			go calcDailySummary(now, config, runningConfig)
		case batchPoints := <-config.BpsChan:
			log.Debug("Received ", len(batchPoints.Points()), " stats")
			key := fmt.Sprintf("%s%s", batchPoints.Database(), batchPoints.RetentionPolicy())
			bp, ok := Bps[key]
			if ok {
				for _, p := range batchPoints.Points() {
					bp.AddPoint(p)
				}
				log.Debug("Aggregating ", len(bp.Points()), " stats to ", key)
			} else {
				Bps[key] = batchPoints
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
	if config.MaxPublishSize == 0 {
		config.MaxPublishSize = defaultMaxPublishSize
	}

	logger, err := log.LoggerFromConfigAsFile(config.SeelogConfig)
	if err != nil {
		return config, fmt.Errorf("error reading Seelog config %s", config.SeelogConfig)
	}
	log.ReplaceLogger(logger)
	log.Info("Replaced logger, see log file according to", config.SeelogConfig)

	if len(config.InfluxURLs) == 0 {
		return config, fmt.Errorf("No InfluxDB urls provided in influxUrls, please provide at least one valid URL.  e.g. \"influxUrls\": [\"http://localhost:8086\"]")
	}
	for _, url := range config.InfluxURLs {
		influxDBProps := InfluxDBProps{
			URL: url,
		}
		config.InfluxDBs = append(config.InfluxDBs, &influxDBProps)
	}

	//Close old connections explicitly
	for _, host := range oldConfig.InfluxDBs {
		if host.InfluxClient != nil {
			host.InfluxClient.Close()
		}
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
		influxClient, err := influxConnect(config)
		if err != nil {
			log.Error("Could not connect to InfluxDb to get daily summary stats!!")
			errHndlr(err, ERROR)
			return
		}

		bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
			Database:        "daily_stats",
			Precision:       "s",
			RetentionPolicy: config.DailySummaryRetentionPolicy,
		})

		calcDailyMaxGbps(influxClient, bp, startTime, endTime, config)
		calcDailyBytesServed(influxClient, bp, startTime, endTime, config)
		log.Info("Collected daily stats @ ", now)
	}
}

func calcDailyMaxGbps(client influx.Client, bp influx.BatchPoints, startTime time.Time, endTime time.Time, config StartupConfig) {
	kilobitsToGigabits := 1000000.00
	queryString := fmt.Sprintf(`select time, cdn, max(value) from "monthly"."bandwidth.cdn.1min" where time > '%s' and time < '%s' group by cdn`, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
	log.Infof("queryString = %v\n", queryString)
	res, err := queryDB(client, queryString, "cache_stats")
	if err != nil {
		log.Errorf("An error occured getting max bandwidth! %v\n", err)
		return
	}
	if res != nil && len(res[0].Series) > 0 {
		for _, row := range res[0].Series {
			for _, record := range row.Values {
				t := record[0].(string)
				if record[1] != nil {
					cdn := record[1].(string)
					value, err := record[2].(json.Number).Float64()
					if err != nil {
						log.Errorf("Couldn't parse value from record %v\n", record)
						continue
					}
					value = value / kilobitsToGigabits
					statTime, _ := time.Parse(time.RFC3339, t)
					log.Infof("max gbps for cdn %v = %v", cdn, value)
					var statsSummary traffic_ops.StatsSummary
					statsSummary.CDNName = cdn
					statsSummary.DeliveryService = "all"
					statsSummary.StatName = "daily_maxgbps"
					statsSummary.StatValue = strconv.FormatFloat(value, 'f', 2, 64)
					statsSummary.SummaryTime = time.Now().Format(time.RFC3339)
					statsSummary.StatDate = statTime.Format("2006-01-02")
					go writeSummaryStats(config, statsSummary)

					//write to influxdb
					tags := map[string]string{"cdn": cdn, "deliveryservice": "all"}
					fields := map[string]interface{}{
						"value": value,
					}
					pt, err := influx.NewPoint(
						"daily_maxgbps",
						tags,
						fields,
						statTime,
					)
					if err != nil {
						log.Errorf("error adding data point for max Gbps...%v\n", err)
						continue
					}
					bp.AddPoint(pt)
				}
			}
		}
	}
	config.BpsChan <- bp
}

func calcDailyBytesServed(client influx.Client, bp influx.BatchPoints, startTime time.Time, endTime time.Time, config StartupConfig) {
	bytesToTerabytes := 1000000000.00
	sampleTimeSecs := 60.00
	bitsTobytes := 8.00
	queryString := fmt.Sprintf(`select mean(value) from "monthly"."bandwidth.cdn.1min" where time > '%s' and time < '%s' group by time(1m), cdn`, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
	log.Infof("queryString = %v\n", queryString)
	res, err := queryDB(client, queryString, "cache_stats")
	if err != nil {
		log.Error("An error occured getting max bandwidth!\n")
		return
	}
	if res != nil && len(res[0].Series) > 0 {
		for _, row := range res[0].Series {
			bytesServed := float64(0)
			cdn := row.Tags["cdn"]
			for _, record := range row.Values {
				if record[1] != nil {
					value, err := record[1].(json.Number).Float64()
					if err != nil {
						log.Errorf("Couldn't parse value from record %v\n", record)
						continue
					}
					bytesServed += value * sampleTimeSecs / bitsTobytes
				}
			}
			bytesServedTB := bytesServed / bytesToTerabytes
			log.Infof("TBytes served for cdn %v = %v", cdn, bytesServedTB)
			//write to Traffic Ops
			var statsSummary traffic_ops.StatsSummary
			statsSummary.CDNName = cdn
			statsSummary.DeliveryService = "all"
			statsSummary.StatName = "daily_bytesserved"
			statsSummary.StatValue = strconv.FormatFloat(bytesServedTB, 'f', 2, 64)
			statsSummary.SummaryTime = time.Now().Format(time.RFC3339)
			statsSummary.StatDate = startTime.Format("2006-01-02")
			go writeSummaryStats(config, statsSummary)
			//write to Influxdb
			tags := map[string]string{"cdn": cdn, "deliveryservice": "all"}
			fields := map[string]interface{}{
				"value": bytesServedTB, //converted to TB
			}
			pt, err := influx.NewPoint(
				"daily_bytesserved",
				tags,
				fields,
				startTime,
			)
			if err != nil {
				log.Errorf("error adding creating data point for max Gbps...%v\n", err)
				continue
			}
			bp.AddPoint(pt)
		}
		config.BpsChan <- bp
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
	to, err := traffic_ops.LoginWithAgent(config.ToURL, config.ToUser, config.ToPasswd, true, UserAgent, false, TrafficOpsRequestTimeout)
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
	to, err := traffic_ops.LoginWithAgent(config.ToURL, config.ToUser, config.ToPasswd, true, UserAgent, false, TrafficOpsRequestTimeout)
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

	runningConfig.CacheMap = make(map[string]traffic_ops.Server)
	for _, server := range servers {
		runningConfig.CacheMap[server.HostName] = server
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
		if server.Type == "RASCAL" && server.Status != config.StatusToMon {
			log.Debugf("Skipping %s%s.  Looking for status %s but got status %s", server.HostName, server.DomainName, config.StatusToMon, server.Status)
			continue
		}

		if server.Type == "RASCAL" && server.Status == config.StatusToMon {
			cdnName := server.CDNName
			if cdnName == "" {
				log.Error("Unable to find CDN name for " + server.HostName + ".. skipping")
				continue
			}

			if runningConfig.HealthUrls[cdnName] == nil {
				runningConfig.HealthUrls[cdnName] = make(map[string]string)
			}
			url := "http://" + server.IPAddress + cacheStatPath
			runningConfig.HealthUrls[cdnName]["CacheStats"] = url
			url = "http://" + server.IPAddress + dsStatPath
			runningConfig.HealthUrls[cdnName]["DsStats"] = url
		}
	}

	lastSummaryTimeStr, err := to.SummaryStatsLastUpdated("daily_maxgbps")
	if err != nil {
		errHndlr(err, ERROR)
	} else {
		lastSummaryTime, err := time.Parse("2006-01-02 15:04:05+00", lastSummaryTimeStr)
		if err != nil {
			errHndlr(err, ERROR)
		} else {
			runningConfig.LastSummaryTime = lastSummaryTime
		}
	}

	configChan <- runningConfig
}

func calcMetrics(cdnName string, url string, cacheMap map[string]traffic_ops.Server, config StartupConfig, runningConfig RunningConfig) {
	sampleTime := int64(time.Now().Unix())
	// get the data from trafficMonitor
	trafMonData, err := getURL(url)
	if err != nil {
		log.Error("Unable to connect to Traffic Monitor @ ", url, " - skipping timeslot")
		return
	}

	if strings.Contains(url, "CacheStats") {
		err = calcCacheValues(trafMonData, cdnName, sampleTime, cacheMap, config)
		errHndlr(err, ERROR)
	} else if strings.Contains(url, "DsStats") {
		err = calcDsValues(trafMonData, cdnName, sampleTime, config)
		errHndlr(err, ERROR)
	} else {
		log.Warn("Don't know what to do with ", url)
	}
}

func errHndlr(err error, severity int) {
	if err != nil {
		switch {
		case severity == WARN:
			log.Warn(err)
		case severity == ERROR:
			log.Error(err)
		case severity == FATAL:
			log.Error(err)
			panic(err)
		}
	}
}

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
	if err != nil {
		return fmt.Errorf("could not unmarshall deliveryservice stats JSON - %v", err)
	}

	statCount := 0
	bps, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:        "deliveryservice_stats",
		Precision:       "ms",
		RetentionPolicy: config.DsRetentionPolicy,
	})
	for dsName, dsData := range jData.DeliveryService {
		for dsMetric, dsMetricData := range dsData {
			var cachegroup, statName string
			tags := map[string]string{
				"deliveryservice": dsName,
				"cdn":             cdnName,
			}

			s := strings.Split(dsMetric, ".")
			if strings.Contains(dsMetric, "type.") {
				cachegroup = "all"
				statName = s[2]
				tags["type"] = s[1]
			} else if strings.Contains(dsMetric, "total.") {
				cachegroup, statName = s[0], s[1]
			} else {
				cachegroup, statName = s[1], s[2]
			}

			tags["cachegroup"] = cachegroup

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

func calcCacheValues(trafmonData []byte, cdnName string, sampleTime int64, cacheMap map[string]traffic_ops.Server, config StartupConfig) error {

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
	if err != nil {
		return fmt.Errorf("could not unmarshall cache stats JSON - %v", err)
	}

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
		cache := cacheMap[cacheName]

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
				"cachegroup": cache.Cachegroup,
				"hostname":   cacheName,
				"cdn":        cdnName,
				"type":       cache.Type,
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
	log.Info("Collected ", statCount, " cache stats values for ", cdnName, " @ ", sampleTime)
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

func influxConnect(config StartupConfig) (influx.Client, error) {
	hosts := config.InfluxDBs
	for len(hosts) > 0 {
		n := rand.Intn(len(hosts))
		host := hosts[n]
		hosts = append(hosts[:n], hosts[n+1:]...)
		parsedURL, _ := url.Parse(host.URL)
		if parsedURL.Scheme == "udp" {
			conf := influx.UDPConfig{
				Addr: parsedURL.Host,
			}
			con, err := influx.NewUDPClient(conf)
			if err != nil {
				errHndlr(fmt.Errorf("An error occurred creating udp client. %v\n", err), ERROR)
				continue
			}
			return con, nil
		}
		//if not udp assume HTTP client
		conf := influx.HTTPConfig{
			Addr:     parsedURL.String(),
			Username: config.InfluxUser,
			Password: config.InfluxPassword,
		}
		con, err := influx.NewHTTPClient(conf)
		if err != nil {
			errHndlr(fmt.Errorf("An error occurred creating HTTP client.  %v\n", err), ERROR)
			continue
		}
		//Close old connections explicitly
		if host.InfluxClient != nil {
			host.InfluxClient.Close()
		}
		host.InfluxClient = con
		_, _, err = con.Ping(10)
		if err != nil {
			errHndlr(err, WARN)
			continue
		}
		return con, nil
	}
	err := errors.New("Could not connect to any of the InfluxDb servers defined in the influxUrls config.")
	return nil, err
}

func sendMetrics(config StartupConfig, runningConfig RunningConfig, bps influx.BatchPoints, retry bool) {
	influxClient, err := influxConnect(config)
	if err != nil {
		if retry {
			config.BpsChan <- bps
		}
		errHndlr(err, ERROR)
		return
	}

	pts := bps.Points()
	for len(pts) > 0 {
		chunkBps, err := influx.NewBatchPoints(influx.BatchPointsConfig{
			Database:        bps.Database(),
			Precision:       bps.Precision(),
			RetentionPolicy: bps.RetentionPolicy(),
		})
		if err != nil {
			if retry {
				config.BpsChan <- chunkBps
			}
			errHndlr(err, ERROR)
		}
		for _, p := range pts[:intMin(config.MaxPublishSize, len(pts))] {
			chunkBps.AddPoint(p)
		}
		pts = pts[intMin(config.MaxPublishSize, len(pts)):]

		err = influxClient.Write(chunkBps)
		if err != nil {
			if retry {
				config.BpsChan <- chunkBps
			}
			errHndlr(err, ERROR)
		} else {
			log.Info(fmt.Sprintf("Sent %v stats for %v", len(chunkBps.Points()), chunkBps.Database()))
		}
	}
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func floatMax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
