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

	"github.com/lib/pq"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

const totalLimitsQuery = `SELECT limit_type, maximum_hard, maximum_soft, ctl.telemetry_id, ctl.telemetry_metric, t.id, t.type, tm.name FROM cdni_total_limits AS ctl LEFT JOIN cdni_telemetry as t ON telemetry_id = t.id LEFT JOIN cdni_telemetry_metrics as tm ON telemetry_metric = tm.name WHERE ctl.capability_id = $1`
const hostLimitsQuery = `SELECT limit_type, maximum_hard, maximum_soft, chl.telemetry_id, chl.telemetry_metric, t.id, t.type, tm.name, host FROM cdni_host_limits AS chl LEFT JOIN cdni_telemetry as t ON telemetry_id = t.id LEFT JOIN cdni_telemetry_metrics as tm ON telemetry_metric = tm.name WHERE chl.capability_id = $1 ORDER BY host DESC`

func GetCapacities(inf *api.APIInfo, ucdn string) (Capabilities, error) {
	capRows, err := inf.Tx.Tx.Query(CapabilityQuery, FciCapacityLimits, ucdn)
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

		tlRows, err := inf.Tx.Tx.Query(totalLimitsQuery, cap.Id)
		if err != nil {
			return Capabilities{}, fmt.Errorf("querying total limits: %w", err)
		}

		defer tlRows.Close()
		totalLimits := []TotalLimitsQueryResponse{}
		for tlRows.Next() {
			var totalLimit TotalLimitsQueryResponse
			if err := tlRows.Scan(&totalLimit.LimitType, &totalLimit.MaximumHard, &totalLimit.MaximumSoft, &totalLimit.TelemetryId, &totalLimit.TelemetryMetic, &totalLimit.Id, &totalLimit.Type, &totalLimit.Name); err != nil {
				return Capabilities{}, fmt.Errorf("scanning db rows: %w", err)
			}
			totalLimits = append(totalLimits, totalLimit)
		}

		hlRows, err := inf.Tx.Tx.Query(hostLimitsQuery, cap.Id)
		if err != nil {
			return Capabilities{}, fmt.Errorf("querying host limits: %w", err)
		}

		defer hlRows.Close()
		hostLimits := []HostLimitsResponse{}
		for hlRows.Next() {
			var hostLimit HostLimitsResponse
			if err := hlRows.Scan(&hostLimit.LimitType, &hostLimit.MaximumHard, &hostLimit.MaximumSoft, &hostLimit.TelemetryId, &hostLimit.TelemetryMetic, &hostLimit.Id, &hostLimit.Type, &hostLimit.Name, &hostLimit.Host); err != nil {
				return Capabilities{}, fmt.Errorf("scanning db rows: %w", err)
			}
			hostLimits = append(hostLimits, hostLimit)
		}

		returnedTotalLimits := []Limit{}
		for _, tl := range totalLimits {
			returnedTotalLimit := Limit{
				LimitType:   CapacityLimitType(tl.LimitType),
				MaximumHard: tl.MaximumHard,
				MaximumSoft: tl.MaximumSoft,
				TelemetrySource: TelemetrySource{
					Id:     tl.TelemetryId,
					Metric: tl.TelemetryMetic,
				},
			}
			returnedTotalLimits = append(returnedTotalLimits, returnedTotalLimit)
		}

		returnedHostLimits := []HostLimit{}
		hostToLimitMap := map[string][]Limit{}
		for _, hl := range hostLimits {
			limit := Limit{
				LimitType:   CapacityLimitType(hl.LimitType),
				MaximumHard: hl.MaximumHard,
				MaximumSoft: hl.MaximumSoft,
				TelemetrySource: TelemetrySource{
					Id:     hl.TelemetryId,
					Metric: hl.TelemetryMetic,
				},
			}

			if val, ok := hostToLimitMap[hl.Host]; ok {
				val = append(val, limit)
				hostToLimitMap[hl.Host] = val
			} else {
				hlList := []Limit{}
				hlList = append(hlList, limit)
				hostToLimitMap[hl.Host] = hlList
			}
		}

		for h, l := range hostToLimitMap {
			returnedHostLimit := HostLimit{
				Host:   h,
				Limits: l,
			}
			returnedHostLimits = append(returnedHostLimits, returnedHostLimit)
		}

		fciCap.CapabilityType = FciCapacityLimits
		fciCap.CapabilityValue = []CapacityCapabilityValue{
			{
				TotalLimits: returnedTotalLimits,
				HostLimits:  returnedHostLimits,
			},
		}

		fciCaps.Capabilities = append(fciCaps.Capabilities, fciCap)
	}

	return fciCaps, nil
}

type CapabilityQueryResponse struct {
	Id   int64  `json:"id" db:"id"`
	Type string `json:"type" db:"type"`
	UCdn string `json:"ucdn" db:"ucdn"`
}

type TotalLimitsQueryResponse struct {
	LimitType      string `json:"limit_type" db:"limit_type"`
	MaximumHard    int64  `json:"maximum_hard" db:"maximum_hard"`
	MaximumSoft    int64  `json:"maximum_soft" db:"maximum_soft"`
	TelemetryId    string `json:"telemetry_id" db:"telemetry_id"`
	TelemetryMetic string `json:"telemetry_metric" db:"telemetry_metric"`
	UCdn           string `json:"ucdn" db:"ucdn"`
	Id             string `json:"id" db:"id"`
	Type           string `json:"type" db:"type"`
	Name           string `json:"name" db:"name"`
}
type HostLimitsResponse struct {
	Host string `json:"host" db:"host"`
	TotalLimitsQueryResponse
}
