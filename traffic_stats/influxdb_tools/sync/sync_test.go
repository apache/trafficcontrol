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
	"fmt"
	"testing"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSync(t *testing.T) {
	Convey("getDailyStats should work as expected", t, func() {
		Convey("getDailyStats should work with no results passed", func() {
			dailyStatsMap := getDailyStats(nil)
			So(dailyStatsMap, ShouldNotBeNil)
			So(dailyStatsMap, ShouldBeEmpty)
		})
		Convey("getDailyStats should work with empty results passed", func() {
			results := []influx.Result{}
			dailyStatsMap := getDailyStats(results)
			So(dailyStatsMap, ShouldNotBeNil)
			So(dailyStatsMap, ShouldBeEmpty)
		})
		Convey("getDailyStats should work as expected on a slice of results", func() {
			results := []influx.Result{
				influx.Result{
					Series: []models.Row{
						models.Row{
							Values: generateDailyValues(0),
						},
					},
				},
				influx.Result{
					Series: []models.Row{
						models.Row{
							Values: generateDailyValues(1),
						},
						models.Row{
							Values: generateDailyValues(2),
						},
					},
				},
				influx.Result{
					Series: []models.Row{
						models.Row{
							Values: generateDailyValues(3),
						},
						models.Row{
							Values: generateDailyValues(4),
						},
						models.Row{
							Values: generateDailyValues(5),
						},
					},
				},
			}
			dailyStatsMap := getDailyStats(results)
			So(dailyStatsMap, ShouldNotBeNil)
			So(dailyStatsMap, ShouldNotBeEmpty)
			So(len(dailyStatsMap), ShouldEqual, 6) // we get one dailyStats object per Values entry
		})
	})
	Convey("getDeliveryServiceStats should work as expected", t, func() {
		Convey("getDeliveryServiceStats should work with no results passed", func() {
			getDeliveryServiceStatsMap := getDeliveryServiceStats(nil)
			So(getDeliveryServiceStatsMap, ShouldNotBeNil)
			So(getDeliveryServiceStatsMap, ShouldBeEmpty)
		})
		Convey("getDeliveryServiceStats should work with empty results passed", func() {
			results := []influx.Result{}
			getDeliveryServiceStatsMap := getDeliveryServiceStats(results)
			So(getDeliveryServiceStatsMap, ShouldNotBeNil)
			So(getDeliveryServiceStatsMap, ShouldBeEmpty)
		})
		Convey("getDeliveryServiceStats should work as expected on a slice of results", func() {
			results := []influx.Result{
				influx.Result{
					Series: []models.Row{
						models.Row{
							Values: generateDeliveryServiceValues(0),
						},
					},
				},
				influx.Result{
					Series: []models.Row{
						models.Row{
							Values: generateDeliveryServiceValues(1),
						},
						models.Row{
							Values: generateDeliveryServiceValues(2),
						},
					},
				},
				influx.Result{
					Series: []models.Row{
						models.Row{
							Values: generateDeliveryServiceValues(3),
						},
						models.Row{
							Values: generateDeliveryServiceValues(4),
						},
						models.Row{
							Values: generateDeliveryServiceValues(5),
						},
					},
				},
			}
			getDeliveryServiceStatsMap := getDeliveryServiceStats(results)
			So(getDeliveryServiceStatsMap, ShouldNotBeNil)
			So(getDeliveryServiceStatsMap, ShouldNotBeEmpty)
			So(len(getDeliveryServiceStatsMap), ShouldEqual, 6) // we get one dailyStats object per Values entry
		})
	})
	Convey("getCacheStats should work as expected", t, func() {
		Convey("getCacheStats should work with no results passed", func() {
			getCacheStatsMap := getCacheStats(nil)
			So(getCacheStatsMap, ShouldNotBeNil)
			So(getCacheStatsMap, ShouldBeEmpty)
		})
		Convey("getCacheStats should work with empty results passed", func() {
			results := []influx.Result{}
			getCacheStatsMap := getCacheStats(results)
			So(getCacheStatsMap, ShouldNotBeNil)
			So(getCacheStatsMap, ShouldBeEmpty)
		})
		Convey("getCacheStats should work as expected on a slice of results", func() {
			results := []influx.Result{
				influx.Result{
					Series: []models.Row{
						models.Row{
							Values: generateCacheValues(0),
						},
					},
				},
				influx.Result{
					Series: []models.Row{
						models.Row{
							Values: generateCacheValues(1),
						},
						models.Row{
							Values: generateCacheValues(2),
						},
					},
				},
				influx.Result{
					Series: []models.Row{
						models.Row{
							Values: generateCacheValues(3),
						},
						models.Row{
							Values: generateCacheValues(4),
						},
						models.Row{
							Values: generateCacheValues(5),
						},
					},
				},
			}
			getCacheStatsMap := getCacheStats(results)
			So(getCacheStatsMap, ShouldNotBeNil)
			So(getCacheStatsMap, ShouldNotBeEmpty)
			So(len(getCacheStatsMap), ShouldEqual, 6) // we get one dailyStats object per Values entry
		})
	})
}

func generateDailyValues(i int) [][]interface{} {
	ret := make([][]interface{}, 1)
	startIdx := i * 4 // this is to have the right amount of difference between the test numbers used
	for j := startIdx; j < startIdx+3; j++ {
		n := startIdx + j
		ret[0] = make([]interface{}, 4)
		ret[0][0] = fmt.Sprintf("t%d", n)
		ret[0][1] = fmt.Sprintf("test.cdn-%d", n)
		ret[0][2] = fmt.Sprintf("testDeliveryService-%d", n)
		ret[0][3] = json.Number(fmt.Sprintf("%d.%d%d", n, n, n))
	}
	return ret
}

func generateDeliveryServiceValues(i int) [][]interface{} {
	ret := make([][]interface{}, 1)
	startIdx := i * 4 // this is to have the right amount of difference between the test numbers used
	for j := startIdx; j < startIdx+3; j++ {
		n := startIdx + j
		ret[0] = make([]interface{}, 5)
		ret[0][0] = fmt.Sprintf("t%d", n)
		ret[0][1] = fmt.Sprintf("cache-group%d", n)
		ret[0][2] = fmt.Sprintf("test.cdn-%d", n)
		ret[0][3] = fmt.Sprintf("testDeliveryService-%d", n)
		ret[0][4] = json.Number(fmt.Sprintf("%d.%d%d", n, n, n))
	}
	return ret
}

func generateCacheValues(i int) [][]interface{} {
	ret := make([][]interface{}, 1)
	startIdx := i * 4 // this is to have the right amount of difference between the test numbers used
	for j := startIdx; j < startIdx+3; j++ {
		n := startIdx + j
		ret[0] = make([]interface{}, 5)
		ret[0][0] = fmt.Sprintf("t%d", n)
		ret[0][1] = fmt.Sprintf("test.cdn-%d", n)
		ret[0][2] = fmt.Sprintf("test.hostname-%d", n)
		ret[0][3] = fmt.Sprintf("test.cacheType-%d", n)
		ret[0][4] = json.Number(fmt.Sprintf("%d.%d%d", n, n, n))
	}
	return ret
}
