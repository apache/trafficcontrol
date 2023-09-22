package httpService

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/test/fakeOrigin/endpoint"
	"github.com/apache/trafficcontrol/v8/test/fakeOrigin/m3u8"
)

type httpEndpoint struct {
	ID                string
	OutputDir         string
	Type              endpoint.Type
	DiskID            string
	IsLive            bool
	DefaultHeaders    http.Header
	NoCache           bool
	HTTPPath          string
	ContentType       string
	LastTranscodeTime time.Time
}

type staticHTTPEndpoint struct {
	FilePath string
	NoCache  bool
	Headers  http.Header
}

// LiveM3U8MinFiles is used to govern the minimum number of files to be present in a live m3u8
const LiveM3U8MinFiles = 20 // TODO make configurable
// LiveM3U8MinDuration is used to govern the minimum number of seconds to be present in a live m2u8
const LiveM3U8MinDuration = time.Duration(40) * time.Second // TODO make configurable

func httpEndpointHandler(hend httpEndpoint) (http.HandlerFunc, error) {
	start := time.Now() // used to compute the live manifest offset
	diskpath := filepath.Join(hend.OutputDir, path.Base(hend.HTTPPath))
	data, err := ioutil.ReadFile(diskpath) // read even if we're not caching, so we can error out before serving
	if err != nil {
		return nil, errors.New("reading file '" + diskpath + "': " + err.Error())
	}

	return func(w http.ResponseWriter, r *http.Request) {
		for key, values := range hend.DefaultHeaders {
			w.Header().Set(key, values[0])
			for _, val := range values[1:] {
				w.Header().Add(key, val)
			}
		}

		if hend.NoCache {
			if data, err = ioutil.ReadFile(diskpath); err != nil {
				w.WriteHeader(http.StatusNotFound)
				// Remove the error message in production, this is just for debug purposes
				w.Write([]byte(strconv.Itoa(http.StatusNotFound) + " " + http.StatusText(http.StatusNotFound) + " - " + err.Error()))
				return
			}
		}

		if hend.IsLive {
			vodM3U8, err := m3u8.Parse(data)
			if err != nil {
				fmt.Println("Error reading file: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Cache-Control", "no-cache")
			if len(vodM3U8.VARs) > 0 {
				w.Header().Set("Content-Type", "application/x-mpegURL")
				w.Write(data)
				return
			}
			handleLiveM3U8(w, r, vodM3U8, LiveM3U8MinFiles, LiveM3U8MinDuration, start)
			return
		}

		w.Header().Set("Content-Type", hend.ContentType)
		w.Write(data)
	}, nil
}

func handleLiveM3U8(w http.ResponseWriter, r *http.Request, vod m3u8.M3U8, minFiles int64, minDuration time.Duration, start time.Time) {
	fmt.Println("path '" + r.URL.Path + "'")
	offset := time.Since(start)
	// TODO cache for 1s?
	live, err := m3u8.TransformVodToLive(vod, offset, minFiles, minDuration)
	if err != nil {
		fmt.Println("Error transforming : " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	liveBts := m3u8.SerializeLive(live)
	w.Header().Set("Content-Type", "application/x-mpegURL")
	w.Write(liveBts)
}

func staticHandler(e staticHTTPEndpoint) (http.HandlerFunc, error) {
	data, err := ioutil.ReadFile(e.FilePath) // read even if we're not caching, so we can error out before serving
	if err != nil {
		return nil, errors.New("reading file '" + e.FilePath + "': " + err.Error())
	}
	return func(w http.ResponseWriter, r *http.Request) {
		for key, values := range e.Headers {
			w.Header().Set(key, values[0])
			for _, val := range values[1:] {
				w.Header().Add(key, val)
			}
		}
		if e.NoCache {
			if data, err = ioutil.ReadFile(e.FilePath); err != nil {
				w.WriteHeader(http.StatusNotFound)
				// Remove the error message in production, this is just for debug purposes
				w.Write([]byte(strconv.Itoa(http.StatusNotFound) + " " + http.StatusText(http.StatusNotFound) + " - " + err.Error()))
				return
			}
		}
		w.Header().Set("Content-Type", http.DetectContentType(data)) // TODO add Content Type to config, to allow users to override detected type
		w.Write(data)
	}, nil

}

func crossdomainHandler(customCrossdomainFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf []byte
		if customCrossdomainFile != "" && strings.HasSuffix(customCrossdomainFile, ".xml") {
			_, ExistErr := os.Stat(customCrossdomainFile)
			if os.IsNotExist(ExistErr) {
				fmt.Println("Error Crossdomain File Does Not Exist")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var err error
			buf, err = ioutil.ReadFile(customCrossdomainFile)
			if err != nil {
				fmt.Println("Error Reading Crossdomain File : " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			tmp := ""
			err = xml.Unmarshal(buf, &tmp)
			if err != nil {
				fmt.Println("Error Crossdomain File is not valid XML: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.Write(buf)
	}
}
