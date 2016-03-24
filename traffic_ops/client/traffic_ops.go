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

	log "github.com/cihub/seelog"
	"golang.org/x/net/publicsuffix"
)

// Session ...
type Session struct {
	UserName  string
	Password  string
	URL       string
	UserAgent *http.Client
	Cache     map[string]cacheentry
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

type cacheentry struct {
	Entered int64
	bytes   []byte
}

// Credentials ..
type Credentials struct {
	Username string `json:"u"`
	Password string `json:"p"`
}

// TODO JvD
const tmPollingInterval = 60

// GetBytesWithTTL - get the path, and cache in the session
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
		var newEntry cacheentry
		newEntry.Entered = time.Now().Unix()
		newEntry.bytes = body
		to.Cache[path] = newEntry
	}
	return body, err
}

// GetBytes - get []bytes array for a certain path on the to session.
// returns the raw body
func (to *Session) getBytes(path string) ([]byte, error) {
	var body []byte
	resp, err := to.UserAgent.Get(fmt.Sprintf("%s%s", to.URL, path))
	if err != nil {
		log.Info(err)
		return body, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info(err)
	}
	return body, err
}

func (to *Session) postJSON(path string, body []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", to.URL, path)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := to.UserAgent.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	// resp, err := to.UserAgent.Post(fmt.Sprintf("%s%s", to.URL, path), "application/json", body)
	// if err != nil {
	// 	log.Error(err)
	// }
	return resp, err
}

// getText
// HTTP GET the path, return the response as a string.
func (to *Session) getText(path string) (string, error) {

	body, err := to.getBytes(path)
	return string(body), err
}

// Login to traffic_ops, the response should set the cookie for this session
// automatically. Start with
// to := traffic_ops.Login("user", "passwd", true)
// subsequent calls like to.GetData("datadeliveryservice") will be authenticated.
func Login(toURL string, toUser string, toPasswd string, insecure bool) (*Session, error) {
	var to Session

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&options)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	to.UserAgent = &http.Client{
		Transport: tr,
		Jar:       jar,
	}

	credentials := Credentials{
		Username: toUser,
		Password: toPasswd,
	}

	jcreds, err := json.Marshal(credentials)

	if err != nil {
		log.Info(err)
		return &to, err
	}

	url := fmt.Sprintf("%s/api/1.1/user/login", toURL)
	resp, err := to.UserAgent.Post(url, "application/json", bytes.NewReader(jcreds))
	if err != nil {
		log.Info(err)
		return &to, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info(err)
		return &to, err
	}

	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Info(err)
		return &to, err
	}

	to.URL = toURL

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
		return &to, err
	}

	to.Cache = make(map[string]cacheentry)

	log.Infof("logged into %s!", toURL)
	return &to, err
}
