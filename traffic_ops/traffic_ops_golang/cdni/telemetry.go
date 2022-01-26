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
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func GetTelemetries(inf *api.APIInfo, ucdn string) (Capabilities, error) {
	capRows, err := inf.Tx.Tx.Query(CapabilityQuery, FciTelemetry, ucdn)
	if err != nil {
		return Capabilities{}, fmt.Errorf("querying capabilities: %w", err)
	}
	defer capRows.Close()
	capabilities := []CapabilityQueryResponse{}
	for capRows.Next() {
		var capability CapabilityQueryResponse
		if err := capRows.Scan(&capability.Id, &capability.Type, &capability.UCdn); err != nil {
			return Capabilities{}, fmt.Errorf("scanning db rows: %w", err)
		}
		capabilities = append(capabilities, capability)
	}

	fciCaps := Capabilities{}

	for _, cap := range capabilities {
		fciCap := Capability{}
		footRows, err := inf.Tx.Tx.Query(FootprintQuery, cap.Id)
		if err != nil {
			return Capabilities{}, fmt.Errorf("querying footprints: %w", err)
		}
		defer footRows.Close()
		footprints := []Footprint{}
		for footRows.Next() {
			var footprint Footprint
			if err := footRows.Scan(&footprint.FootprintType, pq.Array(&footprint.FootprintValue)); err != nil {
				return Capabilities{}, fmt.Errorf("scanning db rows: %w", err)
			}
			footprints = append(footprints, footprint)
		}

		fciCap.Footprints = footprints

		rows, err := inf.Tx.Tx.Query(`SELECT id, type FROM cdni_telemetry WHERE capability_id = $1`, cap.Id)
		if err != nil {
			return Capabilities{}, errors.New("querying cdni telemetry: " + err.Error())
		}
		defer rows.Close()
		returnList := []Telemetry{}
		telemetryList := []Telemetry{}
		for rows.Next() {
			telemetry := Telemetry{}
			if err := rows.Scan(&telemetry.Id, &telemetry.Type); err != nil {
				return Capabilities{}, errors.New("scanning telemetry: " + err.Error())
			}
			telemetryList = append(telemetryList, telemetry)
		}

		for _, t := range telemetryList {
			tmRows, err := inf.Tx.Tx.Query(`SELECT name, time_granularity, data_percentile, latency FROM cdni_telemetry_metrics WHERE telemetry_id = $1`, t.Id)
			if err != nil {
				return Capabilities{}, errors.New("querying cdni telemetry metrics: " + err.Error())
			}
			defer tmRows.Close()
			telemetryMetricsList := []Metric{}
			for tmRows.Next() {
				metric := Metric{}
				if err := tmRows.Scan(&metric.Name, &metric.TimeGranularity, &metric.DataPercentile, &metric.Latency); err != nil {
					return Capabilities{}, errors.New("scanning telemetry metric: " + err.Error())
				}
				telemetryMetricsList = append(telemetryMetricsList, metric)
			}

			t.Metrics = telemetryMetricsList
			returnList = append(returnList, t)
		}

		telemetry := Capability{
			CapabilityType: FciTelemetry,
			CapabilityValue: TelemetryCapabilityValue{
				Sources: returnList,
			},
			Footprints: footprints,
		}

		fciCaps.Capabilities = append(fciCaps.Capabilities, telemetry)
	}

	return fciCaps, nil
}
