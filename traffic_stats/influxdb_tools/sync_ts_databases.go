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
	"flag"
	"fmt"
	"os"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

type cacheStats struct {
	t        string //time
	value    float64
	cdn      string
	hostname string
}
type deliveryServiceStats struct {
	t               string //time
	value           float64
	cdn             string
	deliveryService string
	cacheGroup      string
}
type dailyStats struct {
	t               string //time
	cdn             string
	deliveryService string
	value           float64
}

func main() {

	sourceURL := flag.String("sourceUrl", "http://server1.kabletown.net:8086", "The influxdb url and port")
	targetURL := flag.String("targetUrl", "http://server2.kabletown.net:8086", "The influxdb url and port")
	database := flag.String("database", "all", "Sync a specific database")
	days := flag.Int("days", 0, "Number of days in the past to sync (today - x days), 0 is all")
	flag.Parse()
	fmt.Printf("syncing %v to %v for %v database(s) for the past %v day(s)\n", *sourceURL, *targetURL, *database, *days)
	sourceClient, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr: *sourceURL,
	})
	if err != nil {
		fmt.Printf("Error creating influx sourceClient: %v\n", err)
		os.Exit(1)
	}
	targetClient, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr: *targetURL,
	})
	if err != nil {
		fmt.Printf("Error creating influx targetClient: %v\n", err)
		os.Exit(1)
	}
	chSize := 1
	if *database == "all" {
		chSize = 3
	}

	ch := make(chan string)

	switch *database {
	case "all":
		go syncCsDb(ch, sourceClient, targetClient, *days)
		go syncDsDb(ch, sourceClient, targetClient, *days)
		go syncDailyDb(ch, sourceClient, targetClient, *days)
	case "cache_stats":
		go syncCsDb(ch, sourceClient, targetClient, *days)
	case "deliveryservice_stats":
		go syncDsDb(ch, sourceClient, targetClient, *days)
	case "daily_stats":
		go syncDailyDb(ch, sourceClient, targetClient, *days)
	}

	for i := 1; i <= chSize; i++ {
		fmt.Println(<-ch)
	}

	fmt.Println("Traffic Stats has been synced!")
}

func syncCsDb(ch chan string, sourceClient influx.Client, targetClient influx.Client, days int) {
	db := "cache_stats"
	fmt.Printf("Syncing %s database...\n", db)
	stats := [...]string{
		"bandwidth.cdn.1min",
		"connections.cdn.1min",
	}
	for _, statName := range stats {
		fmt.Printf("Syncing %s database with %s \n", db, statName)
		syncCacheStat(sourceClient, targetClient, statName, days)
	}
	ch <- fmt.Sprintf("Done syncing %s!\n", db)
}

func syncDsDb(ch chan string, sourceClient influx.Client, targetClient influx.Client, days int) {
	db := "deliveryservice_stats"
	fmt.Printf("Syncing %s database...\n", db)
	stats := [...]string{
		"kbps.ds.1min",
		"max.kbps.ds.1day",
		"tps.ds.1min",
		"tps_2xx.ds.1min",
		"tps_3xx.ds.1min",
		"tps_4xx.ds.1min",
		"tps_5xx.ds.1min",
		"tps_total.ds.1min",
	}
	for _, statName := range stats {
		fmt.Printf("Syncing %s database with %s\n", db, statName)
		syncDeliveryServiceStat(sourceClient, targetClient, statName, days)
	}
	ch <- fmt.Sprintf("Done syncing %s!\n", db)
}

func syncDailyDb(ch chan string, sourceClient influx.Client, targetClient influx.Client, days int) {
	db := "daily_stats"
	fmt.Printf("Syncing %s database...\n", db)
	stats := [...]string{
		"daily_bytesserved",
		"daily_maxgbps",
	}

	for _, statName := range stats {
		fmt.Printf("Syncing %s database with %s\n", db, statName)
		syncDailyStat(sourceClient, targetClient, statName, days)
	}
	ch <- fmt.Sprintf("Done syncing %s!\n", db)

}

