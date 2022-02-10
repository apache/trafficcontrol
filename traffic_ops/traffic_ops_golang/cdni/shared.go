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
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"net/http"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"

	"github.com/dgrijalva/jwt-go"
)

const CapabilityQuery = `SELECT id, type, ucdn FROM cdni_capabilities WHERE type = $1 AND ucdn = $2`
const AllFootprintQuery = `SELECT footprint_type, footprint_value::text[], capability_id FROM cdni_footprints`

const totalLimitsQuery = `
SELECT limit_type, maximum_hard, maximum_soft, ctl.telemetry_id, ctl.telemetry_metric, t.id, t.type, tm.name, ctl.capability_id 
FROM cdni_total_limits AS ctl 
LEFT JOIN cdni_telemetry as t ON telemetry_id = t.id 
LEFT JOIN cdni_telemetry_metrics as tm ON telemetry_metric = tm.name`

const hostLimitsQuery = `
SELECT limit_type, maximum_hard, maximum_soft, chl.telemetry_id, chl.telemetry_metric, t.id, t.type, tm.name, host, chl.capability_id 
FROM cdni_host_limits AS chl 
LEFT JOIN cdni_telemetry as t ON telemetry_id = t.id 
LEFT JOIN cdni_telemetry_metrics as tm ON telemetry_metric = tm.name 
ORDER BY host DESC`

func GetCapabilities(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("bearer token header is required"), nil)
		return
	}

	if inf.Config.Cdni == nil || inf.Config.Cdni.JwtDecodingSecret == "" || inf.Config.Cdni.DCdnId == "" {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("cdn.conf does not contain CDNi information"))
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(inf.Config.Cdni.JwtDecodingSecret), nil
	})
	if err != nil {
		api.HandleErr(w, r, nil, http.StatusInternalServerError, fmt.Errorf("parsing claims: %w", err), nil)
		return
	}
	if !token.Valid {
		api.HandleErr(w, r, nil, http.StatusInternalServerError, errors.New("invalid token"), nil)
		return
	}

	var expirationFloat float64
	var ucdn string
	var dcdn string
	for key, val := range claims {
		switch key {
		case "iss":
			if _, ok := val.(string); !ok {
				api.HandleErr(w, r, nil, http.StatusBadRequest, errors.New("invalid token - iss (Issuer) must be a string"), nil)
				return
			}
			ucdn = val.(string)
		case "aud":
			if _, ok := val.(string); !ok {
				api.HandleErr(w, r, nil, http.StatusBadRequest, errors.New("invalid token - aud (Audience) must be a string"), nil)
				return
			}
			dcdn = val.(string)
		case "exp":
			if _, ok := val.(float64); !ok {
				api.HandleErr(w, r, nil, http.StatusBadRequest, errors.New("invalid token - exp (Expiration) must be a float64"), nil)
				return
			}
			expirationFloat = val.(float64)
		}
	}

	expiration := int64(expirationFloat)

	if expiration < time.Now().Unix() {
		api.HandleErr(w, r, nil, http.StatusForbidden, errors.New("token is expired"), nil)
		return
	}
	if dcdn != inf.Config.Cdni.DCdnId {
		api.HandleErr(w, r, nil, http.StatusForbidden, errors.New("invalid token - incorrect dcdn"), nil)
		return
	}
	if ucdn == "" {
		api.HandleErr(w, r, nil, http.StatusForbidden, errors.New("invalid token - empty ucdn field"), nil)
		return
	}

	capacities, err := getCapacities(inf, ucdn)
	if err != nil {
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		return
	}

	telemetries, err := getTelemetries(inf, ucdn)
	if err != nil {
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		return
	}

	fciCaps := Capabilities{}
	capsList := make([]Capability, 0, len(capacities.Capabilities)+len(telemetries.Capabilities))
	capsList = append(capsList, capacities.Capabilities...)
	capsList = append(capsList, telemetries.Capabilities...)

	fciCaps.Capabilities = capsList

	api.WriteRespRaw(w, r, fciCaps)
}

func getFootprintMap(tx *sql.Tx) (map[int][]Footprint, error) {
	footRows, err := tx.Query(AllFootprintQuery)
	if err != nil {
		return nil, fmt.Errorf("querying footprints: %w", err)
	}
	defer log.Close(footRows, "closing foorpint query")
	footprintMap := map[int][]Footprint{}
	for footRows.Next() {
		var footprint Footprint
		if err := footRows.Scan(&footprint.FootprintType, pq.Array(&footprint.FootprintValue), &footprint.CapabilityId); err != nil {
			return nil, fmt.Errorf("scanning db rows: %w", err)
		}

		footprintMap[footprint.CapabilityId] = append(footprintMap[footprint.CapabilityId], footprint)
	}

	return footprintMap, nil
}

func getTotalLimitsMap(tx *sql.Tx) (map[int][]TotalLimitsQueryResponse, error) {
	tlRows, err := tx.Query(totalLimitsQuery)
	if err != nil {
		return nil, fmt.Errorf("querying total limits: %w", err)
	}

	defer log.Close(tlRows, "closing total capacity limits query")
	totalLimitsMap := map[int][]TotalLimitsQueryResponse{}
	for tlRows.Next() {
		var totalLimit TotalLimitsQueryResponse
		if err := tlRows.Scan(&totalLimit.LimitType, &totalLimit.MaximumHard, &totalLimit.MaximumSoft, &totalLimit.TelemetryId, &totalLimit.TelemetryMetic, &totalLimit.Id, &totalLimit.Type, &totalLimit.Name, &totalLimit.CapabilityId); err != nil {
			return nil, fmt.Errorf("scanning db rows: %w", err)
		}

		totalLimitsMap[totalLimit.CapabilityId] = append(totalLimitsMap[totalLimit.CapabilityId], totalLimit)
	}

	return totalLimitsMap, nil
}

