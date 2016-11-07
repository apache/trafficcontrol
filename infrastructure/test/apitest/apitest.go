package apitest

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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// ApiTester is the main testing type. It has functions for all the data needed to test.
type ApiTester interface {
	FQDN() string
	ApiPath() string
	Cookies() []*http.Cookie
}

// GetClient gets a http client object.
// This exists to encapsulate TLS cert verification failure skipping.
// TODO(fix to not skip cert verification, when the api cert is valid)
func GetClient() *http.Client {
	return &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
}

// GetJSONEndpoint makes a GET request to the given endpoint, and returns the parsed JSON object.
func GetJSONEndpoint(t ApiTester, endpoint string) (interface{}, error) {
	uri := t.FQDN() + t.ApiPath() + endpoint
	req, err := http.NewRequest("GET", uri, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	for _, c := range t.Cookies() {
		req.AddCookie(c)
	}

	client := GetClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}

// GetJSONID makes a GET request to the given API endpoint,
// and returns the ID value of the result member with the given Name,
// or an error if no matching name or id exists.
func GetJSONID(t ApiTester, endpoint, name string) (int, error) {

	// nameIdResponse represents an API response of any
	// endpoint whose objects contain a name and string id.
	type nameIdResponse struct {
		Response []struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"response"`
	}

	uri := t.FQDN() + t.ApiPath() + endpoint
	req, err := http.NewRequest("GET", uri, strings.NewReader(""))
	if err != nil {
		return -1, err
	}
	for _, c := range t.Cookies() {
		req.AddCookie(c)
	}

	client := GetClient()
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var jsonData nameIdResponse
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return -1, err
	}

	for _, val := range jsonData.Response {
		if val.Name == name {
			return strconv.Atoi(val.Id)
		}
	}
	return -1, errors.New(name + " does not exist")
}

// TestJSONContains whether the response from the endpoint contains the keys in the expected parameter.
// It expects API endpoints to be of the form {"response": {"mykey": myjsonobject,"mykey2": myjsonobject2}}
// and asserts each key in `expected` exists in the returned response, and has the same value.
// That is, this may be used after POSTing values to assert they are returned with GET, while
// ignoring preexisting values.
func TestJSONContains(t ApiTester, endpoint string, expected map[string]interface{}) error {

	// contains checks whether val is contained in s, comparing with reflect.DeepEqual
	contains := func(s []interface{}, val interface{}) bool {
		for _, sVal := range s {
			if reflect.DeepEqual(sVal, val) {
				return true
			}
		}
		return false
	}

	jsonData, err := GetJSONEndpoint(t, endpoint)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(expected, jsonData) {
		return nil
	}

	jsonDataMap, ok := jsonData.(map[string]interface{})
	if !ok {
		return errors.New("Returned data was not of the expected format: top level value is not an object")
	}
	if _, ok := jsonDataMap["response"]; !ok {
		return errors.New("Returned data was not of the expected format: top level object does not contain a 'response' key")
	}
	jsonDataResponseMaps, ok := jsonDataMap["response"].([]interface{})
	if !ok {
		return errors.New("Returned data was not of the expected format: top level value 'response' member is not an array of objects")
	}

	if _, ok := expected["response"]; !ok {
		return errors.New("Expected data was not of the expected format: top level object does not contain a 'response' key")
	}
	expectedResponseMapSlice, ok := expected["response"].([]interface{})
	if !ok {
		return errors.New("Expected data was not of the expected format: top level value 'response' member is not an array of objects")
	}

	for _, expectedVal := range expectedResponseMapSlice {
		if !contains(jsonDataResponseMaps, expectedVal) {
			return fmt.Errorf("Response %v does not contain expected value '%v'\n", jsonDataResponseMaps, expectedVal)
		}
	}
	return nil
}

// TestJSONEqual makes a GET request to the given endpoint, and returns nil (success)
// if the received value matches the given expected value, or an error if not.
// Note the expected value must match the type returned by `json.Unmarshal` for `interface{}`,
// that is, JSON objects must be represented by map[string]interface{} and
// NOT e.g. map[string]string, even if all members are strings. JSON arrays must be []interface{}.
func TestJSONEqual(t ApiTester, endpoint string, expected map[string]interface{}) error {
	jsonData, err := GetJSONEndpoint(t, endpoint)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, jsonData) {
		expectedBytes, eerr := json.Marshal(expected)
		responseBytes, rerr := json.Marshal(jsonData)
		if eerr == nil && rerr == nil {
			return fmt.Errorf("ERROR:    %s\n  Expected: %v\n  Actual:   %v", endpoint, string(expectedBytes), string(responseBytes))
		}
		if eerr != nil {
			return fmt.Errorf("ERROR:    %s\n  Expected: %v\n  Actual:   %v", endpoint, expected, string(responseBytes))
		}
		if rerr != nil {
			return fmt.Errorf("ERROR:    %s\n  Expected: %v\n  Actual:   %v", endpoint, string(expectedBytes), jsonData)
		}
	}
	return nil
}

// DoPOST executes a POST query to set up test data.
// Note it queries FQDN+endpoint, not FQDN+ApiPath()+endpoint
func DoPOST(t ApiTester, endpoint string, dataMap map[string]string) error {
	data := url.Values{}
	for key, val := range dataMap {
		data.Add(key, val)
	}

	req, err := http.NewRequest("POST", t.FQDN()+endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, c := range t.Cookies() {
		req.AddCookie(c)
	}

	client := GetClient()
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	// contents, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	return nil
}