func queryDB(client influx.Client, cmd string, db string) (res []influx.Result, err error) {
	q := influx.Query{
		Command:  cmd,
		Database: db,
	}
	if response, err := client.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return res, nil
}

func syncCacheStat(sourceClient influx.Client, targetClient influx.Client, statName string, days int) {
	//get records from source DB
	db := "cache_stats"
	bps, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:        db,
		Precision:       "ms",
		RetentionPolicy: "monthly",
	})

	queryString := fmt.Sprintf("select time, cdn, hostname, value from \"monthly\".\"%s\"", statName)
	if days > 0 {
		queryString += fmt.Sprintf(" where time > now() - %dd", days)
	}
	fmt.Println("queryString ", queryString)
	res, err := queryDB(sourceClient, queryString, db)
	if err != nil {
		fmt.Printf("An error occured getting %s records from sourceDb\n", statName)
		return
	}
	sourceStats := getCacheStats(res)

	//get values from target DB
	targetRes, err := queryDB(targetClient, queryString, db)
	if err != nil {
		fmt.Printf("An error occured getting %s record from target db: %v\n", statName, err)
		return
	}
	targetStats := getCacheStats(targetRes)

	for ssKey := range sourceStats {
		ts := targetStats[ssKey]
		ss := sourceStats[ssKey]
		if ts.value > ss.value {
			fmt.Printf("target value %v is at least equal to source value %v\n", ts.value, ss.value)
			continue //target value is bigger so leave it
		}
		statTime, _ := time.Parse(time.RFC3339, ss.t)
		tags := map[string]string{"cdn": ss.cdn}
		fields := map[string]interface{}{
			"value": ss.value,
		}
		pt, err := influx.NewPoint(
			statName,
			tags,
			fields,
			statTime,
		)
		if err != nil {
			fmt.Printf("error adding creating point for %v...%v\n", statName, err)
			continue
		}
		bps.AddPoint(pt)
	}
	targetClient.Write(bps)
}

func syncDeliveryServiceStat(sourceClient influx.Client, targetClient influx.Client, statName string, days int) {

	db := "deliveryservice_stats"
	bps, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:        db,
		Precision:       "ms",
		RetentionPolicy: "monthly",
	})

	queryString := fmt.Sprintf("select time, cachegroup, cdn, deliveryservice, value from \"monthly\".\"%s\"", statName)
	if days > 0 {
		queryString += fmt.Sprintf(" where time > now() - %dd", days)
	}
	fmt.Println("queryString ", queryString)
	res, err := queryDB(sourceClient, queryString, db)
	if err != nil {
		fmt.Printf("An error occured getting %s records from sourceDb: %v\n", statName, err)
		return
	}
	sourceStats := getDeliveryServiceStats(res)
	// get value from target DB
	targetRes, err := queryDB(targetClient, queryString, db)
	if err != nil {
		fmt.Printf("An error occured getting %s record from target db: %v\n", statName, err)
		return
	}
	targetStats := getDeliveryServiceStats(targetRes)

	for ssKey := range sourceStats {
		ts := targetStats[ssKey]
		ss := sourceStats[ssKey]
		if ts.value > ss.value {
			fmt.Printf("target value %v is at least equal to source value %v\n", ts.value, ss.value)
			continue //target value is bigger so leave it
		}
		statTime, _ := time.Parse(time.RFC3339, ss.t)
		tags := map[string]string{
			"cdn":             ss.cdn,
			"cachegroup":      ss.cacheGroup,
			"deliveryservice": ss.deliveryService,
		}
		fields := map[string]interface{}{
			"value": ss.value,
		}
		pt, err := influx.NewPoint(
			statName,
			tags,
			fields,
			statTime,
		)
		if err != nil {
			fmt.Printf("error adding creating point for %v...%v\n", statName, err)
			continue
		}
		bps.AddPoint(pt)
	}
	targetClient.Write(bps)
}
func syncDailyStat(sourceClient influx.Client, targetClient influx.Client, statName string, days int) {

	db := "daily_stats"
	bps, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  db,
		Precision: "s",
	})
	//get records from source DB
	queryString := fmt.Sprintf("select time, cdn, deliveryservice, value from \"%s\"", statName)
	if days > 0 {
		queryString += fmt.Sprintf(" where time > now() - %dd", days)
	}
	res, err := queryDB(sourceClient, queryString, db)
	if err != nil {
		fmt.Printf("An error occured getting %s records from sourceDb: %v\n", statName, err)
		return
	}
	sourceStats := getDailyStats(res)
	// get value from target DB
	targetRes, err := queryDB(targetClient, queryString, db)
	if err != nil {
		fmt.Printf("An error occured getting %s record from target db: %v\n", statName, err)
		return
	}
	targetStats := getDailyStats(targetRes)

	for ssKey := range sourceStats {
		ts := targetStats[ssKey]
		ss := sourceStats[ssKey]
		if ts.value >= ss.value {
			fmt.Printf("target value %v is at least equal to source value %v\n", ts.value, ss.value)
			continue //target value is bigger or equal so leave it
		}
		statTime, _ := time.Parse(time.RFC3339, ss.t)
		tags := map[string]string{
			"cdn":             ss.cdn,
			"deliveryservice": ss.deliveryService,
		}
		fields := map[string]interface{}{
			"value": ss.value,
		}
		pt, err := influx.NewPoint(
			statName,
			tags,
			fields,
			statTime,
		)
		if err != nil {
			fmt.Printf("error adding creating point for %v...%v\n", statName, err)
			continue
		}
		bps.AddPoint(pt)
	}
	targetClient.Write(bps)
}

