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
	"crypto/tls"
	"crypto/x509"
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

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"

	"github.com/Shopify/sarama"
	"github.com/cihub/seelog"
	influx "github.com/influxdata/influxdb/client/v2"
)

const UserAgent = "traffic-stats"

const (
	defaultTrafficOpsRequestTimeout    = 10
	defaultPollingInterval             = 10
	defaultDailySummaryPollingInterval = 60
	defaultConfigInterval              = 300
	defaultPublishingInterval          = 30
	defaultMaxPublishSize              = 10000
)

type Logger struct {
	Error   log.LogLocation `json:"error"`
	Warning log.LogLocation `json:"warning"`
	Info    log.LogLocation `json:"info"`
	Debug   log.LogLocation `json:"debug"`
	Event   log.LogLocation `json:"event"`
}

func (l Logger) ErrorLog() log.LogLocation {
	return l.Error
}
func (l Logger) WarningLog() log.LogLocation {
	return l.Warning
}
func (l Logger) InfoLog() log.LogLocation {
	return l.Info
}
func (l Logger) DebugLog() log.LogLocation {
	return l.Debug
}
func (l Logger) EventLog() log.LogLocation {
	return l.Event
}

var defaultLogger Logger = Logger{
	Error:   log.LogLocationStderr,
	Warning: log.LogLocationStderr,
	Info:    log.LogLocationStderr,
	Debug:   log.LogLocationNull,
	Event:   log.LogLocationStderr,
}

// StartupConfig contains all fields necessary to create a traffic stats session.
type StartupConfig struct {
	ToUser                      string   `json:"toUser"`
	ToPasswd                    string   `json:"toPasswd"`
	ToURL                       string   `json:"toUrl"`
	ToRequestTimeoutSeconds     int      `json:"toRequestTimeout"`
	DisableInflux               bool     `json:"disableInflux"`
	InfluxUser                  string   `json:"influxUser"`
	InfluxPassword              string   `json:"influxPassword"`
	InfluxURLs                  []string `json:"influxUrls"`
	PollingInterval             int      `json:"pollingInterval"`
	DailySummaryPollingInterval int      `json:"dailySummaryPollingInterval"`
	PublishingInterval          int      `json:"publishingInterval"`
	ConfigInterval              int      `json:"configInterval"`
	MaxPublishSize              int      `json:"maxPublishSize"`
	StatusToMon                 string   `json:"statusToMon"`
	SeelogConfig                *string  `json:"seelogConfig"`
	LogConfig                   *Logger  `json:"logs"`
	CacheRetentionPolicy        string   `json:"cacheRetentionPolicy"`
	DsRetentionPolicy           string   `json:"dsRetentionPolicy"`
	DailySummaryRetentionPolicy string   `json:"dailySummaryRetentionPolicy"`
	BpsChan                     chan influx.BatchPoints
	InfluxDBs                   []*InfluxDBProps
	KafkaConfig                 KafkaConfig `json:"kafkaConfig"`
}

type KafkaConfig struct {
	Enable        bool   `json:"enable"`
	Brokers       string `json:"brokers"`
	Topic         string `json:"topic"`
	RequiredAcks  int    `json:"requiredAcks"`
	EnableTls     bool   `json:"enableTls"`
	RootCA        string `json:"rootCA"`
	ClientCert    string `json:"clientCert"`
	ClientCertKey string `json:"clientCertKey"`
}

type KafkaCluster struct {
	producer *sarama.AsyncProducer
	client   *sarama.Client
}

type KafkaJSON struct {
	Name   string                 `json:"name"`
	Tags   map[string]string      `json:"tags"`
	Fields map[string]interface{} `json:"fields"`
	Time   time.Time              `json:"time"`
}

type DataExporter interface {
	ExportData(config StartupConfig, bps influx.BatchPoints, retry bool)
}

func (c *KafkaCluster) ExportData(config StartupConfig, bps influx.BatchPoints, retry bool) {
	err := publishToKafka(config, bps, c)
	if err != nil {
		log.Errorln("Unable to export to Kafka", err)
	}
}

