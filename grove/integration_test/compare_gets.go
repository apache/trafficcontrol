package main

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

import (
	"flag"
	"fmt"
	"github.com/apache/trafficcontrol/v8/grove/web"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type responseType struct {
	Headers http.Header
	Body    []byte
}

func httpGet(URL, headers string) responseType {
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Println("ERROR in httpGet")
	}
	//log.Printf(">>>%v<<< %v\n", headers, len(strings.Split(headers, ".")))
	for _, hdrString := range strings.Split(headers, " ") {
		//log.Println(">>> ", hdrString)
		if hdrString == "" {
			continue
		}
		parts := strings.Split(hdrString, ":")
		if parts[0] == "Host" {
			req.Host = parts[1]
		} else {
			//log.Println("> ", parts)
			req.Header.Set(parts[0], parts[1])
		}
	}
	//log.Printf(">>>> %v", req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR in httpGet")
	}
	defer resp.Body.Close()
	var response responseType
	response.Headers = web.CopyHeader(resp.Header)
	response.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR in httpGet (readall)")
	}
	return response
}

func equalBodies(a, b []byte) bool {
	if a == nil || b == nil {
		return false
	}

	if a == nil && b == nil {
		return true
	}
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func equalStringSlices(a, b []string) bool {
	if a == nil || b == nil {
		return false
	}

	if a == nil && b == nil {
		return true
	}
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func inStringSlice(str string, arr []string) bool {
	for _, strEnt := range arr {
		if strEnt == str {
			return true
		}
	}
	return false
}

func compareResponses(response1 responseType, response2 responseType, ignoreHdrs []string, ignoreMPB bool) bool {

	if ignoreMPB {
		contentTypeHdr := response1.Headers.Get("Content-type")
		ignoreHdrs = append(ignoreHdrs, "Content-Type")   // the boundary will be different
		ignoreHdrs = append(ignoreHdrs, "Content-Length") // the boundary will be different
		//fmt.Println("ignoreing", contentTypeHdr, response1)
		if strings.HasPrefix(contentTypeHdr, "multipart/byteranges") {
			parts := strings.Split(contentTypeHdr, "=")
			MPBoundary := parts[1]
			//log.Println("+++")
			//log.Printf("%s\n", string(response1.Body))
			response1.Body = []byte(strings.Replace(string(response1.Body), MPBoundary, "", -1))
			//log.Printf("%s\n", response1.Body)
		}
		contentTypeHdr = response2.Headers.Get("Content-type")
		if strings.HasPrefix(contentTypeHdr, "multipart/byteranges") {
			parts := strings.Split(contentTypeHdr, "=")
			MPBoundary := parts[1]
			response2.Body = []byte(strings.Replace(string(response2.Body), MPBoundary, "", -1))
		}
	}
	if !equalBodies(response1.Body, response2.Body) {
		return false
	}
	for hdrKey, _ := range response1.Headers {
		if inStringSlice(hdrKey, ignoreHdrs) {
			continue
		}
		if !equalStringSlices(response1.Headers[hdrKey], response2.Headers[hdrKey]) {
			log.Printf("ERROR hdr %v doesn't match: \"%v\" != \"%v\"\n", hdrKey, response1.Headers[hdrKey], response2.Headers[hdrKey])
			return false
		}
		//fmt.Printf(">>>>> %v\n", hdrKey)
	}

	return true
}
func main() {
	originURL := flag.String("org", "http://localhost", "The origin URL (default: \"http://localhost\")")
	cacheURL := flag.String("cache", "http://localhost:8080", "The cache URL (default: \"http://localhost:8080\")")
	path := flag.String("path", "", "The path to GET")
	orgHdrs := flag.String("ohdrs", "", "Comma seperated list of headers to add to origin request")
	cacheHdrs := flag.String("chdrs", "", "Comma separated list of headers to add to cache request")
	ignoreHdrs := flag.String("ignorehdrs", "Server,Date", "Comma separated list of headers to ignore in the compare")
	ignoreMultiPartBoundary := flag.Bool("ignorempb", true, "Ignore multi part boundary in body comparison.")
	flag.Parse()

	resp := httpGet(*originURL+"/"+*path, *orgHdrs)
	cresp := httpGet(*cacheURL+"/"+*path, *cacheHdrs)
	if !compareResponses(resp, cresp, strings.Split(*ignoreHdrs, ","), *ignoreMultiPartBoundary) {
		fmt.Printf("FAIL: Body bytes don't match \n%s\n != \n%s\n", string(resp.Body), string(cresp.Body))
		os.Exit(1)

	}
	fmt.Println("PASS")
	os.Exit(0)
}
