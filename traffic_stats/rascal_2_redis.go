/*
     Copyright 2015 Comcast Cable Communications Management, LLC

     Licensed under the Apache License, Version 2.0 (the "License");
     you may not use this file except in compliance with the License.
     You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

     Unless required by applicable law or agreed to in writing, software
     distributed under the License is distributed on an "AS IS" BASIS,
     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
     See the License for the specific language governing permissions and
     limitations under the License.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/fzzy/radix/redis"
	traffic_ops "github.comcast.com/cdneng/traffic_control/traffic_ops/client"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	FATAL = iota // Exit after printing error
	ERROR = iota // Just keep going, print error
)

const defaultRedisInterval = 10

type StartupConfig struct {
	TmUser        string                       `json:"tmUser"`
	TmPasswd      string                       `json:"tmPasswd"`
	TmUrl         string                       `json:"tmUrl"`
	RedisString   string                       `json:"redisString"`
	RedisInterval int                          `json:"redisInterval"`
	StatusToMon   string                       `json:statusToMon"`
	SeelogConfig  string                       `json:seelogConfig"`
	DsAggregate   map[string]AggregationConfig `json:"dsAggregate"`
}

type AggregationConfig struct {
	RedisKey string `json:"redisKey"` // add more stuff here as necessary
}

type RunningConfig struct {
	HealthUrls      map[string]map[string]string // they 1st map key is CDN_name, the second is DsStats or CacheStats
	CacheGroupMap   map[string]string            // map hostName to cacheGroup
	RetentionPeriod int64                        // how long in seconds to keep the data in the Redis database
}

type redisPool struct {
	connect_string string
	client_chan    chan *redis.Client
	reply_chans    chan chan *redis.Reply
}

func main() {
	configFile := flag.String("cfg", "", "The config file")
	testSummary := flag.Bool("testSummary", false, "Test summary mode")
	flag.Parse()
	file, err := os.Open(*configFile)
	errHndlr(err, FATAL)
	decoder := json.NewDecoder(file)
	config := &StartupConfig{}
	err = decoder.Decode(&config)
	errHndlr(err, FATAL)

	if config.RedisInterval == 0 {
		config.RedisInterval = defaultRedisInterval
	}

	logger, err := log.LoggerFromConfigAsFile(config.SeelogConfig)
	defer log.Flush()

	if err != nil {
		panic("error reading " + config.SeelogConfig)
	}

	fmt.Println("Replacing logger, see log file according to " + config.SeelogConfig)
	if *testSummary {
		fmt.Println("WARNING: testSummary is on!")
	}
	log.ReplaceLogger(logger)

	runtime.GOMAXPROCS(runtime.NumCPU())

	runningConfig, nil := getTmData(config, true)
	go houseKeeping(runningConfig, *testSummary)

	freeList := NewPool(16, config.RedisString)
	<-time.NewTimer(time.Now().Truncate(time.Duration(config.RedisInterval) * time.Second).Add(time.Duration(config.RedisInterval) * time.Second).Sub(time.Now())).C
	tickerChan := time.Tick(time.Duration(config.RedisInterval) * time.Second)
	for now := range tickerChan {
		if now.Second() == 30 {
			trc, err := getTmData(config, false)

			if err == nil {
				runningConfig = trc
			}
		}
		for cdnName, urls := range runningConfig.HealthUrls {
			for _, url := range urls {
				// log.Info(cdnName, "   ", statName, " -> ", url)
				if *testSummary {
					fmt.Println("Skipping stat write - testSummary mode is ON!")
					continue
				}
				go rascalToRedis(cdnName, url, runningConfig.CacheGroupMap, freeList, config)
			}
		}
	}
}

func NewPool(size int, connect_string string) *redisPool {
	rp := new(redisPool)
	rp.connect_string = connect_string
	rp.client_chan = make(chan *redis.Client, size)
	rp.reply_chans = make(chan chan *redis.Reply, size)
	return rp
}

func (rp *redisPool) NewClient() (*redis.Client, error) {
	var redisClient *redis.Client
	var err error
	select {
	case redisClient = <-rp.client_chan:
		// we got an existing connection off the freeList
		// log.Info("Reusing free redis connection")
	default:
		log.Info("Creating new redis connection")
		redisClient, err = redis.DialTimeout("tcp", rp.connect_string, time.Duration(10)*time.Second)
		if err != nil {
			log.Error("ERROR: Unable to connect to redis server ", rp.connect_string, " and no connections on freeList - skipping timeslot")
			return nil, err
		}
	}

	return redisClient, nil
}

// may be needed later, not sure if we need it now.
// func (rp *redisPool) NewClientPing(timeout time.Duration) (*redis.Client, error) {
// 	redisClient, err := rp.NewClient()
// 	// ping connection and open new one if need be
// 	var reply_chan chan *redis.Reply
// 	select {
// 	case reply_chan = <-rp.reply_chans:
// 	default:
// 		reply_chan = make(chan *redis.Reply)
// 	}
// 	go func() { reply_chan <- redisClient.Cmd("PING") }()
// 	select {
// 	case reply := <-reply_chan:
// 		if reply.Err != nil {
// 			go redisClient.Close()
// 			redisClient, err = rp.NewClient()
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 	case <-time.After(timeout):
// 		go redisClient.Close()
// 		redisClient, err = rp.NewClient()
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	select {
// 	case rp.reply_chans <- reply_chan:
// 	default:
// 	}

// 	return redisClient, nil
// }

func (rp *redisPool) FreeClient(c *redis.Client) {
	select {
	case rp.client_chan <- c:
		// redisClient on free list; nothing more to do.
	default:
		// Free list full, close connection and let it get GC'd
		log.Info("Free list is full - closing redis connection")
		go c.Close()
	}
}

func getTmData(config *StartupConfig, init bool) (RunningConfig, error) {
	var runningConfig RunningConfig
	tm, err := traffic_ops.Login(config.TmUrl, config.TmUser, config.TmPasswd, true)
	if err != nil {
		msg := fmt.Sprintf("Error logging in to %v: %v", config.TmUrl, err)

		if init {
			panic(msg)
		} else {
			log.Error(msg)
			return runningConfig, err
		}
	}

	myHostName, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("Error getting my hostname: %v", err))
	}
	myHostName = strings.Split(myHostName, ".")[0]
	log.Info("I am " + myHostName)
	myProfile := ""
	myLocation := ""
	servers, err := tm.Servers()
	if err != nil {
		msg := fmt.Sprintf("Error getting server list from %v: %v ", config.TmUrl, err)

		if init {
			panic(msg)
		} else {
			log.Error(msg)
			return runningConfig, err
		}
	}
	runningConfig.CacheGroupMap = make(map[string]string)
	for _, server := range servers {
		if server.HostName == myHostName {
			log.Info("My location is " + server.Location)
			myLocation = server.Location
			myProfile = server.Profile
		}
		runningConfig.CacheGroupMap[server.HostName] = server.Location
	}

	log.Info("Searching for " + config.StatusToMon + " RASCAL servers in " + myLocation + "...")
	cacheStatPath := "/publish/CacheStats?hc=1&stats="
	dsStatPath := "/publish/DsStats?hc=1&wildcard=1&stats="
	parameters, err := tm.Parameters(myProfile)
	if err != nil {
		msg := fmt.Sprintf("Error getting parameter list from %v: %v", config.TmUrl, err)

		if init {
			panic(msg)
		} else {
			log.Error(msg)
			return runningConfig, err
		}
	}
	runningConfig.RetentionPeriod = 8640 // hardcoded default, if the param doesn't exist, it'll use this
	for _, param := range parameters {
		if param.Name == "DsStats" {
			statName := param.Value
			dsStatPath += "," + statName
		} else if param.Name == "CacheStats" {
			cacheStatPath += "," + param.Value
		} else if param.Name == "RetentionPeriod" {
			runningConfig.RetentionPeriod, err = strconv.ParseInt(param.Value, 10, 64)
			if err != nil {
				log.Error(param.Name, " - error converting ", param.Value, " to Int: ", err)
			}
		}
	}
	cacheStatPath = strings.Replace(cacheStatPath, "=,", "=", 1)
	dsStatPath = strings.Replace(dsStatPath, "=,", "=", 1)

	runningConfig.HealthUrls = make(map[string]map[string]string)
	for _, server := range servers {
		if server.Type == "RASCAL" && server.Status == config.StatusToMon && server.Location == myLocation {
			cdnName := ""
			parameters, _ := tm.Parameters(server.Profile)
			for _, param := range parameters {
				if param.Name == "CDN_name" && param.ConfigFile == "rascal-config.txt" {
					cdnName = param.Value
					break
				}
			}

			if cdnName == "" {
				log.Error("Unable to find CDN name for " + server.HostName + ".. skipping")
				continue
			}

			if runningConfig.HealthUrls[cdnName] == nil {
				runningConfig.HealthUrls[cdnName] = make(map[string]string)
			}
			url := "http://" + server.IpAddress + dsStatPath
			runningConfig.HealthUrls[cdnName]["DsStats"] = url
			log.Info(myHostName, ": ", cdnName, " -> ", url)
			url = "http://" + server.IpAddress + cacheStatPath
			runningConfig.HealthUrls[cdnName]["CacheStats"] = url
			log.Info(myHostName, ": ", cdnName, " -> ", url)
		}
	}
	return runningConfig, nil
}

func summarizeYesterday(redisClient *redis.Client) {

	t := time.Now().Add(-86400 * time.Second)
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()) // reset to start of yesterday 00:00::00
	endTime := startTime.Add(86400 * time.Second)
	endUTime := endTime.Unix()
	startUTime := startTime.Unix()

	log.Info("Summarizing from ", startTime, " (", startUTime, ") to ", endTime, " (", endUTime, ")")

	uTimes := make(map[string][]int64)
	keyList, err := redisClient.Cmd("keys", "*:*:all:all:kbps").List() // all cdns all deliveryservices
	errHndlr(err, ERROR)
	for _, keyName := range keyList {
		log.Info("lrange ", keyName, " ", -9000, " ", -1)
		bwVals, err := redisClient.Cmd("lrange", keyName, -9000, -1).List()
		errHndlr(err, ERROR)
		keyParts := strings.Split(keyName, ":")
		if len(keyParts) != 5 {
			log.Info("Error parsing key: ", keyName)
			continue
		}
		cdnName := keyParts[0]
		deliveryService := keyParts[1]
		cacheGroup := keyParts[2]
		hostName := keyParts[3]
		// statName := keyParts[4]
		if cacheGroup != "all" {
			continue
		}

		// only need to get the times once per CDN
		if uTimes[cdnName] == nil {
			keyName = cdnName + ":tstamp"
			log.Info("lrange ", keyName, " ", -9000, " ", -1)
			uTlist, err := redisClient.Cmd("lrange", keyName, -9000, -1).List()
			errHndlr(err, ERROR)
			for _, tStamp := range uTlist {
				intVal, err := strconv.ParseInt(tStamp, 10, 64)
				if err != nil {
					log.Error(cdnName, " - error converting ", tStamp, " to Int: ", err)
					continue
				}
				uTimes[cdnName] = append(uTimes[cdnName], intVal)
			}
		}

		errorForKey := false
		bytesServed := float64(0)
		maxBps := float64(0)
		prevUtime := startUTime
		for index, sampleTime := range uTimes[cdnName] {
			errHndlr(err, ERROR)
			if sampleTime < startUTime {
				continue
			}
			sampleVal := float64(0)
			if index < len(bwVals) {
				sampleVal, err = strconv.ParseFloat(bwVals[index], 64)
				if err != nil {
					log.Error(keyName, " - error converting ", bwVals[index], " to Float: ", err, " skipping stat summary!")
					errorForKey = true
					break
				}
				if maxBps < sampleVal {
					maxBps = sampleVal
				}
			}
			duration := sampleTime - prevUtime
			bytesServed += float64(duration) * sampleVal / 8
			prevUtime = sampleTime
		}
		if !errorForKey {
			dailyKey := cdnName + ":" + deliveryService + ":" + cacheGroup + ":" + hostName
			log.Info(dailyKey, " bw:", maxBps, " bs:", bytesServed)
			r := redisClient.Cmd("rpush", dailyKey+":daily_maxkbps", strconv.FormatInt(startUTime, 10)+":"+strconv.FormatInt(int64(maxBps), 10))
			errHndlr(r.Err, ERROR)
			r = redisClient.Cmd("rpush", dailyKey+":daily_bytesserved", strconv.FormatInt(startUTime, 10)+":"+strconv.FormatInt(int64(bytesServed), 10))
			errHndlr(r.Err, ERROR)
		}
	}
}

func houseKeeping(runningConfig RunningConfig, testSummary bool) {
	redisClient, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)

	if err != nil {
		errHndlr(err, ERROR)
		return
	}

	defer redisClient.Close()

	minuteChan := time.Tick(time.Minute)
	dayOfWeek := time.Now().Day()
	hourOfDay := time.Now().Hour()

	for now := range minuteChan {
		log.Info("Housekeeping! (", dayOfWeek, ", ", hourOfDay, ")")
		if testSummary || now.Day() != dayOfWeek {
			summarizeYesterday(redisClient)
			dayOfWeek = time.Now().Day()
		}
		/*
			if hourOfDay != time.Now().Hour() {
				log.Info("Saving DB..")
				r := redisClient.Cmd("bgsave")
				errHndlr(r.Err, ERROR)
				hourOfDay = time.Now().Hour()
			}
		*/
		keyList, err := redisClient.Cmd("keys", "*").List()
		errHndlr(err, ERROR)
		for _, keyName := range keyList {
			if !strings.HasSuffix(keyName, "maxKbps") {
				r := redisClient.Cmd("ltrim", keyName, -1*runningConfig.RetentionPeriod, -1) // 1 day
				if r.Err != nil {
					panic(fmt.Sprintf("%v for %v", r.Err.Error(), keyName))
				}
			}
		}
	}
}

