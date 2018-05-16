package main

import (
	"flag"
	"fmt"
	"github.com/apache/incubator-trafficcontrol/grove/web"
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
	for _, hdrString := range strings.Split(headers, ",") {
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

func equal(a, b []byte) bool {
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

func compareResponses(response1 responseType, response2 responseType, ignoreHdrs []string) bool {
	if !equal(response1.Body, response2.Body) {
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
	flag.Parse()

	resp := httpGet(*originURL+"/"+*path, *orgHdrs)
	cresp := httpGet(*cacheURL+"/"+*path, *cacheHdrs)
	if !compareResponses(resp, cresp, strings.Split(*ignoreHdrs, ",")) {
		fmt.Println("FAIL: Body bytes don't match")
		os.Exit(1)

	}
	fmt.Println("PASS")
	os.Exit(0)
}
