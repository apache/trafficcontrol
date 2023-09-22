package client

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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/apache/trafficcontrol/v8/test/router/data"
)

func MustLoadCertificates(cafile string) *x509.CertPool {
	pem, err := ioutil.ReadFile(cafile)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pem) {
		panic("Failed appending certs")
	}

	return certPool
}

func MustGetTlsConfiguration(servername string, cafile string) *tls.Config {
	config := &tls.Config{}
	certPool := MustLoadCertificates(cafile)

	config.RootCAs = certPool
	config.ClientCAs = certPool

	config.ClientAuth = tls.RequireAndVerifyClientCert

	config.CipherSuites = []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256}

	//Use only TLS v1.2
	config.MinVersion = tls.VersionTLS12

	//Don't allow session resumption
	config.SessionTicketsDisabled = true
	return config
}

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwzyz0123456789/"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func DoRequest(isHttps bool, tlsConfig *tls.Config, host string, resultsChan chan<- data.HttpResult) {
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: transport,
	}

	scheme := "http://"
	if isHttps {
		scheme = "https://"
	}

	url := fmt.Sprintf("%vccr.%v/%v/stuff?fakeClientIpAddress=68.87.25.123", scheme, host, RandomString(23))
	fmt.Println("******", url)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("^^^^^^", err)
	}

	req.Header.Set("Connection", "close")

	before := time.Now()
	resp, err := client.Do(req)
	latency := time.Now().Sub(before).Nanoseconds() / int64(1000)

	httpResult := data.HttpResult{
		RequestTime: before,
		Host:        req.Host,
		LatencyUsec: latency,
		Status:      -1,
	}

	if err != nil {
		fmt.Println("*** Error ****", err)
		httpResult.Err = err
		fmt.Println(httpResult)
		return
	}

	if resp != nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 302 {
			fmt.Println("Received unexpected http status", resp.StatusCode)
		}
		httpResult.Status = resp.StatusCode
	}

	resultsChan <- httpResult
}

func ExerciseDeliveryService(isHttps bool, tlsConfig *tls.Config, host string, numRequests int, maxWorkers int, resultsChan chan<- data.HttpResult) {
	var waitGroup sync.WaitGroup

	workers := make(chan int, maxWorkers)
	for worker := 0; worker < maxWorkers; worker++ {
		waitGroup.Add(1)

		go func(ch <-chan int) {
			defer waitGroup.Done()

			for _ = range ch {
				DoRequest(isHttps, tlsConfig, host, resultsChan)
			}
		}(workers)
	}

	for i := 0; i < numRequests; i++ {
		time.Sleep(100 * time.Microsecond)
		workers <- i
	}
	close(workers)

	waitGroup.Wait()
}
