package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const LOGIN string = "/login"
const USERS string = "/users/"

type TokenResponse struct {
	Token string
}

var tokenStr string
var urlStart string

func login(client *http.Client) {
	var jsonStr = []byte(`{"username":"jvd", "password": "secret"}`)
	req, err := http.NewRequest("POST", urlStart+LOGIN, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		log.Println(err)
		return
	}
	tokenStr = "Bearer " + tokenResp.Token
}

func createUSer(client *http.Client) {
	var jsonStr = []byte(`{"username":"jvdtest123", "password": "secret"}`)
	req, err := http.NewRequest("POST", urlStart+USERS, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Creating user", resp.StatusCode, resp.Status, string(body))
}

func editUSer(client *http.Client, userName string) {
	var jsonStr = []byte(`{"password": "secret1212changed"}`)
	req, err := http.NewRequest("PUT", urlStart+USERS+userName, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Editing user", resp.StatusCode, resp.Status, string(body))
}

func deleteUSer(client *http.Client, userName string) {
	req, err := http.NewRequest("DELETE", urlStart+USERS+userName, nil)
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Deleting user", resp.StatusCode, resp.Status, string(body))
}

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	urlStart = "https://localhost:9000"
	login(client)
	log.Println("Token:" + tokenStr)
	log.Print("EXPECT 200 OK")
	createUSer(client)
	log.Print("EXPECT ERROR:")
	createUSer(client)
	log.Print("EXPECT 200 OK")
	editUSer(client, "jvdtest123")
	log.Print("EXPECT 0 ROWS:")
	editUSer(client, "notthere")
	log.Print("EXPECT 200 OK")
	deleteUSer(client, "jvdtest123")
}
