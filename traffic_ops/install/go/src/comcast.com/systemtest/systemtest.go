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

package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
)

const (
	LANDING_PAGE_TITLE    = "Edge Health"
	PROFILES_PAGE_TITLE   = "Profiles"
	PARAMETERS_PAGE_TITLE = "Parameters"
)

var verbose bool

//TODO
// Command-line switches, default host to localhost and verbose to false
// Prompt for password if not supplied
// Optionally use generated cert

func main() {
	var host string
	var username string
	var password string
	if len(os.Args) != 5 {
		fmt.Println("Usage ./systemtest host username password verbose-flag")
		os.Exit(1)
	}

	host = os.Args[1]
	username = os.Args[2]
	password = os.Args[3]
	var err error
	verbose, err = strconv.ParseBool(os.Args[4])
	if err != nil {
		fmt.Println("Usage ./systemtest host username password verbose-flag")
		os.Exit(1)
	}

	if verbose {
		fmt.Println("host:", host, "username:", username, "password:", password, "verbose:", verbose)
	}

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr, Jar: jar}

	site := "https://" + host

	testOkay := true

	err = Login(client, site, username, password)
	if err == nil {
		fmt.Println("Login:           Ok")
	} else {
		testOkay = false
		fmt.Printf("Login:           %s\n", err)
	}

	err = GetProfiles(client, site)
	if err == nil {
		fmt.Println("Get Profiles:    Ok")
	} else {
		testOkay = false
		fmt.Printf("Get Profiles:    %s\n", err)
	}

	err = GetParameters(client, site)
	if err == nil {
		fmt.Println("Get Parameters:  Ok")
	} else {
		testOkay = false
		fmt.Printf("Get Parameters:  %s\n", err)
	}

	if !testOkay {
		os.Exit(1)
	}
}

func Login(client *http.Client, site string, username string, password string) error {
	data := url.Values{}
	data.Set("u", username)
	data.Add("p", password)

	url, err := url.ParseRequestURI(site)
	if err != nil {
		return err
	}

	url.Path = "/login/"

	header := make(http.Header)
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	respBody, err := DoRequest(client, url, "POST", header, data)
	if err != nil {
		return err
	}

	defer respBody.Close()

	doc, err := html.Parse(respBody)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Println("HTML doc parsed ok", "type:", doc.Type, "data:", doc.Data)
	}

	err = CheckHtml(doc, LANDING_PAGE_TITLE)
	if err != nil {
		return err
	}

	return nil
}

func GetProfiles(client *http.Client, site string) error {
	url, err := url.ParseRequestURI(site)
	if err != nil {
		return err
	}

	url.Path = "/profiles"

	respBody, err := DoRequest(client, url, "GET", nil, nil)
	if err != nil {
		return err
	}

	url.Path = "/profiles"

	defer respBody.Close()

	doc, err := html.Parse(respBody)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Println("HTML doc parsed ok", "type:", doc.Type, "data:", doc.Data)
	}

	err = CheckHtml(doc, PROFILES_PAGE_TITLE)
	if err != nil {
		return err
	}

	return nil
}

func GetParameters(client *http.Client, site string) error {
	url, err := url.ParseRequestURI(site)
	if err != nil {
		return err
	}

	url.Path = "/parameters/profile/all"

	respBody, err := DoRequest(client, url, "GET", nil, nil)
	if err != nil {
		return err
	}

	defer respBody.Close()

	doc, err := html.Parse(respBody)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Println("HTML doc parsed ok", "type:", doc.Type, "data:", doc.Data)
	}

	err = CheckHtml(doc, PARAMETERS_PAGE_TITLE)
	if err != nil {
		return err
	}

	return nil
}

func DoRequest(client *http.Client, url *url.URL, method string, header http.Header, data url.Values) (io.ReadCloser, error) {
	urlStr := fmt.Sprintf("%v", url)

	request, err := http.NewRequest(method, urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	if header != nil {
		request.Header = header
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if verbose {
		fmt.Println("resp:", resp)
		fmt.Println("status:", resp.Status)
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Response status code = %d", resp.StatusCode))
	}

	return resp.Body, nil
}

//TODO Make this better before JvD sees it
func CheckHtml(doc *html.Node, pageTitle string) error {
	gotTitle := false
	gotBody := false

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode && n.FirstChild.Data == pageTitle {
				gotTitle = true
			}
		} else if gotTitle && n.Type == html.ElementNode && n.Data == "body" {
			gotBody = true
		}
		if verbose {
			fmt.Println("node", "type:", n.Type, "data:", n.Data)
		}
		if !gotTitle || !gotBody {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
				if gotBody {
					break
				}
			}
		}
	}
	f(doc)

	if !gotTitle || !gotBody {
		return errors.New("Could not locate expected page title and body")
	}

	return nil
}
