package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// LoadX509KeyPair can also load certificate chain with intermediates
	cert, _ := tls.LoadX509KeyPair("../certs/client-intermediate-chain.crt.pem", "../certs/client.key.pem")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	client := http.Client{
		Timeout:   time.Second * 60,
		Transport: transport,
	}

	// Send standard username/password form combo
	// reqBody, err := json.Marshal(map[string]string{
	// 	"u": "userid",
	// 	"p": "exampleuseridpassword",
	// })
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	req, err := http.NewRequest(
		http.MethodPost,
		"https://server.local:8443/api/4.0/user/login",
		bytes.NewBufferString(""),
		// bytes.NewBuffer(reqBody), // username/password
	)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(respBody)) // Verify Success
	fmt.Println(resp.Cookies())   // Verify Cookie(s)
}
