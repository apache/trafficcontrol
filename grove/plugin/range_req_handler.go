package plugin

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
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/grove/web"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

type byteRange struct {
	Start int64
	End   int64
}

const MAXINT64 = 1<<63 - 1

type rangeRequestConfig struct {
	Mode              string `json:"mode"`
	MultiPartBoundary string // not in the json
}

func init() {
	AddPlugin(10000, Funcs{
		load:                rangeReqHandleLoad,
		onRequest:           rangeReqHandlerOnRequest,
		beforeCacheLookUp:   rangeReqHandleBeforeCacheLookup,
		beforeParentRequest: rangeReqHandleBeforeParent,
		beforeRespond:       rangeReqHandleBeforeRespond,
	})
}

// rangeReqHandleLoad loads the configuration
func rangeReqHandleLoad(b json.RawMessage) interface{} {
	cfg := rangeRequestConfig{}
	log.Errorf("rangeReqHandleLoad loading: %s\n", b)

	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Errorln("range_rew_handler  loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	if !(cfg.Mode == "get_full_serve_range" || cfg.Mode == "patch") {
		log.Errorf("Unknown mode for range_req_handler plugin: %s\n", cfg.Mode)
	}

	multipartBoundaryBytes := make([]byte, 16)
	if _, err := rand.Read(multipartBoundaryBytes); err != nil {
		log.Errorf("Error with rand.Read: %v\n", err)
	}
	// Use UUID format
	cfg.MultiPartBoundary = fmt.Sprintf("%x-%x-%x-%x-%x", multipartBoundaryBytes[0:4], multipartBoundaryBytes[4:6], multipartBoundaryBytes[6:8], multipartBoundaryBytes[8:10], multipartBoundaryBytes[10:])

	log.Debugf("range_rew_handler: load success: %+v\n", cfg)
	return &cfg
}

// rangeReqHandlerOnRequest determines if there is a Range header, and puts the ranges in *d.Context as a []byteRanges
func rangeReqHandlerOnRequest(icfg interface{}, d OnRequestData) bool {
	rHeader := d.R.Header.Get("Range")
	if rHeader == "" {
		log.Debugf("No Range header found\n")
		return false
	}
	log.Debugf("Range string is: %s\n", rHeader)
	// put the ranges [] in the context so we can use it later
	byteRanges := parseRangeHeader(rHeader)
	*d.Context = byteRanges
	return false
}

// rangeReqHandleBeforeCacheLookup is used to override the cacheKey when in store_ranges mode.
func rangeReqHandleBeforeCacheLookup(icfg interface{}, d BeforeCacheLookUpData) {
	cfg, ok := icfg.(*rangeRequestConfig)
	if !ok {
		log.Errorf("range_req_handler config '%v' type '%T' expected *rangeRequestConfig\n", icfg, icfg)
		return
	}
	if cfg.Mode == "store_ranges" {
		sep := "?"
		if strings.Contains(d.DefaultCacheKey, "?") {
			sep = "&"
		}
		newKey := d.DefaultCacheKey + sep + "grove_range_req_handler_plugin_data=" + d.Req.Header.Get("Range")
		d.CacheKeyOverrideFunc(newKey)
		log.Debugf("range_req_handler: store_ranges default key:%s, new key:%s\n", d.DefaultCacheKey, newKey)
	}
}

// rangeReqHandleBeforeParent changes the parent request if needed (mode == get_full_serve_range)
func rangeReqHandleBeforeParent(icfg interface{}, d BeforeParentRequestData) {
	log.Debugf("rangeReqHandleBeforeParent calling.")
	rHeader := d.Req.Header.Get("Range")
	if rHeader == "" {
		log.Debugln("No Range header found")
		return
	}
	log.Debugf("Range string is: %s\n", rHeader)
	cfg, ok := icfg.(*rangeRequestConfig)
	if !ok {
		log.Errorf("range_req_handler config '%v' type '%T' expected *rangeRequestConfig\n", icfg, icfg)
		return
	}
	if cfg.Mode == "get_full_serve_range" {
		// get_full_serve_range means get the whole thing from parent/org, but serve the requested range. Just remove the Range header from the upstream request
		d.Req.Header.Del("Range")
	}
	return
}

