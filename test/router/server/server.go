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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/apache/trafficcontrol/v8/test/router/data"
	"github.com/apache/trafficcontrol/v8/test/router/load"
)

var done chan struct{}
var resultChan chan data.HttpResult
var results []data.HttpResult

type credentials struct {
	User     string `json:"u"`
	Password string `json:"p"`
}

var opsCookies []*http.Cookie

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "GET" && r.URL.Path == "/report" {
		data, err := ioutil.ReadFile("foo.json")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		fmt.Fprintf(w, string(data))
	}

	if r.Method == "POST" && r.URL.Path == "/api/4.0/user/login" {

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}

		url := fmt.Sprintf("https://%v/api/4.0/user/login", r.URL.Query().Get("opsHost"))

		resp, err := client.Post(url, "application/json", r.Body)

		if err != nil {
			fmt.Println("Failed to proxy authentication to traffic ops", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Println("Dangit!!!! got non 200", resp)
			w.WriteHeader(500)
			return
		}

		opsCookies = resp.Cookies()
		fmt.Println("ops cookies", opsCookies)

		for _, cookie := range resp.Cookies() {
			fmt.Println("cookie", cookie)
			http.SetCookie(w, cookie)
		}

		fmt.Println("woo-hoo I think I proxied authentication!!!")

		w.Write([]byte(fmt.Sprintf("{\"opsHost\":\"%v\"}", r.URL.Query().Get("opsHost"))))
		return
	}

	opsHost := r.URL.Query().Get("opsHost")
	if len(opsHost) > 0 {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}

		urlString := fmt.Sprintf("https://%v%v", r.URL.Query().Get("opsHost"), r.URL.Path)
		u, err := url.Parse(urlString)

		if err != nil {
			fmt.Println("Failed parsing", urlString, err.Error())
			w.WriteHeader(500)
			return
		}

		if u == nil {
			fmt.Println("Crap")
			w.WriteHeader(500)
			return
		}

		fmt.Println("url", u)

		fmt.Println("AAAAA")
		fmt.Println(opsCookies)
		fmt.Println("BBBB")

		fmt.Println("client", client)
		fmt.Println("jar", client.Jar)

		client.Jar, err = cookiejar.New(nil)

		if err != nil {
			fmt.Println("Failed setting up cookie jar")
			w.WriteHeader(500)
			return
		}

		client.Jar.SetCookies(u, opsCookies)

		fmt.Println(client.Jar)

		resp, err := client.Get(urlString)

		if err != nil {
			fmt.Println("Failed to proxy ", r.URL, "to host", r.URL.Query().Get("opsHost"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		buf, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("Failed reading body of response!")
			w.WriteHeader(500)
		}

		fmt.Println("proxying response", string(buf))
		w.Write(buf)
		return
	}

	if r.Method == "POST" && r.URL.Path == "/loadtest" {
		results = nil
		var lt load.LoadTest
		err := json.NewDecoder(r.Body).Decode(&lt)

		if err != nil {
			fmt.Println("Failed to unmarshal Json!", err.Error())
		}

		done = make(chan struct{})
		resultChan = load.DoLoadTest(lt, done)

		go func() {
			for {
				select {
				case result := <-resultChan:
					results = append(results, result)
				case <-done:
					return
				}
			}
		}()

		w.Write([]byte(`{"status":"started"}`))
	}

	if r.Method == "GET" && r.URL.Path == "/loadtest" {
		b, _ := json.MarshalIndent(results, "", "  ")

		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	}
}

func main() {
	done = make(chan struct{})

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8888", nil)
}