func (influx InfluxClient) ExportData(config StartupConfig, bps influx.BatchPoints, retry bool) {
	sendMetrics(config, bps, retry)
}

type InfluxClient struct {
}

var useSeelog bool = true

// RunningConfig is used to store runtime configuration for Traffic Stats.  This includes information
// about caches, cachegroups, and health urls
type RunningConfig struct {
	HealthUrls      map[string]map[string][]string // the 1st map key is CDN_name, the second is DsStats or CacheStats
	CacheMap        map[string]tc.ServerV5         // map hostName to cache
	LastSummaryTime time.Time
}

// InfluxDBProps contains URL and connection information for InfluxDB servers
type InfluxDBProps struct {
	URL          string
	InfluxClient influx.Client
}

// Timers struct contains all the timers
type Timers struct {
	Poll         <-chan time.Time
	DailySummary <-chan time.Time
	Publish      <-chan time.Time
	Config       <-chan time.Time
}

func info(args ...interface{}) {
	if useSeelog {
		seelog.Info(args...)
	} else {
		log.Infoln(args...)
	}
}
func infof(format string, args ...interface{}) {
	if useSeelog {
		seelog.Infof(format, args...)
	} else {
		log.Infof(format, args...)
	}
}
func errorln(args ...interface{}) {
	if useSeelog {
		seelog.Error(args...)
	} else {
		log.Errorln(args...)
	}
}
func errorf(format string, args ...interface{}) {
	if useSeelog {
		seelog.Errorf(format, args...)
	} else {
		log.Errorf(format, args...)
	}
}
func warn(args ...interface{}) {
	if useSeelog {
		seelog.Warn(args...)
	} else {
		log.Warnln(args...)
	}
}
func warnf(format string, args ...interface{}) {
	if useSeelog {
		seelog.Warnf(format, args...)
	} else {
		log.Warnf(format, args...)
	}
}
func debug(args ...interface{}) {
	if useSeelog {
		seelog.Debug(args...)
	} else {
		log.Debugln(args...)
	}
}
func debugf(format string, args ...interface{}) {
	if useSeelog {
		seelog.Debugf(format, args...)
	} else {
		log.Debugf(format, args...)
	}
}

func main() {
	var Bps map[string]influx.BatchPoints
	var config StartupConfig
	var err error
	var tickers Timers

	configFile := flag.String("cfg", "", "The config file")
	flag.Parse()
	if *configFile == "" {
		flag.Usage()
		panic("-cfg is required")
	}

	config, err = loadStartupConfig(*configFile, config)

	if err != nil {
		err = fmt.Errorf("could not load startup config: %v", err)
		errorln(err)
		panic(err)
	}

	Bps = make(map[string]influx.BatchPoints)
	config.BpsChan = make(chan influx.BatchPoints)

	defer seelog.Flush()

	configChan := make(chan RunningConfig)
	go getToData(config, true, configChan)
	runningConfig := <-configChan

	c := newKakfaCluster(config.KafkaConfig)

	tickers = setTimers(config)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	hupChan := make(chan os.Signal, 1)
	signal.Notify(hupChan, syscall.SIGHUP)

	dataExporters := []DataExporter{}

	if config.KafkaConfig.Enable && c != nil {
		dataExporters = append(dataExporters, c)
	}

	if !config.DisableInflux {
		influx := InfluxClient{}
		dataExporters = append(dataExporters, influx)
	}

	for {
		select {
		case <-hupChan:
			info("HUP Received - reloading config")
			newConfig, err := loadStartupConfig(*configFile, config)
			if err != nil {
				errorf("could not load startup config: %v", err)
			} else {
				config = newConfig
				tickers = setTimers(config)
			}
		case <-termChan:
			info("Shutdown Request Received - Sending stored metrics then quitting")
			for _, val := range Bps {
				for _, dataExporter := range dataExporters {
					dataExporter.ExportData(config, val, false)
				}
			}
			startShutdown(c)
			os.Exit(0)
		case <-tickers.Publish:
			for key, val := range Bps {
				for _, dataExporter := range dataExporters {
					go dataExporter.ExportData(config, val, true)
				}
				delete(Bps, key)
			}
		case runningConfig = <-configChan:
		case <-tickers.Config:
			go getToData(config, false, configChan)
		case <-tickers.Poll:
			for cdnName, urls := range runningConfig.HealthUrls {
				for _, u := range urls {
					debug(cdnName, " -> ", u)
					go calcMetrics(cdnName, u, runningConfig.CacheMap, config)
				}
			}
		case now := <-tickers.DailySummary:
			go calcDailySummary(now, config, runningConfig)
		case batchPoints := <-config.BpsChan:
			debug("Received ", len(batchPoints.Points()), " stats")
			key := fmt.Sprintf("%s%s", batchPoints.Database(), batchPoints.RetentionPolicy())
			bp, ok := Bps[key]
			if ok {
				for _, p := range batchPoints.Points() {
					bp.AddPoint(p)
				}
				debug("Aggregating ", len(bp.Points()), " stats to ", key)
			} else {
				Bps[key] = batchPoints
				debug("Created ", key)
			}
		}
	}
}

