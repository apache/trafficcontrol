package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"

	log "github.com/cihub/seelog"
	traffic_ops "github.com/comcast/traffic_control/traffic_ops/client"
	influx "github.com/influxdb/influxdb/client"
)

const (
	FATAL = iota // Exit after printing error
	ERROR = iota // Just keep going, print error
)

const defaultPollingInterval = 10

type StartupConfig struct {
	ToUser                      string `json:"toUser"`
	ToPasswd                    string `json:"toPasswd"`
	ToUrl                       string `json:"toUrl"`
	InfluxUser                  string `json:"influxUser"`
	InfluxPassword              string `json:"influxPassword"`
	StatusToMon                 string `json:"statusToMon"`
	SeelogConfig                string `json:"seelogConfig"`
	DailySummaryPollingInterval int    `json:"dailySummaryPollingInterval"`
}

type TrafOpsData struct {
	InfluxDbProps   []InfluxDbProps
	LastSummaryTime string
}

type InfluxDbProps struct {
	Fqdn string
	Port int64
}

func main() {
	configFile := flag.String("cfg", "", "The config file")
	test := flag.Bool("test", false, "Test mode")
	flag.Parse()
	file, err := os.Open(*configFile)
	errHndlr(err, FATAL)
	decoder := json.NewDecoder(file)
	config := &StartupConfig{}
	err = decoder.Decode(&config)
	errHndlr(err, FATAL)
	pollingInterval := 60
	if config.DailySummaryPollingInterval > 0 {
		pollingInterval = config.DailySummaryPollingInterval
	}

	logger, err := log.LoggerFromConfigAsFile(config.SeelogConfig)
	defer log.Flush()
	if err != nil {
		panic("error reading Seelog config " + config.SeelogConfig)
	}
	fmt.Println("Replacing logger, see log file according to " + config.SeelogConfig)
	if *test {
		fmt.Println("WARNING: test mode is on!")
	}
	log.ReplaceLogger(logger)

	runtime.GOMAXPROCS(runtime.NumCPU())

	t := time.Now().Add(-86400 * time.Second)
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()) // reset to start of yesterday 00:00::00
	endTime := startTime.Add(86400 * time.Second)
	formatStartTime := startTime.Format("2006-01-02T15:04:05-00:00")
	formatEndTime := endTime.Format("2006-01-02T15:04:05-00:00")
	endUTime := endTime.Unix()
	startUTime := startTime.Unix()
	pts := make([]influx.Point, 0)

	<-time.NewTimer(time.Now().Truncate(time.Duration(pollingInterval) * time.Second).Add(time.Duration(pollingInterval) * time.Second).Sub(time.Now())).C
	tickerChan := time.Tick(time.Duration(pollingInterval) * time.Second)
	for now := range tickerChan {
		//get TrafficOps Data
		trafOpsData, err := getToData(config, false)
		if err != nil {
			errHndlr(err, FATAL)
		}
		lastSummaryTime, err := time.Parse("2006-01-02 15:04:05", trafOpsData.LastSummaryTime)
		if err != nil {
			errHndlr(err, ERROR)
		}
		log.Infof("lastSummaryTime is %v", lastSummaryTime)
		if lastSummaryTime.Day() != now.Day() {
			log.Info("Summarizing from ", startTime, " (", startUTime, ") to ", endTime, " (", endUTime, ")")
			// influx connection
			influxClient, err := influxConnect(config, trafOpsData)
			if err != nil {
				log.Error("Could not connect to InfluxDb to get daily summary stats!!")
				errHndlr(err, ERROR)
				continue
			}
			//create influxdb query
			log.Infof("SELECT sum(value)/6 FROM bandwidth where time > '%v' and time < '%v' group by time(60s), cdn fill(0)", formatStartTime, formatEndTime)
			q := fmt.Sprintf("SELECT sum(value)/6 FROM bandwidth where time > '%v' and time < '%v' group by time(60s), cdn fill(0)", formatStartTime, formatEndTime)
			res, err := queryDB(influxClient, q, "cache_stats")
			if err != nil {
				fmt.Printf("err = %v\n", err)
				errHndlr(err, ERROR)
				continue
			}
			//loop throgh series
			for _, row := range res[0].Series {
				prevUtime := startUTime
				var cdn string
				max := float64(0)
				bytesServed := float64(0)
				cdn = row.Tags["cdn"]
				for _, record := range row.Values {
					kbps, err := record[1].(json.Number).Float64()
					if err != nil {
						errHndlr(err, ERROR)
						continue
					}
					sampleTime, err := time.Parse("2006-01-02T15:04:05Z", record[0].(string))
					if err != nil {
						errHndlr(err, ERROR)
						continue
					}
					sampleUTime := sampleTime.Unix()
					if kbps > max {
						max = kbps
					}
					duration := sampleUTime - prevUtime
					bytesServed += float64(duration) * kbps / 8
					prevUtime = sampleUTime
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
				statsSummary.SummaryTime = now.Format("2006-01-02 15:04:05")
				statsSummary.StatDate = startTime.Format("2006-01-02")
				err = writeSummaryStats(config, statsSummary)
				if err != nil {
					log.Error("Could not store daily_maxgbps stats in traffic ops!")
					errHndlr(err, ERROR)
				}
				//write to influxdb
				pts = append(pts,
					influx.Point{
						Measurement: statsSummary.StatName,
						Tags: map[string]string{
							"deliveryservice": statsSummary.DeliveryService,
							"cdn":             statsSummary.CdnName,
						},
						Fields: map[string]interface{}{
							"value": maxGbps,
						},
						Time:      startTime,
						Precision: "s",
					},
				)
				//write bytes served data to traffic_ops
				statsSummary.StatName = "daily_bytesserved"
				statsSummary.StatValue = strconv.FormatFloat(bytesServedTb, 'f', 2, 64)
				err = writeSummaryStats(config, statsSummary)
				if err != nil {
					log.Error("Could not store daily_bytesserved stats in traffic ops!")
					errHndlr(err, ERROR)
				}
				pts = append(pts,
					influx.Point{
						Measurement: statsSummary.StatName,
						Tags: map[string]string{
							"deliveryservice": statsSummary.DeliveryService,
							"cdn":             statsSummary.CdnName,
						},
						Fields: map[string]interface{}{
							"value": bytesServedTb,
						},
						Time:      startTime,
						Precision: "s",
					},
				)
			}
			log.Infof("Writing daily stats to influxDb")
			bps := influx.BatchPoints{
				Points:          pts,
				Database:        "daily_stats",
				RetentionPolicy: "daily_stats",
			}
			_, err = influxClient.Write(bps)
			if err != nil {
				errHndlr(err, ERROR)
			}
		}
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

func queryDB(con *influx.Client, cmd string, database string) (res []influx.Result, err error) {
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

func influxConnect(config *StartupConfig, trafOps TrafOpsData) (*influx.Client, error) {
	//Connect to InfluxDb
	activeServers := len(trafOps.InfluxDbProps)
	rand.Seed(42)
	//if there is only 1 active, use it
	if activeServers == 1 {
		u, err := url.Parse(fmt.Sprintf("http://%s:%d", trafOps.InfluxDbProps[0].Fqdn, trafOps.InfluxDbProps[0].Port))
		if err != nil {
			return nil, err
		}
		conf := influx.Config{
			URL:      *u,
			Username: config.InfluxUser,
			Password: config.InfluxPassword,
		}
		con, err := influx.NewClient(conf)
		if err != nil {
			return nil, err
		}
		_, _, err = con.Ping()
		if err != nil {
			return nil, err
		}
		return con, nil
	} else if activeServers > 1 {
		//try to connect to all ONLINE servers until we find one that works
		for i := 0; i < activeServers; i++ {
			u, err := url.Parse(fmt.Sprintf("http://%s:%d", trafOps.InfluxDbProps[i].Fqdn, trafOps.InfluxDbProps[i].Port))
			if err != nil {
				errHndlr(err, ERROR)
			} else {
				conf := influx.Config{
					URL:      *u,
					Username: config.InfluxUser,
					Password: config.InfluxPassword,
				}
				con, err := influx.NewClient(conf)
				if err != nil {
					errHndlr(err, ERROR)
					continue
				} else {
					_, _, err = con.Ping()
					if err != nil {
						errHndlr(err, ERROR)
						continue
					} else {
						return con, nil
					}
				}
			}
		}
		err := errors.New("Could not connect to any of the InfluxDb servers that are ONLINE in traffic ops.")
		return nil, err
	} else {
		err := errors.New("No online InfluxDb servers could be found!")
		return nil, err
	}
}

func getToData(config *StartupConfig, init bool) (TrafOpsData, error) {
	var trafOpsData TrafOpsData
	tm, err := traffic_ops.Login(config.ToUrl, config.ToUser, config.ToPasswd, true)
	if err != nil {
		msg := fmt.Sprintf("Error logging in to %v: %v", config.ToUrl, err)
		if init {
			panic(msg)
		} else {
			log.Error(msg)
			return trafOpsData, err
		}
	}

	servers, err := tm.Servers()
	if err != nil {
		msg := fmt.Sprintf("Error getting server list from %v: %v ", config.ToUrl, err)
		if init {
			panic(msg)
		} else {
			log.Error(msg)
			return trafOpsData, err
		}
	}
	for _, server := range servers {
		if server.Type == "INFLUXDB" && server.Status == "ONLINE" {
			fqdn := server.HostName + "." + server.DomainName
			port, err := strconv.ParseInt(server.TcpPort, 10, 32)
			if err != nil {
				port = 8086 //default port
			}
			trafOpsData.InfluxDbProps = append(trafOpsData.InfluxDbProps, InfluxDbProps{Fqdn: fqdn, Port: port})
		}
	}
	lastSummaryTime, err := tm.SummaryStatsLastUpdated("daily_maxgbps")
	if err != nil {
		errHndlr(err, ERROR)
	}
	trafOpsData.LastSummaryTime = lastSummaryTime
	return trafOpsData, nil
}

func writeSummaryStats(config *StartupConfig, statsSummary traffic_ops.StatsSummary) error {
	tm, err := traffic_ops.Login(config.ToUrl, config.ToUser, config.ToPasswd, true)
	if err != nil {
		msg := fmt.Sprintf("Could not store summary stats! Error logging in to %v: %v", config.ToUrl, err)
		log.Error(msg)
		return err
	}
	err = tm.AddSummaryStats(statsSummary)
	if err != nil {
		return err
	}
	return nil
}
