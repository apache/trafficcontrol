package fakesrvrdata

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
	"bytes"
	"strconv"
)

type ThsT *FakeServerData

type FakeServerData struct {
	ATS    FakeATS    `json:"ats"`
	System FakeSystem `json:"system"`
}

// GetSystem returns a FakeServerSystem, which can be serialized as the JSON response of Astats with an "application=system" query
func (s *FakeServerData) GetSystem() FakeServerSystem {
	return FakeServerSystem{
		ATS: FakeSystemATS{
			Server: s.ATS.Server,
		},
		System: s.System,
	}
}

type FakeATS struct {
	Server string `json:"server"`
	Remaps map[string]FakeRemap
}

func copyRemaps(old map[string]FakeRemap) map[string]FakeRemap {
	new := make(map[string]FakeRemap, len(old))
	for k, v := range old {
		new[k] = v
	}
	return new
}

type FakeRemap struct {
	InBytes   uint64
	OutBytes  uint64
	Status2xx uint64
	Status3xx uint64
	Status4xx uint64
	Status5xx uint64
}

func (a FakeATS) MarshalJSON() ([]byte, error) {
	// `copy` is faster than bytes.Buffer, if we precompute the length, which is possible but not trivial.
	var b bytes.Buffer
	b.WriteString("{")
	for remapName, remap := range a.Remaps {
		b.WriteString(`"plugin.remap_stats.`)
		b.WriteString(remapName)
		b.WriteString(`.in_bytes": `)
		b.WriteString(strconv.FormatUint(remap.InBytes, 10))
		b.WriteString(",")
		b.WriteString(`"plugin.remap_stats.`)
		b.WriteString(remapName)
		b.WriteString(`.out_bytes": `)
		b.WriteString(strconv.FormatUint(remap.OutBytes, 10))
		b.WriteString(",")
		b.WriteString(`"plugin.remap_stats.`)
		b.WriteString(remapName)
		b.WriteString(`.status_2xx": `)
		b.WriteString(strconv.FormatUint(remap.Status2xx, 10))
		b.WriteString(",")
		b.WriteString(`"plugin.remap_stats.`)
		b.WriteString(remapName)
		b.WriteString(`.status_3xx": `)
		b.WriteString(strconv.FormatUint(remap.Status3xx, 10))
		b.WriteString(",")
		b.WriteString(`"plugin.remap_stats.`)
		b.WriteString(remapName)
		b.WriteString(`.status_4xx": `)
		b.WriteString(strconv.FormatUint(remap.Status4xx, 10))
		b.WriteString(",")
		b.WriteString(`"plugin.remap_stats.`)
		b.WriteString(remapName)
		b.WriteString(`.status_5xx": `)
		b.WriteString(strconv.FormatUint(remap.Status5xx, 10))
		b.WriteString(",")
	}
	b.WriteString(`"server": "`)
	b.WriteString(a.Server)
	b.WriteString(`"}`)
	return b.Bytes(), nil
}

type FakeServerSystem struct {
	ATS    FakeSystemATS `json:"ats"`
	System FakeSystem    `json:"system"`
}

type FakeSystemATS struct {
	Server string `json:"server"`
}

type FakeSystem struct {
	Name                 string          `json:"inf.name"`
	Speed                int             `json:"inf.speed"`
	ProcNetDev           FakeProcNetDev  `json:"proc.net.dev"`
	ProcLoadAvg          FakeProcLoadAvg `json:"proc.loadavg"`
	ConfgiReloadRequests uint64          `json:"configReloadRequests"`
	LastReloadRequest    uint64          `json:"lastReloadRequest"`
	ConfigReloads        uint64          `json:"configReloads"`
	LastReload           uint64          `json:"lastReload"`
	AstatsLoad           uint64          `json:"astatsLoad"`
	Something            string          `json:"something"`
	Version              string          `json:"application_version"`
}

type FakeProcLoadAvg struct {
	CPU1m        float64
	CPU5m        float64
	CPU10m       float64
	RunningProcs int
	TotalProcs   int
	LastPIDUsed  int
}

func (p FakeProcLoadAvg) MarshalJSON() ([]byte, error) {
	return []byte("\"" +
		strconv.FormatFloat(p.CPU1m, 'f', -1, 64) + " " +
		strconv.FormatFloat(p.CPU5m, 'f', -1, 64) + " " +
		strconv.FormatFloat(p.CPU10m, 'f', -1, 64) + " " +
		strconv.Itoa(p.RunningProcs) + "/" +
		strconv.Itoa(p.TotalProcs) + " " +
		strconv.Itoa(p.LastPIDUsed) + "\""), nil
}

type FakeProcNetDev struct {
	Interface      string
	RcvBytes       uint64
	RcvPackets     uint64
	RcvErrs        uint64
	RcvDropped     uint64
	RcvFIFOErrs    uint64
	RcvFrameErrs   uint64
	RcvCompressed  uint64
	RcvMulticast   uint64
	SndBytes       uint64
	SndPackets     uint64
	SndErrs        uint64
	SndDropped     uint64
	SndFIFOErrs    uint64
	SndCollisions  uint64
	SndCarrierErrs uint64
	SndCompressed  uint64
}

func (p FakeProcNetDev) MarshalJSON() ([]byte, error) {
	return []byte("\"" + p.Interface + ": " +
		strconv.FormatUint(p.RcvBytes, 10) + " " +
		strconv.FormatUint(p.RcvPackets, 10) + " " +
		strconv.FormatUint(p.RcvErrs, 10) + " " +
		strconv.FormatUint(p.RcvDropped, 10) + " " +
		strconv.FormatUint(p.RcvFIFOErrs, 10) + " " +
		strconv.FormatUint(p.RcvFrameErrs, 10) + " " +
		strconv.FormatUint(p.RcvCompressed, 10) + " " +
		strconv.FormatUint(p.RcvMulticast, 10) + " " +
		strconv.FormatUint(p.SndBytes, 10) + " " +
		strconv.FormatUint(p.SndPackets, 10) + " " +
		strconv.FormatUint(p.SndErrs, 10) + " " +
		strconv.FormatUint(p.SndDropped, 10) + " " +
		strconv.FormatUint(p.SndFIFOErrs, 10) + " " +
		strconv.FormatUint(p.SndCollisions, 10) + " " +
		strconv.FormatUint(p.SndCarrierErrs, 10) + " " +
		strconv.FormatUint(p.SndCompressed, 10) + "\""), nil
}
