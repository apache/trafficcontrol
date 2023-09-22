package load

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
	"strings"
	"sync"

	"fmt"
	"github.com/apache/trafficcontrol/v8/test/router/client"
	"github.com/apache/trafficcontrol/v8/test/router/data"
)

type LoadTest struct {
	CaFile                string   `json:"caFile"`
	Cdn                   string   `json:"cdn"`
	TxCount               int      `json:"txCount"`
	Connections           int      `json:"connections"`
	HttpDeliveryServices  []string `json:"httpDeliveryServices"`
	HttpsDeliveryServices []string `json:"httpsDeliveryServices"`
}

func DoLoadTest(loadtest LoadTest, done chan struct{}) chan data.HttpResult {
	resultsChan := make(chan data.HttpResult)

	go func() {
		fmt.Println("Starting load test", loadtest)
		defer close(done)
		var waitGroup sync.WaitGroup
		for _, deliveryService := range loadtest.HttpDeliveryServices {
			waitGroup.Add(1)
			go func(ds string) {
				defer waitGroup.Done()
				host := strings.Join([]string{ds, loadtest.Cdn}, ".")
				tlsConfig := client.MustGetTlsConfiguration(host, loadtest.CaFile)
				client.ExerciseDeliveryService(false, tlsConfig, host, loadtest.TxCount, loadtest.Connections, resultsChan)
			}(deliveryService)
		}

		for _, deliveryService := range loadtest.HttpsDeliveryServices {
			waitGroup.Add(1)
			go func(ds string) {
				defer waitGroup.Done()
				host := strings.Join([]string{ds, loadtest.Cdn}, ".")
				tlsConfig := client.MustGetTlsConfiguration(host, loadtest.CaFile)
				client.ExerciseDeliveryService(true, tlsConfig, host, loadtest.TxCount, loadtest.Connections, resultsChan)
			}(deliveryService)
		}

		waitGroup.Wait()
	}()

	return resultsChan
}
