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
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/v8/traffic_monitor/tools/testcaches/fakesrvrdata"
)

func News(portStart int, numPorts int, remaps []string) ([]*http.Server, error) {
	servers := []*http.Server{}
	for i := 0; i < numPorts; i++ {
		port := portStart + i
		server, err := New(port, remaps)
		if err != nil {
			// TODO stop all servers already created?
			return nil, errors.New("making server on port " + strconv.Itoa(port) + ": " + err.Error())
		}
		servers = append(servers, server)
	}
	return servers, nil
}

func New(port int, remaps []string) (*http.Server, error) {
	serverData, remapIncrements := newData(remaps)
	fakeServerThs, err := fakesrvrdata.Run(serverData, remapIncrements)
	if err != nil {
		return nil, errors.New("running FakeServer: " + err.Error())
	}
	fmt.Println("Starting Serving on port " + strconv.Itoa(port)) // debug
	srvr := Serve(port, fakeServerThs)
	fmt.Println("Serving on port " + strconv.Itoa(port)) // debug
	return srvr, nil
}

func newData(remaps []string) (fakesrvrdata.FakeServerData, map[string]fakesrvrdata.BytesPerSec) {
	serverDataRemap := map[string]fakesrvrdata.FakeRemap{}
	for _, remap := range remaps {
		serverDataRemap[remap] = fakesrvrdata.FakeRemap{}
	}
	serverData := fakesrvrdata.FakeServerData{
		ATS: fakesrvrdata.FakeATS{
			Server: "7.1.4",
			Remaps: serverDataRemap,
		},
		System: fakesrvrdata.FakeSystem{
			Name:       "bond0",
			Speed:      20000,
			ProcNetDev: fakesrvrdata.FakeProcNetDev{Interface: "bond0"},
			ProcLoadAvg: fakesrvrdata.FakeProcLoadAvg{
				CPU1m:        4.52,
				CPU5m:        4.64,
				CPU10m:       4.47,
				RunningProcs: 2,
				TotalProcs:   1592,
				LastPIDUsed:  32174,
			},
			ConfgiReloadRequests: 10,
			LastReloadRequest:    1512768709,
			ConfigReloads:        1,
			LastReload:           1512678516,
			AstatsLoad:           1512678516,
			Something:            "here",
		},
	}

	remapIncrements := map[string]fakesrvrdata.BytesPerSec{}
	for i, remap := range remaps {
		scale := uint64(float64(i) / float64(len(remaps)) * 100.0)
		remapIncrements[remap] = fakesrvrdata.BytesPerSec{
			Min: fakesrvrdata.FakeRemap{
				InBytes:   1,
				OutBytes:  1,
				Status2xx: 1,
				Status3xx: 0,
				Status4xx: 0,
				Status5xx: 0,
			},
			Max: fakesrvrdata.FakeRemap{
				InBytes:   scale/2 + 1,
				OutBytes:  scale + 1,
				Status2xx: scale/10 + 1,
				Status3xx: scale / 40,
				Status4xx: scale / 20,
				Status5xx: scale / 60,
			},
		}
	}
	return serverData, remapIncrements
}