func getCacheStats(res []influx.Result) map[string]cacheStats {
	response := make(map[string]cacheStats)
	if res != nil && len(res[0].Series) > 0 {
		for _, row := range res[0].Series {
			for _, record := range row.Values {
				data := new(cacheStats)
				t := record[0].(string)
				data.t = t
				data.cdn = record[1].(string)
				var err error
				data.value, err = record[3].(json.Number).Float64()
				if err != nil {
					fmt.Printf("Couldn't parse value from record %v\n", record)
					continue
				}
				key := data.t + data.cdn + data.hostname
				response[key] = *data
			}
		}
	}
	return response
}

func getDeliveryServiceStats(res []influx.Result) map[string]deliveryServiceStats {
	response := make(map[string]deliveryServiceStats)
	if len(res[0].Series) > 0 {
		for _, row := range res[0].Series {
			for _, record := range row.Values {
				data := new(deliveryServiceStats)
				data.t = record[0].(string)
				if record[1] != nil {
					data.cacheGroup = record[1].(string)
				}
				data.cdn = record[2].(string)
				if record[3] != nil {
					data.deliveryService = record[3].(string)
				}
				var err error
				data.value, err = record[4].(json.Number).Float64()
				if err != nil {
					fmt.Printf("Couldn't parse value from record %v\n", record)
					continue
				}
				key := data.t + data.cacheGroup + data.cdn + data.deliveryService
				response[key] = *data
			}
		}
	}
	return response
}

func getDailyStats(res []influx.Result) map[string]dailyStats {
	response := make(map[string]dailyStats)
	if len(res[0].Series) > 0 {
		for _, row := range res[0].Series {
			for _, record := range row.Values {
				data := new(dailyStats)
				data.t = record[0].(string)
				data.cdn = record[1].(string)
				data.deliveryService = record[2].(string)
				var err error
				data.value, err = record[3].(json.Number).Float64()
				if err != nil {
					fmt.Printf("Couldn't parse value from record %v\n", record)
					continue
				}
				key := data.t + data.cdn + data.deliveryService
				response[key] = *data
			}
		}
	}
	return response
}
