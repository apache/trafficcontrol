package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const LOGIN string = "/login"
const USERS string = "/users/"
const REGISTER string = "/register/"
const CAMERAS string = "/cameras"
const CAMERACONTROL string = "/control/"
const FEED string = "/feed/"
const LIVESTREAM string = "/livestream"
const VIDEO string = "/video"

type TokenResponse struct {
	Token string
}

var tokenStr string
var urlStart string

func login(client *http.Client, userName string, password string) {
	var jsonStr = []byte(`{"username":"` + userName + `", "password": "` + password + `"}`)
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
	// req.Header.Set("Authorization", tokenStr)
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

func controlCameraPosition(client *http.Client, userName string, camerName string, actionString string) {
	req, err := http.NewRequest("POST", urlStart+CAMERACONTROL+userName+"/"+camerName+"?"+actionString, nil)
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
	log.Print("Control camera ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func controlCameraRecording(client *http.Client, userName string, camerName string, actionString string) {
	req, err := http.NewRequest(actionString, urlStart+FEED+userName+"/"+camerName, nil)
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
	log.Print("Control camera recording ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func liveStream(client *http.Client, cameraOwner string, cameraName string) {
	req, err := http.NewRequest("GET", urlStart+LIVESTREAM+"/"+cameraOwner+"/"+cameraName, nil)
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
	log.Print("Get live stream ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func getVideo(client *http.Client, cameraOwner string, cameraName string, start string, end string) {
	req, err := http.NewRequest("GET", urlStart+VIDEO+"/"+cameraOwner+"/"+cameraName+"?start="+start+"&stop="+end, nil)
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
	log.Print("Get video ", resp.StatusCode, " ", resp.Status, " ", string(body))
}

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	urlStart = "https://localhost:9000"
	// urlStart = "https://ec2-52-37-126-44.us-west-2.compute.amazonaws.com:9000"
	log.Println("API Server: " + urlStart)

	// clean up
	login(client, "root", "supersecret")
	deleteCamera(client, "jvdtest123", "livingroom")
	deleteUser(client, "jvdtest123")

	log.Print("1 EXPECT 200 OK")
	createUser(client)
	// getUsers(client)

	login(client, "jvdtest123", "secret")
	log.Println("Token:" + tokenStr)
	log.Print("2 EXPECT ERROR")
	createUser(client)
	log.Print("3 EXPECT 200 OK")
	editUser(client, "jvdtest123")
	log.Print("4 EXPECT 200 OK")
	getUser(client, "jvdtest123")
	log.Print("6 EXPECT 401 Not Authorized")
	editUser(client, "root")
	log.Print("9 EXPECT 200 OK")
	createCamera(client)
	log.Print("10 EXPECT ERROR")
	createCamera(client)
	log.Print("11 EXPECT 200 OK")
	editCamera(client, "jvdtest123", "livingroom")
	log.Print("12 EXPECT 200 OK (multiple items)")
	getCameras(client, "jvdtest123")
	log.Print("13 EXPECT 200 OK")
	getCamera(client, "jvdtest123", "livingroom")
	log.Print("14 EXPECT 401 Not Authorized")
	editCamera(client, "notthere", "livingroom")
	log.Println("15 Don't know what to expect")
	controlCameraPosition(client, "jvdtest123", "livingroom", "action=start&direction=Up&velocity=3")
	time.Sleep(3 * time.Second)
	log.Println("16 Don't know what to expect")
	controlCameraPosition(client, "jvdtest123", "livingroom", "action=start&direction=down&velocity=3")
	log.Println("17 expect 200 OK")
	controlCameraRecording(client, "jvdtest123", "livingroom", "POST")
	time.Sleep(3 * time.Second)
	log.Println("17 expect 200 OK")
	controlCameraRecording(client, "jvdtest123", "livingroom", "DELETE")
	log.Println("18 expect 200 OK")
	liveStream(client, "jvdtest123", "livingroom")
	t := time.Now()
	end := t.Add(-3600 * time.Second)
	start := t.Add(-2 * 3600 * time.Second)
	// log.Println("start:", start.Local().Format(time.RFC3339), " end:", end.Local().Format(time.RFC3339))
	log.Println("19 Not sure what to expect")
	getVideo(client, "jvdtest123", "livingroom", start.Format(time.RFC3339), end.Format(time.RFC3339))
	// log.Print("15 EXPECT 200 OK")
	// deleteCamera(client, "jvdtest123", "livingroom")
	// deleteUser(client, "jvdtest123")
}
