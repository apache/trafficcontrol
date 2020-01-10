package cache

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
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc/enum"
	"io/ioutil"
	"math/rand"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/traffic_monitor/todata"

	"github.com/json-iterator/go"
)

func TestAstats(t *testing.T) {
	t.Log("Running Astats Tests")

	text, err := ioutil.ReadFile("astats.json")
	if err != nil {
		t.Log(err)
	}
	aStats := Astats{}
	json := jsoniter.ConfigFastest
	err = json.Unmarshal(text, &aStats)
	fmt.Printf("aStats ---> %v\n", aStats)
	if err != nil {
		t.Log(err)
	}
	fmt.Printf("Found %v key/val pairs in ats\n", len(aStats.Ats))
}

func getMockTODataDSNameDirectMatches() map[enum.DeliveryServiceName]string {
	return map[enum.DeliveryServiceName]string{
		"ds0": "ds0.example.invalid",
		"ds1": "ds1.example.invalid",
	}
}

func getMockTOData(dsNameFQDNs map[enum.DeliveryServiceName]string) todata.TOData {
	tod := todata.New()
	for dsName, dsDirectMatch := range dsNameFQDNs {
		tod.DeliveryServiceRegexes.DirectMatches[dsDirectMatch] = dsName
	}
	return *tod
}

func getMockRawStats(cacheName string, dsNameFQDNs map[enum.DeliveryServiceName]string) map[string]interface{} {
	st := map[string]interface{}{}
	for _, dsFQDN := range dsNameFQDNs {
		st["plugin.remap_stats."+dsFQDN+".in_bytes"] = float64(rand.Uint64())
		st["plugin.remap_stats."+dsFQDN+".out_bytes"] = float64(rand.Uint64())
		st["plugin.remap_stats."+dsFQDN+".status_2xx"] = float64(rand.Uint64())
		st["plugin.remap_stats."+dsFQDN+".status_3xx"] = float64(rand.Uint64())
		st["plugin.remap_stats."+dsFQDN+".status_4xx"] = float64(rand.Uint64())
		st["plugin.remap_stats."+dsFQDN+".status_5xx"] = float64(rand.Uint64())
	}
	return st
}

func getMockSystem(infSpeed int, outBytes int) AstatsSystem {
	infName := randStr()
	return AstatsSystem{
		InfName:           infName,
		InfSpeed:          9876554433210,
		ProcNetDev:        infName + ":12234567 8901234    1    2    3     4          5   12345 " + strconv.Itoa(outBytes) + " 923412341234    6    7    8     9       10          11",
		ProcLoadavg:       "1.2 2.34 5.67 1/876 1234",
		ConfigLoadRequest: rand.Int(),
		LastReloadRequest: rand.Int(),
		ConfigReloads:     rand.Int(),
		LastReload:        rand.Int(),
		AstatsLoad:        rand.Int(),
		NotAvailable:      randBool(),
	}

}

func TestAstatsPrecompute(t *testing.T) {
	dsNameFQDNs := getMockTODataDSNameDirectMatches()
	toData := getMockTOData(dsNameFQDNs)
	cacheName := "cache0"
	rawStats := getMockRawStats(cacheName, dsNameFQDNs)
	outBytes := 987655443321
	infSpeedMbps := 9876554433210
	system := getMockSystem(infSpeedMbps, outBytes)

	prc := astatsPrecompute(enum.CacheName(cacheName), toData, rawStats, system)

	if len(prc.Errors) != 0 {
		t.Fatalf("astatsPrecompute Errors expected 0, actual: %+v\n", prc.Errors)
	}
	if prc.OutBytes != int64(outBytes) {
		t.Fatalf("astatsPrecompute OutBytes expected 987655443321, actual: %+v\n", prc.OutBytes)
	}
	if prc.MaxKbps != int64(infSpeedMbps*1000) {
		t.Fatalf("astatsPrecompute MaxKbps expected 9876554433210000, actual: %+v\n", prc.MaxKbps)
	}

	for dsName, dsFQDN := range dsNameFQDNs {
		dsStat, ok := prc.DeliveryServiceStats[dsName]
		if !ok {
			t.Fatalf("astatsPrecompute DeliveryServiceStats expected %+v, actual: missing\n", dsName)
		}
		if statName := "plugin.remap_stats." + dsFQDN + ".in_bytes"; dsStat.InBytes != uint64(rawStats[statName].(float64)) {
			t.Fatalf("astatsPrecompute DeliveryServiceStats[%+v].InBytes expected %+v, actual: %+v\n", dsName, uint64(rawStats[statName].(float64)), dsStat.InBytes)
		}
		if statName := "plugin.remap_stats." + dsFQDN + ".out_bytes"; dsStat.OutBytes != uint64(rawStats[statName].(float64)) {
			t.Fatalf("astatsPrecompute DeliveryServiceStats[%+v].OutBytes expected %+v, actual: %+v\n", dsName, uint64(rawStats[statName].(float64)), dsStat.OutBytes)
		}
		if statName := "plugin.remap_stats." + dsFQDN + ".status_2xx"; dsStat.Status2xx != uint64(rawStats[statName].(float64)) {
			t.Fatalf("astatsPrecompute DeliveryServiceStats[%+v].Status2xx expected %+v, actual: %+v\n", dsName, uint64(rawStats[statName].(float64)), dsStat.Status2xx)
		}
		if statName := "plugin.remap_stats." + dsFQDN + ".status_3xx"; dsStat.Status3xx != uint64(rawStats[statName].(float64)) {
			t.Fatalf("astatsPrecompute DeliveryServiceStats[%+v].Status3xx expected %+v, actual: %+v\n", dsName, uint64(rawStats[statName].(float64)), dsStat.Status3xx)
		}
		if statName := "plugin.remap_stats." + dsFQDN + ".status_4xx"; dsStat.Status4xx != uint64(rawStats[statName].(float64)) {
			t.Fatalf("astatsPrecompute DeliveryServiceStats[%+v].Status4xx expected %+v, actual: %+v\n", dsName, uint64(rawStats[statName].(float64)), dsStat.Status4xx)
		}
		if statName := "plugin.remap_stats." + dsFQDN + ".status_5xx"; dsStat.Status5xx != uint64(rawStats[statName].(float64)) {
			t.Fatalf("astatsPrecompute DeliveryServiceStats[%+v].Status5xx expected %+v, actual: %+v\n", dsName, uint64(rawStats[statName].(float64)), dsStat.Status5xx)
		}
	}
}
