package main

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

	"time"
	"fmt"
	"github.com/apache/trafficcontrol/traffic_ops/client"
	"sync"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"strconv"
	"testing"
)

const (
	opsUrl         = "https://yourTO.net"
	opsUser        = ""
	opsPass        = ""
)

var to *client.Session

func CreateToConnection(trafficOpsURL string, username string, password string) error {
	toLocal, _, err := client.LoginWithAgent(trafficOpsURL, username, password, true, "cdn-bgp-consumer",
		false, time.Second*time.Duration(30))
	if err != nil {
		fmt.Errorf("Unable to login to TO: %v", err)
		return err
	}
	to = toLocal

	return nil
}

//
// go test -run TestCertificateType
//
func TestCertificateType(t *testing.T) {
	CreateToConnection(opsUrl, opsUser, opsPass)

	deliveryServices, _, err :=  to.GetDeliveryServices()

	if err != nil {
		fmt.Errorf("Unable to get delivery services: %v", err)
		return
	}

	messages := make(chan string)
	var wg sync.WaitGroup
	totalSecDS := 0

	for _, deliveryService := range deliveryServices {
		if deliveryService.Protocol > 0 {
			wg.Add(1)
			totalSecDS += 1
		}
	}

	certTypes := make(map[string]int)

	fmt.Printf("Total DS: %d, ds with protocol != http_only: %d\n", len(deliveryServices), totalSecDS)
	fmt.Printf("Getting DS certificates, need try multiple times for some DS\n")

	errorsCount := 0
	errorsString := ""

	for _, deliveryService := range deliveryServices {
		if deliveryService.Protocol == 0 {
			continue
		}
		deliveryServiceCopy := tc.DeliveryService(deliveryService)

		go func() {
			defer wg.Done()
			keepTrying := true

			for keepTrying {

				fmt.Print("Trying " + deliveryServiceCopy.XMLID + "\n")
				// riak has trouble services too many requests at the same time
				deliveryServiceSSLKeys, _, error := to.GetDeliveryServiceSSLKeysByID(deliveryServiceCopy.XMLID)

				if error != nil {
					fmt.Print("Could not get ssl key for " + deliveryServiceCopy.XMLID + ", trying again\n")
					time.Sleep(time.Second)
					continue
				}

				keepTrying = false

				//messages <- spew.Sdump(deliveryServiceSSLKeys)

				if certsChain, err := decodeCertificate(deliveryServiceSSLKeys.Certificate.Crt); err != nil {
					errorStr := fmt.Sprintf("ERROR: could not decodeCertificate for %v, %v\n", deliveryServiceCopy.XMLID, err)
					fmt.Print(errorStr)
					errorsCount += 1
					errorsString = errorsString + errorStr
					return
				} else {
					certsChainStr := ""
					for index, cert := range certsChain {
						certsChainStr = certsChainStr + strconv.Itoa(index) + ". " + cert.SignatureAlgorithm.String() + "\n"
					}
					fmt.Print(deliveryServiceCopy.XMLID +": \n" + certsChainStr )
					if _, ok := certTypes[certsChainStr]; !ok {
						certTypes[certsChainStr] = 0
					}
					certTypes[certsChainStr] += 1
				}
			}
		}()
	}

	go func() {
		for message := range messages {
			fmt.Print(message)
		}
	}()

	wg.Wait()

	fmt.Printf("\nTotal DS: %d, ds with protocol != http_only: %d\n\n", len(deliveryServices), totalSecDS)

	for certType, num := range certTypes {
		fmt.Printf("%d delivery services has sig:\n%s", num, certType)
	}
	fmt.Printf("Had %d errors:\n%s", errorsCount, errorsString)
}


