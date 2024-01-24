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

func getCapacities(inf *api.Info, ucdn string) (Capabilities, error) {
	capRows, err := inf.Tx.Tx.Query(CapabilityQuery, FciCapacityLimits, ucdn)
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

	limitsMap, err := getLimitsMap(inf.Tx.Tx)
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
		limits := limitsMap[cap.Id]
		if limits == nil {
			limits = []LimitsQueryResponse{}
		}

		returnedLimits := []Limit{}
		for _, l := range limits {
			returnedTotalLimit := Limit{
				Id:          l.LimitId,
				LimitType:   CapacityLimitType(l.LimitType),
				MaximumHard: l.MaximumHard,
				MaximumSoft: l.MaximumSoft,
				TelemetrySource: TelemetrySource{
					Id:     l.TelemetryId,
					Metric: l.TelemetryMetic,
				},
			}

			returnedTotalLimit.Scope = l.Scope

			returnedLimits = append(returnedLimits, returnedTotalLimit)
		}

		fciCap.CapabilityType = FciCapacityLimits
		fciCap.CapabilityValue = []CapacityCapabilityValue{
			{
				Limits: returnedLimits,
			},
		}

		fciCaps.Capabilities = append(fciCaps.Capabilities, fciCap)
	}

	return fciCaps, nil
}

// CapabilityQueryResponse contains data about the capability query.
type CapabilityQueryResponse struct {
	Id   int    `json:"id" db:"id"`
	Type string `json:"type" db:"type"`
	UCdn string `json:"ucdn" db:"ucdn"`
}

// LimitsQueryResponse contains information about the limits query.
type LimitsQueryResponse struct {
	Scope          *LimitScope `json:"scope,omitempty"`
	LimitId        string      `json:"limitId" db:"limit_id"`
	LimitType      string      `json:"limitType" db:"limit_type"`
	MaximumHard    int64       `json:"maximum_hard" db:"maximum_hard"`
	MaximumSoft    int64       `json:"maximum_soft" db:"maximum_soft"`
	TelemetryId    string      `json:"telemetry_id" db:"telemetry_id"`
	TelemetryMetic string      `json:"telemetry_metric" db:"telemetry_metric"`
	UCdn           string      `json:"ucdn" db:"ucdn"`
	Id             string      `json:"id" db:"id"`
	Type           string      `json:"type" db:"type"`
	Name           string      `json:"name" db:"name"`
	CapabilityId   int         `json:"-"`
}

// LimitScope contains information for a specific limit.
type LimitScope struct {
	ScopeType  *string  `json:"type" db:"scope_type"`
	ScopeValue []string `json:"value" db:"scope_value"`
}
