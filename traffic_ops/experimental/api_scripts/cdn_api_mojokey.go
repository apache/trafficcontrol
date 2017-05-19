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

//to run-> go run cdn_api_mojokey.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

//ERROR HANDLER -------------------------------------------------------------
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

//READ IN CREDENTIALS FROM FILE--------------------------------------------------------
func read_creds_file() (user, pw string) {
	path := os.Getenv("HOME") + "/Downloads/scripts/"
	file, err := os.Open(path + "CDN_API_Credentials.txt")
	checkError("Failed to open creds file", err)
	contents, err := ioutil.ReadAll(file) //read file to bytes
	checkError("Failed to read creds file", err)
	creds := strings.Split(string(contents), "\n") //create a slice of the contents, must convert from bytes to string
	file.Close()                                   //closes file
	user = creds[0]                                // extract user from file strings
	pw = creds[1]

	return
}

//API Request for Mojolicious Key-------------------------------------------------------------
func api_request(user, pw string) (mojo_key string) {
	user_pw := map[string]string{"p": pw, "u": user}                                // create a map of user and pw for inclusion in request
	user_pw_json, _ := json.Marshal(user_pw)                                        // convert user/password map to json structure for inclusion in request
	key_url := "https://cdnportal.comcast.net/api/1.2/user/login"                   // add key to url
	key_req, err := http.NewRequest("POST", key_url, bytes.NewBuffer(user_pw_json)) // post request to get mojolicious cookie
	key_req.Header.Add("Accept", "application/json")
	key_client := &http.Client{}
	key_resp, err := key_client.Do(key_req)
	checkError("Failed key http request: ", err)
	key_resp.Body.Close()

	//Extract Mojolicious cookie from response headers map
	pattern := regexp.MustCompile(`mojolicious=([A-Za-z0-9\-\_]+);`) //compile regex to extract cookie from response
	mojo_cookie := key_resp.Header["Set-Cookie"]                     // get cookie header from http response
	mojo_key = pattern.FindStringSubmatch(mojo_cookie[0])[1]         //extract cookie from response
	return
}

func main() {

	//GET CREDENTIALS FROM FILE IF NOT ENTERED ON COMMAND LINE
	fmt.Println("Getting user credentials from file...")
	user, pw := read_creds_file()
	fmt.Println("USER =", user)
	//GET MOJOLICIOUS KEY FROM API
	fmt.Println("Sending API call to get mojolicious key...")
	mojo_key := api_request(user, pw)
	fmt.Println("MOJO KEY =", mojo_key)

}