func startShutdown(c *KafkaCluster) {
	if c == nil {
		return
	}

	infof("Starting cluster shutdown, closing producer and client")
	if err := (*c.producer).Close(); err != nil {
		warnf("Error closing producer for cluster:  %v", err)
	}
	if err := (*c.client).Close(); err != nil {
		warnf("Error closing client for cluster:  %v", err)
	}
	c.producer = nil
	infof("Finished cluster shutdown")
}

func TlsConfig(config KafkaConfig) (*tls.Config, error) {
	if !config.EnableTls {
		return nil, nil
	}
	c := &tls.Config{}
	if config.RootCA != "" {
		infof("Loading TLS root CA certificate from %s", config.RootCA)
		caCert, err := ioutil.ReadFile(config.RootCA)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		c.RootCAs = caCertPool
	} else {
		warnf("No TLS root CA defined")
	}
	if config.ClientCert != "" {
		infof("Loading TLS client certificate")
		if config.ClientCertKey == "" {
			return nil, fmt.Errorf("client cert path defined without key path")
		}
		certPEM, err := ioutil.ReadFile(config.ClientCert)
		if err != nil {
			return nil, err
		}
		keyPEM, err := ioutil.ReadFile(config.ClientCertKey)
		if err != nil {
			return nil, err
		}
		cert, err := tls.X509KeyPair(certPEM, keyPEM)
		if err != nil {
			return nil, err
		}
		c.Certificates = []tls.Certificate{cert}
	}
	return c, nil
}

func newKakfaCluster(config KafkaConfig) *KafkaCluster {

	if !config.Enable {
		return nil
	}

	brokers := strings.Split(config.Brokers, ",")
	sc := sarama.NewConfig()

	tlsConfig, err := TlsConfig(config)
	if err != nil {
		errorln("Unable to create TLS config", err)
		return nil
	}

	sc.Producer.RequiredAcks = sarama.RequiredAcks(config.RequiredAcks)
	if config.EnableTls && tlsConfig != nil {
		sc.Net.TLS.Enable = true
		sc.Net.TLS.Config = tlsConfig
	}

	cl, err := sarama.NewClient(brokers, sc)

	if err != nil {
		errorln("Unable to create client", err)
		return nil
	}

	p, err := sarama.NewAsyncProducerFromClient(cl)

	if err != nil {
		errorln("Unable to create producer", err)
		return nil
	}
	c := &KafkaCluster{
		producer: &p,
		client:   &cl,
	}

	return c
}

