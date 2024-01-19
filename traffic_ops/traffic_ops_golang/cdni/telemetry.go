package cdni

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

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
)

func getTelemetries(inf *api.Info, ucdn string) (Capabilities, error) {
	capRows, err := inf.Tx.Tx.Query(CapabilityQuery, FciTelemetry, ucdn)
	if err != nil {
		return Capabilities{}, fmt.Errorf("querying capabilities: %w", err)
	}
	defer log.Close(capRows, "closing capabilities query")
	capabilities := []CapabilityQueryResponse{}
	for capRows.Next() {
		var capability CapabilityQueryResponse
		if err := capRows.Scan(&capability.Id, &capability.Type, &capability.UCdn); err != nil {
			return Capabilities{}, fmt.Errorf("scanning db rows: %w", err)
		}
		capabilities = append(capabilities, capability)
	}

	footprintMap, err := getFootprintMap(inf.Tx.Tx)
	if err != nil {
		return Capabilities{}, err
	}

	telemetryMap, err := getTelemetriesMap(inf.Tx.Tx)
	if err != nil {
		return Capabilities{}, err
	}

	telemetryMetricMap, err := getTelemetryMetricsMap(inf.Tx.Tx)
	if err != nil {
		return Capabilities{}, err
	}

	fciCaps := Capabilities{}

	for _, cap := range capabilities {
		fciCap := Capability{}
		fciCap.Footprints = footprintMap[cap.Id]
		if fciCap.Footprints == nil {
			fciCap.Footprints = []Footprint{}
		}
		returnList := []Telemetry{}
		telemetryList := telemetryMap[cap.Id]
		if telemetryList == nil {
			telemetryList = []Telemetry{}
		}

		for _, t := range telemetryList {
			telemetryMetricsList := telemetryMetricMap[t.Id]
			if telemetryMetricsList == nil {
				telemetryMetricsList = []Metric{}
			}
			t.Metrics = telemetryMetricsList
			returnList = append(returnList, t)
		}

		telemetry := Capability{
			CapabilityType: FciTelemetry,
			CapabilityValue: TelemetryCapabilityValue{
				Sources: returnList,
			},
			Footprints: fciCap.Footprints,
		}

		fciCaps.Capabilities = append(fciCaps.Capabilities, telemetry)
	}

	return fciCaps, nil
}
