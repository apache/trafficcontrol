/*

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
	"errors"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type Session struct {
	UserName  string
	Password  string
	Url       string
	UserAgent *http.Client
	Cache     map[string]cacheentry
}

// {"alerts":[{"level":"success","text":"Successfully logged in."}],"version":"1.1"}
type Result struct {
	Alerts  []Alert
	Version float64 `json:"version"`
}

type Alert struct {
	Level string `json:"level"`
	Text  string `json:"text"`
}

type cacheentry struct {
	Entered int64
	bytes   []byte
}

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
	resp, err := to.UserAgent.Get(to.Url + path)
	if err != nil {
		log.Println(err)
		return body, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return body, err
}

func (to *Session) PostJson(path string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", to.Url+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := to.UserAgent.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	body, err = ioutil.ReadAll(resp.Body)
	return body, err
}

func (to *Session) PutJson(path string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("PUT", to.Url+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := to.UserAgent.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	return respBody, err
}

func (to *Session) Delete(path string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", to.Url+path, nil)
	resp, err := to.UserAgent.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	return respBody, err
}

// getText
// HTTP GET the path, return the response as a string.
func (to *Session) GetText(path string) (string, error) {

	body, err := to.getBytes(path)
	return string(body), err
}

// Login to traffic_ops, the response should set the cookie for this session
// automatically. Start with
// client := client.Login("user", "passwd", true)
// subsequent calls like to.GetData("datadeliveryservice") will be authenticated.
func Login(toUrl string, toUser string, toPasswd string, insecure bool) (*Session, error) {
	var to Session

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&options)
	if strings.HasPrefix(toUrl, "https://") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		}
		to.UserAgent = &http.Client{
			Transport: tr,
			Jar:       jar,
		}
	} else {
		to.UserAgent = &http.Client{
			Jar: jar,
		}
	}

	credentials := Credentials{
		Username: toUser,
		Password: toPasswd,
	}

	jcreds, err := json.Marshal(credentials)

	if err != nil {
		log.Println(err)
		return &to, err
	}

	resp, err := to.UserAgent.Post(toUrl+"/login", "application/json", bytes.NewReader(jcreds))
	if err != nil {
		log.Println(err)
		return &to, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return &to, err
	}

	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
		return &to, err
	}

	to.Url = toUrl

	success := false

	for _, alert := range result.Alerts {
		if alert.Level == "success" && alert.Text == "Successfully logged in." {
			success = true
			break
		}
	}

	if !success {
		fmt.Println("NO SUCCESS")
		err := errors.New("Login failed, result string: " + fmt.Sprintf("%+v", result))
		return &to, err
	}

	to.Cache = make(map[string]cacheentry)

	log.Println("logged into " + toUrl + "!")
	return &to, err
}