func publishToKafka(config StartupConfig, bps influx.BatchPoints, c *KafkaCluster) error {

	input := (*c.producer).Input()

	for _, point := range bps.Points() {

		var KafkaJSON KafkaJSON
		KafkaJSON.Name = point.Name()
		KafkaJSON.Tags = point.Tags()
		KafkaJSON.Fields, _ = point.Fields()
		KafkaJSON.Time = point.Time()

		message, err := json.Marshal(KafkaJSON)

		if err != nil {
			return err
		}

		topic := config.KafkaConfig.Topic

		input <- &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(message),
		}
	}
	return nil
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

	if config.ToRequestTimeoutSeconds <= 0 {
		config.ToRequestTimeoutSeconds = defaultTrafficOpsRequestTimeout
	}

	if config.LogConfig != nil {
		if err = log.InitCfg(config.LogConfig); err != nil {
			return config, fmt.Errorf("initializing logging configuration: %w", err)
		}
		useSeelog = false
	} else if config.SeelogConfig != nil {
		logger, err := seelog.LoggerFromConfigAsFile(*config.SeelogConfig)
		if err != nil {
			return config, fmt.Errorf("error reading Seelog config '%s': %w", *config.SeelogConfig, err)
		}
		seelog.ReplaceLogger(logger)
		useSeelog = true
		infof("Replaced logger, see seelog file according to '%s'", *config.SeelogConfig)
		warn("seelog-based logging is deprecated, please switch to using the 'logs' property")
	} else {
		if err = log.InitCfg(defaultLogger); err != nil {
			return config, fmt.Errorf("initializing default logger: %w", err)
		}
		useSeelog = false
		warn("No logging configuration found in configuration file - default logging to stderr will be used")
	}

	if config.DisableInflux {
		return config, nil
	}

	if len(config.InfluxURLs) == 0 {
		return config, fmt.Errorf("No InfluxDB urls provided in influxUrls, please provide at least one valid URL.  e.g. \"influxUrls\": [\"http://localhost:8086\"]")
	}
	for _, u := range config.InfluxURLs {
		influxDBProps := InfluxDBProps{
			URL: u,
		}
		config.InfluxDBs = append(config.InfluxDBs, &influxDBProps)
	}

	//Close old connections explicitly
	for _, host := range oldConfig.InfluxDBs {
		if host.InfluxClient != nil {
			if err := host.InfluxClient.Close(); err != nil {
				errorf("closing influx client: %s", err)
			}
		}
	}

	return config, nil
}

func calcDailySummary(now time.Time, config StartupConfig, runningConfig RunningConfig) {
	if config.DisableInflux {
		info("Skipping daily stats since InfluxDB is not enabled")
		return
	}
	infof("lastSummaryTime is %v", runningConfig.LastSummaryTime)
	if runningConfig.LastSummaryTime.Day() != now.Day() {
		startTime := now.Truncate(24 * time.Hour).Add(-24 * time.Hour)
		endTime := startTime.Add(24 * time.Hour)
		info("Summarizing from ", startTime, " (", startTime.Unix(), ") to ", endTime, " (", endTime.Unix(), ")")

		// influx connection
		influxClient, err := influxConnect(config)
		if err != nil {
			errorf("could not connect to InfluxDb to get daily summary stats: %v", err)
			return
		}

		bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
			Database:        "daily_stats",
			Precision:       "s",
			RetentionPolicy: config.DailySummaryRetentionPolicy,
		})

		calcDailyMaxGbps(influxClient, bp, startTime, endTime, config)
		calcDailyBytesServed(influxClient, bp, startTime, endTime, config)
		info("Collected daily stats @ ", now)
	}
}