func rascalToRedis(cdnName string, url string, cacheGroupMap map[string]string, freeList *redisPool, config *StartupConfig) {
	sampleTime := int64(time.Now().Unix())
	// get the data from rascal
	rascalData, err := getUrl(url)
	if err != nil {
		log.Info("ERROR: Unable to connect to rascal @ ", url, " - skipping timeslot")
		return
	}

	// get a connection to redis
	redisClient, _ := freeList.NewClient()
	// store the data to redis
	if strings.Contains(url, "CacheStats") {
		err = storeCacheValues(rascalData, cdnName, sampleTime, cacheGroupMap, redisClient)
	} else if strings.Contains(url, "DsStats") {
		err = storeDsValues(rascalData, cdnName, sampleTime, redisClient, config.DsAggregate)
	} else {
		log.Info("Don't know what to do with ", url)
	}
	// return the redis connection to the pool
	freeList.FreeClient(redisClient)
}

func errHndlr(err error, severity int) {
	if err != nil {
		switch {
		case severity == ERROR:
			log.Info(err)
		case severity == FATAL:
			panic(err)
		}
	}
}

/* the ds json looks like:
{
  "deliveryService": {
    "linear-gbr-hls-sbr": {
      "location.us-ma-woburn.kbps": [{
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
func storeDsValues(rascalData []byte, cdnName string, sampleTime int64, redisClient *redis.Client, dsAggregate map[string]AggregationConfig) error {
	type DsStatsJson struct {
		Pp              string `json:"pp"`
		Date            string `json:"date"`
		DeliveryService map[string]map[string][]struct {
			Index uint64 `json:"index"`
			Time  uint64 `json:"time"`
			Value string `json:"value"`
			Span  uint64 `json:"span"`
		} `json:"deliveryService"`
	}

	var jData DsStatsJson
	err := json.Unmarshal(rascalData, &jData)
	errHndlr(err, ERROR)
	statCount := 0
	statTotals := make(map[string]float64)
	for dsName, dsData := range jData.DeliveryService {
		for dsMetric, dsMetricData := range dsData {
			keyPart := strings.Replace(dsMetric, "location.", "", -1)
			keyPart = strings.Replace(keyPart, ".kbps", ":all:kbps", -1)
			keyPart = strings.Replace(keyPart, ".tps", ":all:tps", -1)
			keyPart = strings.Replace(keyPart, ".status", ":all:status", -1)
			keyPart = strings.Replace(keyPart, "total:all:", "all:all:", -1) // for consistency all everywhere
			redisKey := cdnName + ":" + dsName + ":" + keyPart
			statValue := dsMetricData[0].Value
			//fmt.Printf("%s  ->%s\n", redisKey, statValue)
			statCount++

			aggConfig, exists := dsAggregate[dsMetric]

			if exists {
				statFloatValue, err := strconv.ParseFloat(statValue, 64)

				if err != nil {
					statFloatValue = 0.0
				}

				statTotals[cdnName+":all:all:all:"+aggConfig.RedisKey] += statFloatValue
			}

			r := redisClient.Cmd("rpush", redisKey, statValue)
			errHndlr(r.Err, ERROR)
		}
	}
	for totalKey, totalVal := range statTotals {
		r := redisClient.Cmd("rpush", totalKey, strconv.FormatFloat(totalVal, 'f', 2, 64))
		errHndlr(r.Err, ERROR)
		statCount++
	}
	log.Info("Saved ", statCount, " ds values for ", cdnName, " @ ", sampleTime)
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

func storeCacheValues(rascalData []byte, cdnName string, sampleTime int64, cacheGroupMap map[string]string, redisClient *redis.Client) error {
	/* note about the redis data:
	keys are cdnName:deliveryService:cacheGroup:cacheName:statName
	*/

	type CacheStatsJson struct {
		Pp     string `json:"pp"`
		Date   string `json:"date"`
		Caches map[string]map[string][]struct {
			Index uint64 `json:"index"`
			Time  uint64 `json:"time"`
			Value string `json:"value"`
			Span  uint64 `json:"span"`
		} `json:"caches"`
	}

	var jData CacheStatsJson
	err := json.Unmarshal(rascalData, &jData)
	errHndlr(err, ERROR)
	statCount := 0
	statTotals := make(map[string]float64)
	for cacheName, cacheData := range jData.Caches {
		for statName, statData := range cacheData {
			redisKey := cdnName + ":all:" + cacheGroupMap[cacheName] + ":" + cacheName + ":" + statName
			redisKey = strings.Replace(redisKey, ":bandwidth", ":kbps", 1)
			statValue := statData[0].Value
			//fmt.Printf("%s  ->%s\n", redisKey, statValue)
			statCount++
			statFloatValue, err := strconv.ParseFloat(statValue, 64)
			if err != nil {
				statFloatValue = 0.0
			}
			statTotals[cdnName+":all:"+cacheGroupMap[cacheName]+":all:"+statName] += statFloatValue
			statTotals[cdnName+":all:all:all:"+statName] += statFloatValue
			if statName == "maxKbps" {
				r := redisClient.Cmd("zadd", redisKey, sampleTime, statValue) // only care for the last val here.
				errHndlr(r.Err, ERROR)
			} else {
				r := redisClient.Cmd("rpush", redisKey, statValue)
				errHndlr(r.Err, ERROR)
			}
		}
	}
	for totalKey, totalVal := range statTotals {
		totalKey = strings.Replace(totalKey, ":bandwidth", ":kbps", 1)
		if strings.Contains(totalKey, "maxKbps") {
			r := redisClient.Cmd("zadd", totalKey, sampleTime, strconv.FormatFloat(totalVal, 'f', 2, 64))
			errHndlr(r.Err, ERROR)
		} else {
			r := redisClient.Cmd("rpush", totalKey, strconv.FormatFloat(totalVal, 'f', 2, 64))
			errHndlr(r.Err, ERROR)
		}
		statCount++
	}
	r := redisClient.Cmd("rpush", cdnName+":tstamp", sampleTime)
	errHndlr(r.Err, ERROR)
	log.Info("Saved ", statCount, " values for ", cdnName, " @ ", sampleTime)
	return nil
}

func getUrl(url string) ([]byte, error) {

	// log.Info(url, "  >>>>>")
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