// rangeReqHandleBeforeRespond builds the 206 response
// Assume all the needed ranges have been put in cache before, which is the truth for "get_full_serve_range" mode which gets the whole object into cache.
// If mode == store_ranges, do nothing, we just return the object stored-as is
func rangeReqHandleBeforeRespond(icfg interface{}, d BeforeRespondData) {
	log.Debugf("rangeReqHandleBeforeRespond calling\n")
	ictx := d.Context
	ctx, ok := (*ictx).([]byteRange)
	if !ok {
		log.Errorf("Invalid context: %v\n", ictx)
	}
	if len(ctx) == 0 {
		return // there was no (valid) range header
	}

	cfg, ok := icfg.(*rangeRequestConfig)
	if !ok {
		log.Errorf("range_req_handler config '%v' type '%T' expected *rangeRequestConfig\n", icfg, icfg)
		return
	}
	if cfg.Mode == "store_ranges" {
		return // no need to do anything here.
	}

	// mode != store_ranges
	multipartBoundaryString := cfg.MultiPartBoundary
	multipart := false
	originalContentType := d.Hdr.Get("Content-type")
	*d.Hdr = web.CopyHeader(*d.Hdr) // copy the headers, we don't want to mod the cacheObj
	if len(ctx) > 1 {
		multipart = true
		multipartBoundaryString = cfg.MultiPartBoundary
		d.Hdr.Set("Content-Type", fmt.Sprintf("multipart/byteranges; boundary=%s", multipartBoundaryString))
	}
	totalContentLength, err := strconv.ParseInt(d.Hdr.Get("Content-Length"), 10, 64)
	if err != nil {
		log.Errorf("Invalid Content-Length header: %v\n", d.Hdr.Get("Content-Length"))
	}
	body := make([]byte, 0)
	for _, thisRange := range ctx {
		if thisRange.End == MAXINT64 || thisRange.End >= totalContentLength { // if the end range is "", or too large serve until the end
			thisRange.End = totalContentLength - 1
		}
		if thisRange.Start == -1 {
			thisRange.Start = totalContentLength - thisRange.End
			thisRange.End = totalContentLength - 1
		}

		rangeString := "bytes " + strconv.FormatInt(thisRange.Start, 10) + "-" + strconv.FormatInt(thisRange.End, 10)
		log.Debugf("range:%d-%d\n", thisRange.Start, thisRange.End)
		if multipart {
			body = append(body, []byte("\r\n--"+multipartBoundaryString+"\r\n")...)
			body = append(body, []byte("Content-type: "+originalContentType+"\r\n")...)
			body = append(body, []byte("Content-range: "+rangeString+"/"+strconv.FormatInt(totalContentLength, 10)+"\r\n\r\n")...)
		} else {
			d.Hdr.Add("Content-Range", rangeString+"/"+strconv.FormatInt(totalContentLength, 10))
		}
		bSlice := (*d.Body)[thisRange.Start : thisRange.End+1]
		body = append(body, bSlice...)
	}
	if multipart {
		body = append(body, []byte("\r\n--"+multipartBoundaryString+"--\r\n")...)
	}
	d.Hdr.Set("Content-Length", strconv.Itoa(len(body)))
	*d.Body = body
	*d.Code = http.StatusPartialContent
	return
}

func parseRange(rangeString string) (byteRange, error) {
	parts := strings.Split(rangeString, "-")

	var bRange byteRange
	if parts[0] == "" {
		bRange.Start = -1 // -1 means from the end
	} else {
		start, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			log.Errorf("Error converting rangeString start '%v' to numbers\n", rangeString)
			return byteRange{}, err
		}
		bRange.Start = start
	}
	if parts[1] == "" {
		bRange.End = MAXINT64 // MAXINT64 means till the end
	} else {
		end, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Errorf("Error converting rangeString end '%v' to numbers\n", rangeString)
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
		log.Errorf("Not a valid Range type: \"%s\"\n", rangeStringParts[0])
	}

	for _, thisRangeString := range strings.Split(rangeStringParts[1], ",") {
		thisRange, err := parseRange(thisRangeString)
		if err != nil {
			return nil
		}
		byteRanges = append(byteRanges, thisRange)
	}

	// if there is just one range, return, and don't incur the overhead of determining overlaps
	if len(byteRanges) <= 1 {
		return byteRanges
	}

	// Collapse overlapping byte range requests, first sort the array by Start
	sort.Slice(byteRanges, func(i, j int) bool {
		return byteRanges[i].Start < byteRanges[j].Start
	})

	// Then, copy ranges into collapsedRanges if applicable, collapse as needed
	collapsedRanges := make([]byteRange, 0)
	j := 0
	collapsedRanges = append(collapsedRanges, byteRanges[j])
	for i := 1; i < len(byteRanges); i++ {
		if collapsedRanges[j].End < byteRanges[i].Start {
			// Most normal case, the ranges are not overlapping; add the range to the collapsedRanges array
			collapsedRanges = append(collapsedRanges, byteRanges[i])
			j++
			continue
		}
		if collapsedRanges[j].Start <= byteRanges[i].Start && collapsedRanges[j].End >= byteRanges[i].End {
			// Don't add the entry at i, it is part of the entry at i-1
			continue
		}
		if collapsedRanges[j].Start <= byteRanges[i].Start && collapsedRanges[j].End < byteRanges[i].End {
			// Overlapping ranges, combine into one
			collapsedRanges[j].End = byteRanges[i].End
			continue
		}
	}

	return collapsedRanges
}