func calcDailyMaxGbps(client influx.Client, bp influx.BatchPoints, startTime time.Time, endTime time.Time, config StartupConfig) {
	kilobitsToGigabits := 1000000.00
	queryString := fmt.Sprintf(`select time, cdn, max(value) from "monthly"."bandwidth.cdn.1min" where time > '%s' and time < '%s' group by cdn`, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
	infof("queryString = %v\n", queryString)
	res, err := queryDB(client, queryString, "cache_stats")
	if err != nil {
		errorf("An error occured getting max bandwidth! %v\n", err)
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
						errorf("Couldn't parse value from record %v\n", record)
						continue
					}
					value = value / kilobitsToGigabits
					statTime, _ := time.Parse(time.RFC3339, t)
					infof("max gbps for cdn %v = %v", cdn, value)
					var statsSummary tc.StatsSummaryV5
					statsSummary.CDNName = util.StrPtr(cdn)
					statsSummary.DeliveryService = util.StrPtr("all")
					statsSummary.StatName = util.StrPtr("daily_maxgbps")
					statsSummary.StatValue = util.FloatPtr(value)
					statsSummary.SummaryTime = time.Now()
					statsSummary.StatDate = &statTime
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
						errorf("error adding data point for max Gbps...%v\n", err)
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
	infof("queryString = %s", queryString)
	res, err := queryDB(client, queryString, "cache_stats")
	if err != nil {
		errorln("An error occured getting max bandwidth: ", err)
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
						errorf("Couldn't parse value from record %v", record)
						continue
					}
					bytesServed += value * sampleTimeSecs / bitsTobytes
				}
			}
			bytesServedTB := bytesServed / bytesToTerabytes
			infof("TBytes served for cdn %v = %v", cdn, bytesServedTB)
			//write to Traffic Ops
			var statsSummary tc.StatsSummaryV5
			statsSummary.CDNName = util.StrPtr(cdn)
			statsSummary.DeliveryService = util.StrPtr("all")
			statsSummary.StatName = util.StrPtr("daily_bytesserved")
			statsSummary.StatValue = util.FloatPtr(bytesServedTB)
			statsSummary.SummaryTime = time.Now()
			statsSummary.StatDate = &startTime
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
				errorf("error adding creating data point for max Gbps...%v\n", err)
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

func writeSummaryStats(config StartupConfig, statsSummary tc.StatsSummaryV5) {
	to, _, err := client.LoginWithAgent(config.ToURL, config.ToUser, config.ToPasswd, true, UserAgent, false, time.Duration(config.ToRequestTimeoutSeconds)*time.Second)
	if err != nil {
		newErr := fmt.Errorf("Could not store summary stats! Error logging in to %v: %v", config.ToURL, err)
		errorln(newErr)
		return
	}
	_, _, err = to.CreateSummaryStats(statsSummary, client.RequestOptions{})
	if err != nil {
		errorf("could not create summary stats: %v", err)
	}
}

func getToData(config StartupConfig, init bool, configChan chan RunningConfig) {
	var runningConfig RunningConfig
	to, _, err := client.LoginWithAgent(config.ToURL, config.ToUser, config.ToPasswd, true, UserAgent, false, time.Duration(config.ToRequestTimeoutSeconds)*time.Second)
	if err != nil {
		msg := fmt.Sprintf("Error logging in to %v: %v", config.ToURL, err)
		if init {
			panic(msg)
		}
		errorln(msg)
		return
	}

	servers, _, err := to.GetServers(client.RequestOptions{})
	if err != nil {
		msg := fmt.Sprintf("Error getting server list from %v: %v ", config.ToURL, err)
		if init {
			panic(msg)
		}
		errorln(msg)
		return
	}

	runningConfig.CacheMap = make(map[string]tc.ServerV5)
	for _, server := range servers.Response {
		runningConfig.CacheMap[server.HostName] = server
	}

	cacheStatPath := "/publish/CacheStats?hc=1&wildcard=1&stats="
	dsStatPath := "/publish/DsStats?hc=1&wildcard=1&stats="
	parameters, _, err := to.GetParametersByProfileName("TRAFFIC_STATS", client.RequestOptions{})
	if err != nil {
		msg := fmt.Sprintf("Error getting parameter list from %v: %v", config.ToURL, err)
		if init {
			panic(msg)
		}
		errorln(msg)
		return
	}

	for _, param := range parameters.Response {
		if param.Name == "DsStats" {
			statName := param.Value
			dsStatPath += "," + statName
		} else if param.Name == "CacheStats" {
			cacheStatPath += "," + param.Value
		}
	}
	cacheStatPath = strings.Replace(cacheStatPath, "=,", "=", 1)
	dsStatPath = strings.Replace(dsStatPath, "=,", "=", 1)

	setHealthURLs(config, &runningConfig, cacheStatPath, dsStatPath)

	opts := client.RequestOptions{
		QueryParameters: url.Values{
			"lastSummaryDate": {"true"},
			"statName":        {"daily_maxgbps"},
		},
	}
	lastSummaryTimeResponse, _, err := to.GetSummaryStatsLastUpdated(opts)
	if err != nil {
		errorf("unable to get summary stats last updated: %v", err)
	} else if lastSummaryTimeResponse.Response.SummaryTime == nil {
		warn("unable to get last updated stats summary timestamp: daily_maxgbps stats summary not reported yet")
	} else {
		runningConfig.LastSummaryTime = *lastSummaryTimeResponse.Response.SummaryTime
	}

	configChan <- runningConfig
}