func getHostLimitsMap(tx *sql.Tx) (map[int][]HostLimitsResponse, error) {
	hlRows, err := tx.Query(hostLimitsQuery)
	if err != nil {
		return nil, fmt.Errorf("querying host limits: %w", err)
	}

	defer log.Close(hlRows, "closing host capacity limits query")
	hostLimitsMap := map[int][]HostLimitsResponse{}
	for hlRows.Next() {
		var hostLimit HostLimitsResponse
		if err := hlRows.Scan(&hostLimit.LimitType, &hostLimit.MaximumHard, &hostLimit.MaximumSoft, &hostLimit.TelemetryId, &hostLimit.TelemetryMetic, &hostLimit.Id, &hostLimit.Type, &hostLimit.Name, &hostLimit.Host, &hostLimit.CapabilityId); err != nil {
			return nil, fmt.Errorf("scanning db rows: %w", err)
		}

		hostLimitsMap[hostLimit.CapabilityId] = append(hostLimitsMap[hostLimit.CapabilityId], hostLimit)
	}

	return hostLimitsMap, nil
}

func getTelemetriesMap(tx *sql.Tx) (map[int][]Telemetry, error) {
	rows, err := tx.Query(`SELECT id, type, capability_id FROM cdni_telemetry`)
	if err != nil {
		return nil, errors.New("querying cdni telemetry: " + err.Error())
	}
	defer log.Close(rows, "closing telemetry query")

	telemetryMap := map[int][]Telemetry{}
	for rows.Next() {
		telemetry := Telemetry{}
		if err := rows.Scan(&telemetry.Id, &telemetry.Type, &telemetry.CapabilityId); err != nil {
			return nil, errors.New("scanning telemetry: " + err.Error())
		}

		telemetryMap[telemetry.CapabilityId] = append(telemetryMap[telemetry.CapabilityId], telemetry)
	}

	return telemetryMap, nil
}

func getTelemetryMetricsMap(tx *sql.Tx) (map[string][]Metric, error) {
	tmRows, err := tx.Query(`SELECT name, time_granularity, data_percentile, latency, telemetry_id FROM cdni_telemetry_metrics`)
	if err != nil {
		return nil, errors.New("querying cdni telemetry metrics: " + err.Error())
	}
	defer log.Close(tmRows, "closing telemetry metrics query")

	telemetryMetricMap := map[string][]Metric{}
	for tmRows.Next() {
		metric := Metric{}
		if err := tmRows.Scan(&metric.Name, &metric.TimeGranularity, &metric.DataPercentile, &metric.Latency, &metric.TelemetryId); err != nil {
			return nil, errors.New("scanning telemetry metric: " + err.Error())
		}

		telemetryMetricMap[metric.TelemetryId] = append(telemetryMetricMap[metric.TelemetryId], metric)
	}

	return telemetryMetricMap, nil
}

type Capabilities struct {
	Capabilities []Capability `json:"capabilities"`
}

type Capability struct {
	CapabilityType  SupportedCapabilities `json:"capability-type"`
	CapabilityValue interface{}           `json:"capability-value"`
	Footprints      []Footprint           `json:"footprints"`
}

type CapacityCapabilityValue struct {
	TotalLimits []Limit     `json:"total-limits"`
	HostLimits  []HostLimit `json:"host-limits"`
}

type HostLimit struct {
	Host   string  `json:"host"`
	Limits []Limit `json:"limits"`
}

type Limit struct {
	LimitType       CapacityLimitType `json:"limit-type"`
	MaximumHard     int64             `json:"maximum-hard"`
	MaximumSoft     int64             `json:"maximum-soft"`
	TelemetrySource TelemetrySource   `json:"telemetry-source"`
}

type TelemetrySource struct {
	Id     string `json:"id"`
	Metric string `json:"metric"`
}

type TelemetryCapabilityValue struct {
	Sources []Telemetry `json:"sources"`
}

type Telemetry struct {
	Id           string              `json:"id"`
	Type         TelemetrySourceType `json:"type"`
	CapabilityId int                 `json:"-"`
	Metrics      []Metric            `json:"metrics"`
}

type Metric struct {
	Name            string `json:"name"`
	TimeGranularity int    `json:"time-granularity"`
	DataPercentile  int    `json:"data-percentile"`
	Latency         int    `json:"latency"`
	TelemetryId     string `json:"-"`
}

type Footprint struct {
	FootprintType  FootprintType `json:"footprint-type" db:"footprint_type"`
	FootprintValue []string      `json:"footprint-value" db:"footprint_value"`
	CapabilityId   int           `json:"-"`
}

type CapacityLimitType string

const (
	Egress         CapacityLimitType = "egress"
	Requests                         = "requests"
	StorageSize                      = "storage-size"
	StorageObjects                   = "storage-objects"
	Sessions                         = "sessions"
	CacheSize                        = "cache-size"
)

type SupportedCapabilities string

const (
	FciTelemetry      SupportedCapabilities = "FCI.Telemetry"
	FciCapacityLimits                       = "FCI.CapacityLimits"
)

type TelemetrySourceType string

const (
	Generic TelemetrySourceType = "generic"
)

type FootprintType string

const (
	Ipv4Cidr    FootprintType = "ipv4cidr"
	Ipv6Cidr                  = "ipv6cidr"
	Asn                       = "asn"
	CountryCode               = "countrycode"
)
