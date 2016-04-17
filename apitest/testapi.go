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
const REGISTER string = "/register/"
const CAMERAS string = "/cameras"

type TokenResponse struct {
	Token string
}

var tokenStr string
var urlStart string

func login(client *http.Client) {
	var jsonStr = []byte(`{"username":"jvdtest123", "password": "secret"}`)
	req, err := http.NewRequest("POST", urlStart+LOGIN, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
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

func createUser(client *http.Client) {
	var jsonStr = []byte(`{"username":"jvdtest123", "password": "secret"}`)
	req, err := http.NewRequest("POST", urlStart+REGISTER, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("Creating user ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func editUser(client *http.Client, userName string) {
	var jsonStr = []byte(`{"password": "secret1212changed"}`)
	req, err := http.NewRequest("PUT", urlStart+USERS+userName, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("Editing user ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func deleteUser(client *http.Client, userName string) {
	req, err := http.NewRequest("DELETE", urlStart+USERS+userName, nil)
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("Deleting user ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func getUser(client *http.Client, userName string) {
	req, err := http.NewRequest("GET", urlStart+USERS+userName, nil)
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("Get user ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func getUsers(client *http.Client) {
	req, err := http.NewRequest("GET", urlStart+USERS, nil)
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("Get all users ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func createCamera(client *http.Client) {
	var jsonStr = []byte(`{"name":"livingroom", "owner": "jvdtest123", "url":"http://camera.com", "username":"cameo", "password": "secret"}`)
	req, err := http.NewRequest("POST", urlStart+CAMERAS+"/jvdtest123", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("Creating camera ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func editCamera(client *http.Client, cameraOwner string, cameraName string) {
	var jsonStr = []byte(`{"password": "camerasecret1212changed"}`)
	req, err := http.NewRequest("PUT", urlStart+CAMERAS+"/"+cameraOwner+"/"+cameraName, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("Editing camera ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func deleteCamera(client *http.Client, cameraOwner string, cameraName string) {
	req, err := http.NewRequest("DELETE", urlStart+CAMERAS+"/"+cameraOwner+"/"+cameraName, nil)
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("Deleting camera ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func getCamera(client *http.Client, cameraOwner string, cameraName string) {
	req, err := http.NewRequest("GET", urlStart+CAMERAS+"/"+cameraOwner+"/"+cameraName, nil)
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("Get camera ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func getCameras(client *http.Client, cameraOwner string) {
	req, err := http.NewRequest("GET", urlStart+CAMERAS+"/"+cameraOwner, nil)
	req.Header.Set("Authorization", tokenStr)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("Get all cameras ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	urlStart = "https://localhost:9000"
	// urlStart = "https://ec2-52-37-126-44.us-west-2.compute.amazonaws.com:9000"
	log.Println("API Server: " + urlStart)
	log.Print("EXPECT 200 OK")
	createUser(client)

	login(client)
	log.Println("Token:" + tokenStr)
	log.Print("EXPECT ERROR:")
	createUser(client)
	log.Print("EXPECT 200 OK")
	editUser(client, "jvdtest123")
	log.Print("EXPECT 200 OK (multiple items)")
	getUsers(client)
	log.Print("EXPECT 200 OK")
	getUser(client, "jvdtest123")
	log.Print("EXPECT 0 ROWS:")
	editUser(client, "notthere")
	log.Print("EXPECT 200 OK")
	deleteUser(client, "jvdtest123")

	log.Println()
	log.Print("EXPECT 200 OK")
	createUser(client)

	log.Print("EXPECT 200 OK")
	createCamera(client)
	log.Print("EXPECT ERROR:")
	createCamera(client)
	log.Print("EXPECT 200 OK")
	editCamera(client, "jvdtest123", "livingroom")
	log.Print("EXPECT 200 OK (multiple items)")
	getCameras(client, "jvdtest123")
	log.Print("EXPECT 200 OK")
	getCamera(client, "jvdtest123", "livingroom")
	log.Print("EXPECT 0 ROWS:")
	editCamera(client, "notthere", "livingroom")
	log.Print("EXPECT 200 OK")
	deleteCamera(client, "jvdtest123", "livingroom")

	log.Print("EXPECT 200 OK")
	deleteUser(client, "jvdtest123")
}