func setHealthURLs(config StartupConfig, runningConfig *RunningConfig, cacheStatPath string, dsStatPath string) {
	runningConfig.HealthUrls = make(map[string]map[string][]string)
	for _, server := range runningConfig.CacheMap {
		if server.Type == tc.MonitorTypeName && server.Status != config.StatusToMon {
			debugf("Skipping %s.%s.  Looking for status %s but got status %s", server.HostName, server.DomainName, config.StatusToMon, server.Status)
			continue
		}

		if server.Type == tc.MonitorTypeName && server.Status == config.StatusToMon {
			cdnName := server.CDN
			if cdnName == "" {
				errorln("Unable to find CDN name for " + server.HostName + ".. skipping")
				continue
			}

			if runningConfig.HealthUrls[cdnName] == nil {
				runningConfig.HealthUrls[cdnName] = make(map[string][]string)
			}
			tcpPort := util.CoalesceToDefault(server.TCPPort)
			healthURL := "http://" + server.HostName + "." + server.DomainName + ":" + strconv.Itoa(tcpPort) + cacheStatPath
			runningConfig.HealthUrls[cdnName]["CacheStats"] = append(runningConfig.HealthUrls[cdnName]["CacheStats"], healthURL)
			healthURL = "http://" + server.HostName + "." + server.DomainName + ":" + strconv.Itoa(tcpPort) + dsStatPath
			runningConfig.HealthUrls[cdnName]["DsStats"] = append(runningConfig.HealthUrls[cdnName]["DsStats"], healthURL)
		}
	}
}

func calcMetrics(cdnName string, urls []string, cacheMap map[string]tc.ServerV5, config StartupConfig) {
	sampleTime := time.Now().Unix()
	// get the data from trafficMonitor
	var trafMonData []byte
	var err error
	var healthURL string
	for _, u := range urls {
		trafMonData, err = getURL(u)
		if err != nil {
			errorf("error getting %s stats URL %s: %v", cdnName, u, err)
			continue
		}
		healthURL = u
		infof("successfully got %s stats URL %s", cdnName, u)
		break
	}
	if healthURL == "" {
		errorf("unable to get any %s stats URL - skipping timeslot", cdnName)
		return
	}

	if strings.Contains(healthURL, "CacheStats") {
		err = calcCacheValues(trafMonData, cdnName, sampleTime, cacheMap, config)
		if err != nil {
			errorf("error calculating cache metric values for CDN %s: %v", cdnName, err)
		}
	} else if strings.Contains(healthURL, "DsStats") {
		err = calcDsValues(trafMonData, cdnName, sampleTime, config)
		if err != nil {
			errorf("error calculating delivery service metric values for CDN %s: %v", cdnName, err)
		}
	} else {
		warn("Don't know what to do with given ", cdnName, " stats URL: ", healthURL)
	}
}

func calcDsValues(tmData []byte, cdnName string, sampleTime int64, config StartupConfig) error {
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
	err := json.Unmarshal(tmData, &jData)
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
			//Get the stat time and make sure it's greater than the time 24 hours ago. If not, skip it so influxdb doesn't throw retention policy errors.
			validTime := time.Now().AddDate(0, 0, -1).UnixNano() / 1000000
			timeStamp := int64(dsMetricData[0].Time)
			if timeStamp < validTime {
				info(fmt.Sprintf("Skipping %v %v: %v is greater than 24 hours old.", dsName, dsMetric, timeStamp))
				continue
			}
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
				errorf("calculating delivery service metric values: error parsing stat time: %v", err)
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
				errorf("calculating delivery service metric values: error creating new influxDB point: %v", err)
				continue
			}
			bps.AddPoint(pt)
			statCount++
		}
	}
	config.BpsChan <- bps
	info("Collected ", statCount, " deliveryservice stats values for ", cdnName, " @ ", sampleTime)
	return nil
}

