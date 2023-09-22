package fakesrvr

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
	"html"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"unsafe"

	"github.com/apache/trafficcontrol/v8/traffic_monitor/tools/testcaches/fakesrvrdata"
)

type CmdFunc = func(http.ResponseWriter, *http.Request, fakesrvrdata.Ths)

var cmds = map[string]CmdFunc{
	"setstat":   cmdSetStat,
	"setdelay":  cmdSetDelay,
	"setsystem": cmdSetSystem,
}

// cmdSetStat sets the rate of the given stat increase for the given remap.
//
// query parameters:
//
//	remap: string; required; the full name of the remap whose kbps to set.
//	stat: string; required; the stat to set (in_bytes, out_bytes, status_2xx, status_3xx, status_4xx, status_5xx).
//	min:   unsigned integer; required; new minimum of kbps increase of InBytes stat for the given remap.
//	max:   unsigned integer; required; new maximum of kbps increase of InBytes stat for the given remap.
func cmdSetStat(w http.ResponseWriter, r *http.Request, fakeSrvrDataThs fakesrvrdata.Ths) {
	urlQry := r.URL.Query()

	newMinStr := html.EscapeString(urlQry.Get("min"))
	newMin, err := strconv.ParseUint(newMinStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error parsing query parameter 'min': must be a positive integer: " + err.Error() + "\n"))
		return
	}

	newMaxStr := html.EscapeString(urlQry.Get("max"))
	newMax, err := strconv.ParseUint(newMaxStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error parsing query parameter 'max': must be a positive integer: " + err.Error() + "\n"))
		return
	}

	remap := html.EscapeString(urlQry.Get("remap"))
	if remap == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing query parameter 'remap': must specify a remap to set\n"))
		return
	}

	stat := html.EscapeString(urlQry.Get("stat"))

	validStats := map[string]struct{}{
		"in_bytes":   {},
		"out_bytes":  {},
		"status_2xx": {},
		"status_3xx": {},
		"status_4xx": {},
		"status_5xx": {},
	}

	if _, ok := validStats[stat]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		statNames := []string{}
		for statName := range validStats {
			statNames = append(statNames, statName)
		}
		w.Write([]byte("error with query parameter 'stat' '" + stat + "': not found. Valid stats are: [" + strings.Join(statNames, ",") + "\n"))
		return
	}

	srvr := (*fakesrvrdata.FakeServerData)(fakeSrvrDataThs.Get())
	if _, ok := srvr.ATS.Remaps[remap]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		remapNames := []string{}
		for remapName := range srvr.ATS.Remaps {
			remapNames = append(remapNames, remapName)
		}
		w.Write([]byte("error with query parameter 'remap' '" + remap + "': not found. Valid remaps are: [" + strings.Join(remapNames, ",") + "\n"))
		return
	}

	incs := <-fakeSrvrDataThs.GetIncrementsChan
	inc := incs[remap]

	switch stat {
	case "in_bytes":
		inc.Min.InBytes = newMin
		inc.Max.InBytes = newMax
	case "out_bytes":
		inc.Min.OutBytes = newMin
		inc.Max.OutBytes = newMax
	case "status_2xx":
		inc.Min.Status2xx = newMin
		inc.Max.Status2xx = newMax
	case "status_3xx":
		inc.Min.Status3xx = newMin
		inc.Max.Status3xx = newMax
	case "status_4xx":
		inc.Min.Status4xx = newMin
		inc.Max.Status4xx = newMax
	case "status_5xx":
		inc.Min.Status5xx = newMin
		inc.Max.Status5xx = newMax
	default:
		panic("unknown stat; should never happen")
	}

	fakeSrvrDataThs.IncrementChan <- fakesrvrdata.IncrementChanT{RemapName: remap, BytesPerSec: inc}

	w.WriteHeader(http.StatusNoContent)
}

func cmdSetDelay(w http.ResponseWriter, r *http.Request, fakeSrvrDataThs fakesrvrdata.Ths) {
	urlQry := r.URL.Query()

	newMinStr := urlQry.Get("min")
	newMin, err := strconv.ParseUint(newMinStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error parsing query parameter 'min': must be a non-negative integer: " + err.Error() + "\n"))
		return
	}

	newMaxStr := urlQry.Get("max")
	newMax, err := strconv.ParseUint(newMaxStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error parsing query parameter 'max': must be a non-negative integer: " + err.Error() + "\n"))
		return
	}

	newMinMax := fakesrvrdata.MinMaxUint64{Min: newMin, Max: newMax}
	newMinMaxPtr := &newMinMax

	p := (unsafe.Pointer)(newMinMaxPtr)
	atomic.StorePointer(fakeSrvrDataThs.DelayMS, p)
	w.WriteHeader(http.StatusNoContent)
}

func cmdSetSystem(w http.ResponseWriter, r *http.Request, fakeSrvrDataThs fakesrvrdata.Ths) {
	urlQry := r.URL.Query()

	if newSpeedStr := urlQry.Get("speed"); newSpeedStr != "" {
		newSpeed, err := strconv.ParseInt(newSpeedStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error parsing query parameter 'speed': must be a non-negative integer: " + err.Error() + "\n"))
			return
		}

		srvr := (*fakesrvrdata.FakeServerData)(fakeSrvrDataThs.Get())
		srvr.System.Speed = int(newSpeed)
		fakeSrvrDataThs.Set(srvr)
	}

	if newLoadAvg1MStr := urlQry.Get("loadavg1m"); newLoadAvg1MStr != "" {
		newLoadAvg1M, err := strconv.ParseFloat(newLoadAvg1MStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error parsing query parameter 'loadavg1m': must be a number: " + err.Error() + "\n"))
			return
		}

		srvr := (*fakesrvrdata.FakeServerData)(fakeSrvrDataThs.Get())
		srvr.System.ProcLoadAvg.CPU1m = newLoadAvg1M
		fakeSrvrDataThs.Set(srvr)
	}

	if newLoadAvg5MStr := urlQry.Get("loadavg5m"); newLoadAvg5MStr != "" {
		newLoadAvg5M, err := strconv.ParseFloat(newLoadAvg5MStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error parsing query parameter 'loadavg5m': must be a number: " + err.Error() + "\n"))
			return
		}

		srvr := (*fakesrvrdata.FakeServerData)(fakeSrvrDataThs.Get())
		srvr.System.ProcLoadAvg.CPU5m = newLoadAvg5M
		fakeSrvrDataThs.Set(srvr)
	}

	if newLoadAvg10MStr := urlQry.Get("loadavg10m"); newLoadAvg10MStr != "" {
		newLoadAvg10M, err := strconv.ParseFloat(newLoadAvg10MStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error parsing query parameter 'loadavg10m': must be a non-negative integer: " + err.Error() + "\n"))
			return
		}

		srvr := (*fakesrvrdata.FakeServerData)(fakeSrvrDataThs.Get())
		srvr.System.ProcLoadAvg.CPU10m = newLoadAvg10M
		fakeSrvrDataThs.Set(srvr)
	}

	w.WriteHeader(http.StatusNoContent)
}
