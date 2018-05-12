package plugin

import (
	"encoding/json"

	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/apache/incubator-trafficcontrol/grove/web"
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"net/http"
	"strconv"
	"strings"
)

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

type byteRange struct {
	Start int
	End   int
}

type rangeRequestConfig struct {
	Mode string `json::"mode"`
}

func init() {
	AddPlugin(10000, Funcs{load: rangeReqHandleLoad, onRequest: rangeReqHandlerOnRequest, beforeParentRequest: rangeReqHandleBeforeParent, beforeRespond: rangeReqHandleBeforeRespond})
}

// rangeReqHandleLoad loads the configuration
func rangeReqHandleLoad(b json.RawMessage) interface{} {
	cfg := rangeRequestConfig{}
	log.Errorf("rangeReqHandleLoad loading: %s", b)

	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Errorln("range_rew_handler  loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	log.Debugf("range_rew_handler: load success: %+v\n", cfg)
	return &cfg
}

// rangeReqHandlerOnRequest determines if there is a Range header, and puts the ranges in *d.Context as a []byteRanges
func rangeReqHandlerOnRequest(icfg interface{}, d OnRequestData) bool {
	rHeader := d.R.Header.Get("Range")
	if rHeader == "" {
		log.Debugf("No Range header found")
		return false
	}
	log.Debugf("Range string is: %s", rHeader)
	byteRanges := parseRangeHeader(rHeader)
	*d.Context = byteRanges
	return false
}

// rangeReqHandleBeforeParent changes the parent request if needed (mode == getfull)
func rangeReqHandleBeforeParent(icfg interface{}, d BeforeParentRequestData) {
	log.Debugf("rangeReqHandleBeforeParent calling.")
	rHeader := d.Req.Header.Get("Range")
	if rHeader == "" {
		log.Debugf("No Range header found")
		return
	}
	log.Debugf("Range string is: %s", rHeader)
	cfg, ok := icfg.(*rangeRequestConfig)
	if !ok {
		log.Errorf("range_req_handler config '%v' type '%T' expected *rangeRequestConfig\n", icfg, icfg)
		return
	}
	if cfg.Mode == "getfull" {
		// getfull means get the whole thing from parent/org, but serve the requested range. Just remove the Range header from the upstream request, and put the range header in the context so we can use it in the response.
		d.Req.Header.Del("Range")
	}
	return
}

//--
//> GET /10Mb.txt HTTP/1.1
//> Host: localhost
//> Range: bytes=0-1,500000-500009
//> User-Agent: curl/7.54.0
//> Accept: */*
//>
//< HTTP/1.1 206 Partial Content
//< Date: Sat, 12 May 2018 14:21:56 GMT
//< Server: Apache/2.4.29 (Unix)
//< Last-Modified: Sat, 12 May 2018 14:19:11 GMT
//< ETag: "a00000-56c02efb675c0"
//< Accept-Ranges: bytes
//< Content-Length: 216
//< Cache-Control: max-age=60, public
//< Expires: Sat, 12 May 2018 14:22:56 GMT
//< Content-Type: multipart/byteranges; boundary=4d583e1cee600e82
//<
//
//--4d583e1cee600e82
//Content-type: text/plain
//Content-range: bytes 0-1/10485760
//
//00
//--4d583e1cee600e82
//Content-type: text/plain
//Content-range: bytes 500000-500009/10485760
//
//000500000
//
//--4d583e1cee600e82--
//--

// https://httpwg.org/specs/rfc7230.html#message.body.length

// rangeReqHandleBeforeRespond builds the 206 response
// assume all the needed ranges have been put in cache before, which is the truth for "getfull" mode which gets the whole object into cache.
func rangeReqHandleBeforeRespond(icfg interface{}, d BeforeRespondData) {
	log.Debugf("rangeReqHandleBeforeRespond calling\n")
	ictx := d.Context
	ctx, ok := (*ictx).([]byteRange)
	if !ok {
		log.Errorf("Invalid context: %v", ictx)
	}
	if len(ctx) == 0 {
		return // there was no (valid) range header
	}

	multipartBoundaryString := ""
	originalContentType := d.Hdr.Get("Content-type")
	*d.Hdr = web.CopyHeader(*d.Hdr) // copy the headers, we don't want to mod the cacheObj
	if len(ctx) > 1 {
		//multipart = true
		multipartBoundaryBytes := make([]byte, 8)
		if _, err := rand.Read(multipartBoundaryBytes); err != nil {
			log.Errorf("Error with rand.Read: %v", err)
		}
		multipartBoundaryString = hex.EncodeToString(multipartBoundaryBytes)
		d.Hdr.Set("Content-Type", fmt.Sprintf("multipart/byteranges; boundary=%s", multipartBoundaryString))
	}
	totalContentLength, err := strconv.Atoi(d.Hdr.Get("Content-Length"))
	if err != nil {
		log.Errorf("Invalid Content-Length header: %v", d.Hdr.Get("Content-Length"))
	}
	body := make([]byte, 0)
	for _, thisRange := range ctx {
		if thisRange.End == -1 {
			thisRange.End = totalContentLength
		}
		log.Debugf("range:%d-%d", thisRange.Start, thisRange.End)
		if multipartBoundaryString != "" {
			body = append(body, []byte(fmt.Sprintf("\r\n--%s\r\n", multipartBoundaryString))...)
			body = append(body, []byte(fmt.Sprintf("Content-type: %s\r\n", originalContentType))...)
			body = append(body, []byte(fmt.Sprintf("Content-range: bytes %d-%d/%d\r\n\r\n", thisRange.Start, thisRange.End, totalContentLength))...)
		} else {
			byteRangeString := fmt.Sprintf("bytes %d-%d/%d", thisRange.Start, thisRange.End, totalContentLength)
			d.Hdr.Add("Content-Range", byteRangeString)
		}
		bSlice := (*d.Body)[thisRange.Start : thisRange.End+1]
		body = append(body, bSlice...)
	}
	if multipartBoundaryString != "" {
		body = append(body, []byte(fmt.Sprintf("\r\n--%s--\r\n", multipartBoundaryString))...)
	}
	d.Hdr.Set("Content-Length", fmt.Sprintf("%d", len(body)))
	*d.Body = body
	*d.Code = http.StatusPartialContent
	return
}

func parseRange(rangeString string) (byteRange, error) {
	parts := strings.Split(rangeString, "-")

	var bRange byteRange
	if parts[0] == "" {
		bRange.Start = 0
	} else {
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Errorf("Error converting rangeString \"%\" to numbers", rangeString)
			return byteRange{}, err
		}
		bRange.Start = start
	}
	if parts[1] == "" {
		bRange.End = -1 // -1 means till the end
	} else {
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Errorf("Error converting rangeString \"%\" to numbers", rangeString)
			return byteRange{}, err
		}
		bRange.End = end
	}
	return bRange, nil
}

func parseRangeHeader(rHdrVal string) []byteRange {
	byteRanges := make([]byteRange, 0)
	rangeStringParts := strings.Split(rHdrVal, "=")
	if rangeStringParts[0] != "bytes" {
		log.Errorf("Not a valid Range type: \"%s\"", rangeStringParts[0])
	}

	for _, thisRangeString := range strings.Split(rangeStringParts[1], ",") {
		log.Debugf("bRangeStr: %s", thisRangeString)
		thisRange, err := parseRange(thisRangeString)
		if err != nil {
			return nil
		}
		byteRanges = append(byteRanges, thisRange)
	}
	return byteRanges
}