func calcCacheValues(trafmonData []byte, cdnName string, sampleTime int64, cacheMap map[string]tc.ServerV5, config StartupConfig) error {
	var jData tc.LegacyStats
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
		errorf("calculating cache metric values: creating new influxDB batch points: %v", err)
	}

	for cacheName, cacheData := range jData.Caches {
		cache := cacheMap[string(cacheName)]

		for statName, statData := range cacheData {
			if len(statData) == 0 {
				continue
			}
			//Get the stat time and make sure it's greater than the time 24 hours ago.  If not, skip it so influxdb doesn't throw retention policy errors.
			validTime := time.Now().AddDate(0, 0, -1)
			if statData[0].Time.Before(validTime) {
				info(fmt.Sprintf("Skipping %v %v: %v is greater than 24 hours old.", cacheName, statName, statData[0].Time))
				continue
			}
			dataKey := statName
			dataKey = strings.Replace(dataKey, ".bandwidth", ".kbps", 1)
			dataKey = strings.Replace(dataKey, "-", "_", -1)

			//Get the stat value and convert to float
			statFloatValue := 0.0
			if statsValue, ok := statData[0].Val.(string); !ok {
				warnf("stat data %s with value %v couldn't be converted into string", statName, statData[0].Val)
			} else {
				statFloatValue, err = strconv.ParseFloat(statsValue, 64)
				if err != nil {
					warnf("stat %s with value %v couldn't be converted into a float", statName, statsValue)
				}
			}

			tags := map[string]string{
				"cachegroup": cache.CacheGroup,
				"hostname":   string(cacheName),
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
				statData[0].Time,
			)
			if err != nil {
				errorf("calculating cache metric values: error creating new influxDB point: %v", err)
				continue
			}
			bps.AddPoint(pt)
			statCount++
		}
	}
	config.BpsChan <- bps
	info("Collected ", statCount, " cache stats values for ", cdnName, " @ ", sampleTime)
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
		if host.InfluxClient != nil && parsedURL.Scheme == "http" {
			// NOTE: closing an http client just closes idle connections -- the client can still make new requests
			if err := host.InfluxClient.Close(); err != nil {
				errorf("closing http influx client: %s", err)
			}
			_, _, err := host.InfluxClient.Ping(10)
			if err != nil {
				warnf("pinging InfluxDB: %v", err)
				continue
			}
			return host.InfluxClient, nil
		}
		if parsedURL.Scheme == "udp" {
			conf := influx.UDPConfig{
				Addr: parsedURL.Host,
			}
			con, err := influx.NewUDPClient(conf)
			if err != nil {
				errorf("An error occurred creating InfluxDB UDP client: %v", err)
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
			errorf("An error occurred creating InfluxDB HTTP client: %v", err)
			continue
		}
		host.InfluxClient = con
		_, _, err = con.Ping(10)
		if err != nil {
			warnf("pinging InfluxDB: %v", err)
			continue
		}
		return con, nil
	}
	err := errors.New("could not connect to any of the InfluxDb servers defined in the influxUrls config")
	return nil, err
}

func sendMetrics(config StartupConfig, bps influx.BatchPoints, retry bool) {
	influxClient, err := influxConnect(config)
	if err != nil {
		if retry {
			config.BpsChan <- bps
		}
		errorf("sending metrics to InfluxDB: unable to get InfluxDB client: %v", err)
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
			errorf("sending metrics to InfluxDB: error creating new batch points: %v", err)
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
			errorf("sending metrics to InfluxDB: error writing batch points: %v", err)
		} else {
			info(fmt.Sprintf("Sent %v stats for %v", len(chunkBps.Points()), chunkBps.Database()))
		}
	}
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
