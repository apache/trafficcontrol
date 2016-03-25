/*
   Copyright 2015 Comcast Cable Communications Management, LLC

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/prometheus/log"
	"golang.org/x/net/publicsuffix"
)

// Session ...
type Session struct {
	UserName  string
	Password  string
	URL       string
	UserAgent *http.Client
	Cache     map[string]cacheEntry
}

// Result {"alerts":[{"level":"success","text":"Successfully logged in."}],"version":"1.1"}
type Result struct {
	Alerts  []Alert
	Version string `json:"version"`
}

// Alert ...
type Alert struct {
	Level string `json:"level"`
	Text  string `json:"text"`
}

type cacheEntry struct {
	Entered int64
	bytes   []byte
}

// Credentials contains Traffic Ops login credentials
type Credentials struct {
	Username string `json:"u"`
	Password string `json:"p"`
}

// TODO JvD
const tmPollingInterval = 60

// loginCreds gathers login credentials for Traffic Ops.
func loginCreds(toUser string, toPasswd string) ([]byte, error) {
	credentials := Credentials{
		Username: toUser,
		Password: toPasswd,
	}

	js, err := json.Marshal(credentials)
	if err != nil {
		err := fmt.Errorf("Error creating login json: %v", err)
		return nil, err
	}
	return js, nil
}

// Login to traffic_ops, the response should set the cookie for this session
// automatically. Start with
//     to := traffic_ops.Login("user", "passwd", true)
// subsequent calls like to.GetData("datadeliveryservice") will be authenticated.
func Login(toURL string, toUser string, toPasswd string, insecure bool) (*Session, error) {
	credentials, err := loginCreds(toUser, toPasswd)
	if err != nil {
		return nil, err
	}

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, err
	}

	to := Session{
		UserAgent: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
			},
			Jar: jar,
		},
		URL: toURL,
	}

	uri := "/api/1.1/user/login"
	resp, err := to.request(uri, credentials)
	if err != nil {
		return nil, err
	}

	var result Result
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	success := false
	for _, alert := range result.Alerts {
		if alert.Level == "success" && alert.Text == "Successfully logged in." {
			success = true
			break
		}
	}

	if !success {
		fmt.Println("NO SUCCESS")
		err := fmt.Errorf("Login failed, result string: %+v", result)
		return nil, err
	}

	log.Infof("logged into %s!", toURL)
	return &to, nil
}

// request performs the actual HTTP request to Traffic Ops
func (to *Session) request(path string, body []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", to.URL, path)

	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
	}

	resp, err := to.UserAgent.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}

// getBytesWithTTL - get the path, and cache in the session
// return from cache is found and the ttl isn't expired, otherwise get it and
// store it in cache
func (to *Session) getBytesWithTTL(path string, ttl int64) ([]byte, error) {

	var body []byte
	var err error
	getFresh := false
	if cacheEntry, ok := to.Cache[path]; ok {
		if cacheEntry.Entered > time.Now().Unix()-ttl {
			fmt.Println("Cache HIT for", path)
			body = cacheEntry.bytes
		} else {
			fmt.Println("Cache HIT but EXPIRED for", path)
			getFresh = true
		}
	} else {
		fmt.Println("Cache MISS for", path)
		getFresh = true
	}

	if getFresh {
		body, err = to.getBytes(path)
		if err != nil {
			return nil, err
		}

		var newEntry cacheEntry
		newEntry.Entered = time.Now().Unix()
		newEntry.bytes = body
		to.Cache[path] = newEntry
	}
	return body, nil
}

// GetBytes - get []bytes array for a certain path on the to session.
// returns the raw body
func (to *Session) getBytes(path string) ([]byte, error) {
	resp, err := to.request(path, nil)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
