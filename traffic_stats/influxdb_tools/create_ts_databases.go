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
	"fmt"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func createFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Usage: "The influxdb url and port",
			Value: "http://localhost:8086",
		},
		cli.IntFlag{
			Name:  "replication",
			Usage: "The number of nodes in the cluster",
			Value: 3,
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "The influxdb username used to create DBs",
			Value: "",
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "The influxdb password used to create DBs",
			Value: "",
		},
	}
}

func create(c *cli.Context) error {
	influxURL := c.String("url")
	replication := c.Int("replication")
	user := c.String("user")
	password := c.String("password")

	fmt.Printf("creating datbases for influxUrl: %s with a replication of %d using user %s\n", influxURL, replication, user)
	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     influxURL,
		Username: user,
		Password: password,
	})
	if err != nil {
		return errors.Wrap(err, "Error creating influx client")
	}
	_, _, err = client.Ping(10)
	if err != nil {
		return errors.Wrap(err, "Error creating influx client")
	}

	createCacheStats(client, replication)
	createDailyStats(client, replication)
	createDeliveryServiceStats(client, replication)
	return nil
}

// queryDB takes a variadic argument for the target database so as to make
// passing the variable optional, however, if passed, only the first db passed
// in will be used
func queryDB(client influx.Client, cmd string, dbs ...string) (res []influx.Result, err error) {
	db := ""
	if len(dbs) > 0 {
		db = dbs[0]
	}
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

func createCacheStats(client influx.Client, replication int) {
	db := cache
	createDatabase(client, db)
	createRetentionPolicy(client, db, "daily", "26h", replication, true)
	createRetentionPolicy(client, db, "monthly", "30d", replication, false)
	createRetentionPolicy(client, db, "indefinite", "INF", replication, false)
	createContinuousQuery(client, "bandwidth_1min", `CREATE CONTINUOUS QUERY bandwidth_1min ON cache_stats RESAMPLE FOR 2m BEGIN SELECT mean(value) AS "value" INTO "cache_stats"."monthly"."bandwidth.1min" FROM "cache_stats"."daily".bandwidth GROUP BY time(1m), * END`)
	createContinuousQuery(client, "connections_1min", `CREATE CONTINUOUS QUERY connections_1min ON cache_stats RESAMPLE FOR 2m BEGIN SELECT mean(value) AS "value" INTO "cache_stats"."monthly"."connections.1min" FROM "cache_stats"."daily"."ats.proxy.process.http.current_client_connections" GROUP BY time(1m), * END`)
	createContinuousQuery(client, "bandwidth_cdn_1min", `CREATE CONTINUOUS QUERY bandwidth_cdn_1min ON cache_stats RESAMPLE FOR 5m BEGIN SELECT sum(value) AS "value" INTO "cache_stats"."monthly"."bandwidth.cdn.1min" FROM "cache_stats"."monthly"."bandwidth.1min" GROUP BY time(1m), cdn END`)
	createContinuousQuery(client, "connections_cdn_1min", `CREATE CONTINUOUS QUERY connections_cdn_1min ON cache_stats RESAMPLE FOR 5m BEGIN SELECT sum(value) AS "value" INTO "cache_stats"."monthly"."connections.cdn.1min" FROM "cache_stats"."monthly"."connections.1min" GROUP BY time(1m), cdn END`)
	createContinuousQuery(client, "bandwidth_cdn_type_1min", `CREATE CONTINUOUS QUERY bandwidth_cdn_type_1min ON cache_stats RESAMPLE FOR 5m BEGIN SELECT sum(value) AS "value" INTO "cache_stats"."monthly"."bandwidth.cdn.type.1min" FROM "cache_stats"."monthly"."bandwidth.1min" GROUP BY time(1m), cdn, type END`)
	createContinuousQuery(client, "connections_cdn_type_1min", `CREATE CONTINUOUS QUERY connections_cdn_type_1min ON cache_stats RESAMPLE FOR 5m BEGIN SELECT sum(value) AS "value" INTO "cache_stats"."monthly"."connections.cdn.type.1min" FROM "cache_stats"."monthly"."connections.1min" GROUP BY time(1m), cdn, type END`)
	createContinuousQuery(client, "maxKbps_1min", `CREATE CONTINUOUS QUERY maxKbps_1min ON cache_stats RESAMPLE FOR 2m BEGIN SELECT mean(value) AS value INTO cache_stats.monthly."maxkbps.1min" FROM cache_stats.daily.maxKbps GROUP BY time(1m), * END`)
	createContinuousQuery(client, "maxkbps_cdn_1min", `CREATE CONTINUOUS QUERY maxkbps_cdn_1min ON cache_stats RESAMPLE FOR 5m BEGIN SELECT sum(value) AS value INTO cache_stats.monthly."maxkbps.cdn.1min" FROM cache_stats.monthly."maxkbps.1min" GROUP BY time(1m), cdn END`)
	createContinuousQuery(client, "wrap_count_vol1_1m", `CREATE CONTINUOUS QUERY wrap_count_vol1_1m ON cache_stats RESAMPLE FOR 2m BEGIN SELECT mean(value) AS vol1_wrap_count INTO cache_stats.monthly."wrap_count.1min" FROM cache_stats.daily."ats.proxy.process.cache.volume_1.wrap_count" GROUP BY time(1m), * END`)
	createContinuousQuery(client, "wrap_count_vol2_1m", `CREATE CONTINUOUS QUERY wrap_count_vol2_1m ON cache_stats RESAMPLE FOR 2m BEGIN SELECT mean(value) AS vol2_wrap_count INTO cache_stats.monthly."wrap_count.1min" FROM cache_stats.daily."ats.proxy.process.cache.volume_2.wrap_count" GROUP BY time(1m), * END`)
}

func createDeliveryServiceStats(client influx.Client, replication int) {
	db := deliveryService
	createDatabase(client, db)
	createRetentionPolicy(client, db, "daily", "26h", replication, true)
	createRetentionPolicy(client, db, "monthly", "30d", replication, false)
	createRetentionPolicy(client, db, "indefinite", "INF", replication, false)
	createContinuousQuery(client, "tps_2xx_ds_1min", `CREATE CONTINUOUS QUERY tps_2xx_ds_1min ON deliveryservice_stats RESAMPLE FOR 2m BEGIN SELECT mean(value) AS "value" INTO "deliveryservice_stats"."monthly"."tps_2xx.ds.1min" FROM "deliveryservice_stats"."daily".tps_2xx WHERE cachegroup = 'total' GROUP BY time(1m), * END`)
	createContinuousQuery(client, "tps_3xx_ds_1min", `CREATE CONTINUOUS QUERY tps_3xx_ds_1min ON deliveryservice_stats RESAMPLE FOR 2m BEGIN SELECT mean(value) AS "value" INTO "deliveryservice_stats"."monthly"."tps_3xx.ds.1min" FROM "deliveryservice_stats"."daily".tps_3xx WHERE cachegroup = 'total' GROUP BY time(1m), * END`)
	createContinuousQuery(client, "tps_4xx_ds_1min", `CREATE CONTINUOUS QUERY tps_4xx_ds_1min ON deliveryservice_stats RESAMPLE FOR 2m BEGIN SELECT mean(value) AS "value" INTO "deliveryservice_stats"."monthly"."tps_4xx.ds.1min" FROM "deliveryservice_stats"."daily".tps_4xx WHERE cachegroup = 'total' GROUP BY time(1m), * END`)
	createContinuousQuery(client, "tps_5xx_ds_1min", `CREATE CONTINUOUS QUERY tps_5xx_ds_1min ON deliveryservice_stats RESAMPLE FOR 2m BEGIN SELECT mean(value) AS "value" INTO "deliveryservice_stats"."monthly"."tps_5xx.ds.1min" FROM "deliveryservice_stats"."daily".tps_5xx WHERE cachegroup = 'total' GROUP BY time(1m), * END`)
	createContinuousQuery(client, "tps_total_ds_1min", `CREATE CONTINUOUS QUERY tps_total_ds_1min ON deliveryservice_stats RESAMPLE FOR 2m BEGIN SELECT mean(value) AS "value" INTO "deliveryservice_stats"."monthly"."tps_total.ds.1min" FROM "deliveryservice_stats"."daily".tps_total WHERE cachegroup = 'total' GROUP BY time(1m), * END`)
	createContinuousQuery(client, "kbps_ds_1min", `CREATE CONTINUOUS QUERY kbps_ds_1min ON deliveryservice_stats RESAMPLE FOR 2m BEGIN SELECT mean(value) AS "value" INTO "deliveryservice_stats"."monthly"."kbps.ds.1min" FROM "deliveryservice_stats"."daily".kbps WHERE cachegroup = 'total' GROUP BY time(1m), * END`)
	createContinuousQuery(client, "kbps_cg_1min", `CREATE CONTINUOUS QUERY kbps_cg_1min ON deliveryservice_stats RESAMPLE FOR 2m BEGIN SELECT mean(value) AS "value" INTO "deliveryservice_stats"."monthly"."kbps.cg.1min" FROM "deliveryservice_stats"."daily".kbps WHERE cachegroup != 'total' GROUP BY time(1m), * END`)
	createContinuousQuery(client, "max_kbps_ds_1day", `CREATE CONTINUOUS QUERY max_kbps_ds_1day ON deliveryservice_stats RESAMPLE FOR 2d BEGIN SELECT max(value) AS "value" INTO "deliveryservice_stats"."indefinite"."max.kbps.ds.1day" FROM "deliveryservice_stats"."monthly"."kbps.ds.1min" GROUP BY time(1d), deliveryservice, cdn END`)
}

func createDailyStats(client influx.Client, replication int) {
	db := daily
	createDatabase(client, db)
	createRetentionPolicy(client, db, "indefinite", "INF", replication, true)
}

func createDatabase(client influx.Client, db string) {
	_, err := queryDB(client, fmt.Sprintf("CREATE DATABASE %s", db))
	if err != nil {
		fmt.Printf("An error occured creating the %v database: %v\n", db, err)
		return
	}
	fmt.Println("Successfully created database: ", db)
}

func createRetentionPolicy(client influx.Client, db string, name string, duration string, replication int, isDefault bool) {
	qString := fmt.Sprintf("CREATE RETENTION POLICY %s ON %s DURATION %s REPLICATION %d", name, db, duration, replication)
	if isDefault {
		qString += " DEFAULT"
	}
	_, err := queryDB(client, qString)
	if err != nil {
		fmt.Printf("An error occured creating the retention policy %s on database: %s:  %v\n", name, db, err)
		return
	}
	fmt.Printf("Successfully created retention policy %s for database: %s\n", name, db)
}

func createContinuousQuery(client influx.Client, name string, query string) {
	_, err := queryDB(client, query)
	if err != nil {
		fmt.Printf("An error occured creating continuous query %s: %v\n", name, err)
		return
	}
	fmt.Println("Successfully created continuous query ", name)
}
