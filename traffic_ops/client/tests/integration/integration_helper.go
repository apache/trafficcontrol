package integration

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

var (
	to *traffic_ops.Session
)

func init() {
	toURL := flag.String("toURL", "http://localhost:3000", "Traffic Ops URL")
	toUser := flag.String("toUser", "admin", "Traffic Ops user")
	toPass := flag.String("toPass", "password", "Traffic Ops password")
	flag.Parse()
	var loginErr error
	to, loginErr = traffic_ops.Login(*toURL, *toUser, *toPass, true)
	if loginErr != nil {
		fmt.Printf("\nError logging in to %v: %v\nMake sure toURL, toUser, and toPass flags are included and correct.\nExample:  go test -toUser=admin -toPass=pass -toURL=http://localhost:3000\n\n", *toURL, loginErr)
		os.Exit(1)
	}
}

//GetCdn returns a Cdn struct
func GetCdn() (traffic_ops.CDN, error) {
	cdns, err := to.CDNs()
	if err != nil {
		return *new(traffic_ops.CDN), err
	}
	cdn := cdns[0]
	if cdn.Name == "ALL" {
		cdn = cdns[1]
	}
	return cdn, nil
}

//GetProfile returns a Profile Struct
func GetProfile() (traffic_ops.Profile, error) {
	profiles, err := to.Profiles()
	if err != nil {
		return *new(traffic_ops.Profile), err
	}
	return profiles[0], nil
}

//Request sends a request to TO and returns a response.
//This is basically a copy of the private "request" method in the traffic_ops.go \
//but I didn't want to make that one public.
func Request(to traffic_ops.Session, method, path string, body []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", to.URL, path)

	var req *http.Request
	var err error

	if body != nil && method != "GET" {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
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

	if resp.StatusCode != http.StatusOK {
		e := traffic_ops.HTTPError{
			HTTPStatus:     resp.Status,
			HTTPStatusCode: resp.StatusCode,
			URL:            url,
		}
		return nil, &e
	}

	return resp, nil
}
