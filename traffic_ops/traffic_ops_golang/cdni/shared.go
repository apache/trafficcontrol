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
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

const CapabilityQuery = `SELECT id, type, ucdn FROM cdni_capabilities WHERE type = $1 AND ucdn = $2`
const FootprintQuery = `SELECT footprint_type, footprint_value::text[] FROM cdni_footprints WHERE capability_id = $1`

func GetCapabilities(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	bearerToken := r.Header.Get("Authorization")

	if inf.Config.Cdni == nil || inf.Config.Cdni.JwtDecodingSecret == "" || inf.Config.Cdni.DCdnId == "" {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("cdn.conf does not contain CDNi information"))
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(inf.Config.Cdni.JwtDecodingSecret), nil
	})
	if err != nil {
		api.HandleErr(w, r, nil, http.StatusInternalServerError, errors.New("parsing claims"), nil)
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
			ucdn = val.(string)
		case "aud":
			dcdn = val.(string)
		case "exp":
			expirationFloat = val.(float64)
		}
	}

	expiration := int64(expirationFloat)

	if expiration < time.Now().Unix() {
		api.HandleErr(w, r, nil, http.StatusForbidden, errors.New("token is expired"), nil)
		return
	}
	if dcdn != inf.Config.Cdni.DCdnId {
		api.HandleErr(w, r, nil, http.StatusForbidden, errors.New("invalid token"), nil)
		return
	}
	if ucdn == "" {
		api.HandleErr(w, r, nil, http.StatusForbidden, errors.New("invalid token"), nil)
		return
	}

	capacities, err := GetCapacities(inf, ucdn)
	if err != nil {
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		return
	}

	telemetries, err := GetTelemetries(inf, ucdn)
	if err != nil {
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		return
	}

	fciCaps := Capabilities{}
	capsList := []Capability{}
	capsList = append(capsList, capacities.Capabilities...)
	capsList = append(capsList, telemetries.Capabilities...)

	fciCaps.Capabilities = capsList

	api.WriteRespRaw(w, r, fciCaps)
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
	Id      string              `json:"id"`
	Type    TelemetrySourceType `json:"type"`
	Metrics []Metric            `json:"metrics"`
}

type Metric struct {
	Name            string `json:"name"`
	TimeGranularity int    `json:"time-granularity"`
	DataPercentile  int    `json:"data-percentile"`
	Latency         int    `json:"latency"`
}

type Footprint struct {
	FootprintType  FootprintType `json:"footprint-type" db:"footprint_type"`
	FootprintValue []string      `json:"footprint-value" db:"footprint_value"`
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
