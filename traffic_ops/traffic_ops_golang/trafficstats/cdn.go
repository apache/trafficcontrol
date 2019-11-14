package trafficstats

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"errors"
	"fmt"
	"net/http"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/lib/pq"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

const (
	cdnStatsQuery = `
SELECT last(value) FROM "%s"."monthly"."%s"
	WHERE cdn = $cdn`
)

// GetCurrentStats handler for getting current stats for CDNs
func GetCurrentStats(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx

	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	client, err := inf.CreateInfluxClient()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	} else if client == nil {
		sysErr = errors.New("Traffic Stats is not configured, but DS stats were requested")
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}
	defer (*client).Close()

	currentStats := []interface{}{}

	// Get CDN names
	cdns := []string{}
	if err := tx.QueryRow(`SELECT ARRAY(SELECT name FROM cdn)`).Scan(pq.Array(&cdns)); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("querying cdn names"))
		return
	}
	totalStats := tc.TrafficStatsTotalStats{
		CDN: "total",
	}
	// util func to add to total
	addToTotal := func(v, t *float64) *float64 {
		if v == nil {
			return t
		}
		if t == nil {
			t = util.FloatPtr(0.0)
		}
		*t += *v
		return t
	}

	for _, cdn := range cdns {
		cdnStats := tc.TrafficStatsCDNStats{
			CDN: cdn,
		}
		bw, err := getCDNStat(client, cdn, "bandwidth.cdn.1min", inf.Config.ConfigInflux.CacheDBName)
		if err != nil {
			sysErr = fmt.Errorf("getting bandwidth from cdn %v: %v", cdn, err)
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
			return
		}
		totalStats.Bandwidth = addToTotal(bw, totalStats.Bandwidth)

		if bw != nil {
			cdnStats.Bandwidth = util.FloatPtr(*bw / 1000000)
		}

		con, err := getCDNStat(client, cdn, "connections.cdn.1min", inf.Config.ConfigInflux.CacheDBName)
		if err != nil {
			sysErr = fmt.Errorf("getting connections from cdn %v: %v", cdn, err)
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
			return
		}
		totalStats.Connnections = addToTotal(con, totalStats.Connnections)
		cdnStats.Connnections = con

		cap, err := getCDNStat(client, cdn, "maxkbps.cdn.1min", inf.Config.ConfigInflux.CacheDBName)
		if err != nil {
			sysErr = fmt.Errorf("getting maxkbps from cdn %v: %v", cdn, err)
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
			return
		}
		if cap != nil {
			cdnStats.Capacity = util.FloatPtr(*cap / 1000000 * .85)
		}

		currentStats = append(currentStats, cdnStats)
	}

	if totalStats.Bandwidth != nil {
		*totalStats.Bandwidth /= 1000000
	}
	currentStats = append(currentStats, totalStats)

	api.WriteResp(w, r, currentStats)
}

func getCDNStat(client *influx.Client, cdnName, metricName, db string) (*float64, error) {
	qStr := fmt.Sprintf(cdnStatsQuery, db, metricName)
	q := influx.NewQueryWithParameters(qStr,
		db,
		"rfc3339",
		map[string]interface{}{
			"cdn": cdnName,
		})
	series, err := getSeries(db, q, client)
	if err != nil {
		return nil, err
	}
	if series == nil {
		return nil, nil
	}
	vals := series.Values[0]
	mappedValues := map[string]interface{}{}
	for i, v := range vals {
		mappedValues[series.Columns[i]] = v
	}
	val, err := extractFloat64("last", mappedValues)
	if err != nil {
		return nil, err
	}
	return &val, nil
}
