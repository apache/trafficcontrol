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

import "fmt"

func ExampleLoadavgFromRawLine() {
	loadavg, err := LoadavgFromRawLine("0.30 0.12 0.21 1/863 1421")
	fmt.Println(err)
	fmt.Printf("%.2f %.2f %.2f %d/%d %d", loadavg.One, loadavg.Five, loadavg.Fifteen, loadavg.CurrentProcesses, loadavg.TotalProcesses, loadavg.LatestPID)
	// Output: <nil>
	// 0.30 0.12 0.21 1/863 1421
}

func ExampleStatistics_AddInterfaceFromRawLine() {
	var s Statistics
	raw := "eth0:47907832129 14601260    0    0    0     0          0   790726 728207677726 10210700052    0    0    0     0       0          0"

	if err := s.AddInterfaceFromRawLine(raw); err != nil {
		fmt.Println(err)
		return
	}

	iface, ok := s.Interfaces["eth0"]
	if !ok {
		fmt.Printf("Error, no 'eth0' interface!\n%+v", s.Interfaces)
		return
	}
	fmt.Printf("eth0: {BytesOut: %d, BytesIn: %d}", iface.BytesOut, iface.BytesIn)
	// Output: eth0: {BytesOut: 728207677726, BytesIn: 47907832129}
}
