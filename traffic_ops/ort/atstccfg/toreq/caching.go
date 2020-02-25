package toreq

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
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
)

// DefaultCacheFormat is the encoder, decoder, and file extension to use for caching.
// This may be changed at compile-time. See CacheFormatJSON.
// Note this is not used for all cache files. Notably, Delivery Service Servers use a custom CSV format, which is faster.
var DefaultCacheFormat = CacheFormatJSON

// GetCached attempts to get the given object from tempDir/cacheFileName.
// If the cache file doesn't exist, is too old, or is malformed, it uses getter to get the object, and stores it in cacheFileName.
// The object is placed in obj (which must be a pointer to the type of object to decode), and the error from getter is returned.
// The cache format is defined by CacheEncoder and CacheDecoder, which may be changed at compile-time.
func GetCached(cfg config.TCCfg, cacheFileName string, obj interface{}, getter func(obj interface{}) error) error {
	return GetCachedWithFormat(cfg, cacheFileName, obj, getter, DefaultCacheFormat)
}

// GetCachedDSS is like GetCached, but optimized for Delivery Service Servers, which is massive.
func GetCachedDSS(cfg config.TCCfg, cacheFileName string, obj *[]tc.DeliveryServiceServer, getter func(obj interface{}) error) error {
	return GetCachedWithFormat(cfg, cacheFileName, obj, getter, CacheFormatDSS)
}

func GetCachedWithFormat(cfg config.TCCfg, cacheFileName string, obj interface{}, getter func(obj interface{}) error, cacheFormat CacheFormat) error {
	cacheFileName += cacheFormat.Extension
	start := time.Now()
	err := ReadCache(cacheFormat.Decoder, cfg.TempDir, cacheFileName, cfg.CacheFileMaxAge, obj)
	if err == nil {
		log.Infof("ReadCache %v (hit) took %v\n", cacheFileName, time.Since(start).Round(time.Millisecond))
		return nil
	}

	log.Infoln("ReadCache failed to get object from '" + cfg.TempDir + "/" + cacheFileName + "', calling getter: " + err.Error())

	currentRetry := 0
	for {
		err := getter(obj)
		if err == nil {
			break
		}
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			// if the server returned a 404, retrying won't help
			return errors.New("getting uncached: " + err.Error())
		}
		if currentRetry == cfg.NumRetries {
			return errors.New("getting uncached: " + err.Error())
		}

		sleepSeconds := config.RetryBackoffSeconds(currentRetry)
		log.Warnf("getting '%v', sleeping for %v seconds: %v\n", cacheFileName, sleepSeconds, err)
		currentRetry++
		time.Sleep(time.Second * time.Duration(sleepSeconds)) // TODO make backoff configurable?
	}

	WriteCache(cacheFormat.Encoder, cfg.TempDir, cacheFileName, obj)
	log.Infof("GetCachedJSON %v (miss) took %v\n", cacheFileName, time.Since(start).Round(time.Millisecond))
	return nil
}

// WriteCache attempts to write obj to tempDir/cacheFileName.
// If there is an error, it is written to stderr but not returned.
func WriteCache(encode EncodeFunc, tempDir string, cacheFileName string, obj interface{}) {
	if encode == nil {
		log.Errorln("object '" + cacheFileName + "': nil encode func! Should never happen! Can't write file!")
		return
	}

	objPath := filepath.Join(tempDir, cacheFileName)
	// Use os.OpenFile, not os.Create, in order to set perms to 0600 - cookies allow login, therefore the file MUST only allow access by the current user, for security reasons
	objFile, err := os.OpenFile(objPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Errorln("creating object cache file '" + objPath + "': " + err.Error())
		return
	}
	defer objFile.Close()

	if err := encode(objFile, obj); err != nil {
		log.Errorln("writing '" + cacheFileName + "': " + err.Error())
		return
	}
}

// ReadCache gets obj from tempDir/cacheFileName, if it exists and isn't older than CacheFileMaxAge.
// The obj must be a non-nil pointer to the object to decode into.
func ReadCache(decode DecodeFunc, tempDir string, cacheFileName string, cacheFileMaxAge time.Duration, obj interface{}) error {
	if decode == nil {
		return errors.New("nil decode func")
	}
	objPath := filepath.Join(tempDir, cacheFileName)
	objFile, err := os.Open(objPath)
	if err != nil {
		return errors.New("opening object file '" + objPath + "':" + err.Error())
	}
	defer objFile.Close()

	objFileInfo, err := objFile.Stat()
	if err != nil {
		return errors.New("getting object file info '" + objPath + "':" + err.Error())
	}
	if objFileAge := time.Now().Sub(objFileInfo.ModTime()); objFileAge > cacheFileMaxAge {
		return fmt.Errorf("object file too old, max age %dms less than file age %dms", int(cacheFileMaxAge/time.Millisecond), int(objFileAge/time.Millisecond))
	}
	if err := decode(objFile, obj); err != nil {
		log.Errorln("GetJSONObjFromFile loaded '" + cacheFileName + "' but failed to parse JSON: " + err.Error())
		return errors.New("unmarshalling object from file '" + objPath + "':" + err.Error())
	}
	return nil
}

type CacheFormat struct {
	Encoder   EncodeFunc
	Decoder   DecodeFunc
	Extension string
}

var CacheFormatJSON = CacheFormat{JSONEncode, JSONDecode, JSONExtension}

// CacheFormatDSS is a special format encoder/decoder optimized for Delivery Service Servers.
// The encoder and decoder both return errors if obj is not a *[]tc.DeliveryServiceServer.
var CacheFormatDSS = CacheFormat{DSSEncode, DSSDecode, DSSExtension}

type EncodeFunc func(w io.Writer, obj interface{}) error
type DecodeFunc func(r io.Reader, obj interface{}) error

const JSONExtension = ".json"

func JSONDecode(r io.Reader, obj interface{}) error { return json.NewDecoder(r).Decode(obj) }
func JSONEncode(w io.Writer, obj interface{}) error { return json.NewEncoder(w).Encode(obj) }

func StringToCookies(cookiesStr string) []*http.Cookie {
	hdr := http.Header{}
	hdr.Add("Cookie", cookiesStr)
	req := http.Request{Header: hdr}
	return req.Cookies()
}

func CookiesToString(cookies []*http.Cookie) string {
	strs := []string{}
	for _, cookie := range cookies {
		strs = append(strs, cookie.String())
	}
	return strings.Join(strs, "; ")
}

// WriteCookiesFile writes the given cookies to the temp file. On error, returns nothing, but writes to stderr.
func WriteCookiesToFile(cookiesStr string, tempDir string) {
	cookiePath := filepath.Join(tempDir, config.TempCookieFileName)
	// Use os.OpenFile, not os.Create, in order to set perms to 0600 - cookies allow login, therefore the file MUST only allow access by the current user, for security reasons
	if cookieFile, err := os.OpenFile(cookiePath, os.O_RDWR|os.O_CREATE, 0600); err != nil {
		log.Errorln("creating cookie file '" + cookiePath + "': " + err.Error())
	} else {
		defer cookieFile.Close()
		if _, err := cookieFile.WriteString(cookiesStr + "\n"); err != nil {
			log.Errorln("writing cookie file '" + cookiePath + "': " + err.Error())
		}
	}
}

func GetCookiesFromFile(tempDir string, cacheFileMaxAge time.Duration) (string, error) {
	cookiePath := filepath.Join(tempDir, config.TempCookieFileName)

	cookieFile, err := os.Open(cookiePath)
	if err != nil {
		return "", errors.New("opening cookie file '" + cookiePath + "':" + err.Error())
	}
	defer cookieFile.Close()

	cookieFileInfo, err := cookieFile.Stat()
	if err != nil {
		return "", errors.New("getting cookie file info '" + cookiePath + "':" + err.Error())
	}

	cookieFileAge := time.Now().Sub(cookieFileInfo.ModTime())
	if cookieFileAge > cacheFileMaxAge {
		return "", fmt.Errorf("cookie file too old, max age %dms less than file age %dms", int(cacheFileMaxAge/time.Millisecond), int(cookieFileAge/time.Millisecond))
	}

	bts, err := ioutil.ReadAll(cookieFile)
	if err != nil {
		return "", errors.New("reading cookies from file '" + cookiePath + "':" + err.Error())
	}
	return string(bts), nil
}

// DSSExtension is the file extension for the format read and written by DSSEncode and DSSDecode.
const DSSExtension = ".csv"

// DSSEncode is an EncodeFunc optimized for Delivery Service Servers.
// If iObj is not a *[]tc.DeliveryServiceServer it immediately returns an error.
func DSSEncode(w io.Writer, iObj interface{}) error {
	// Please don't change this to use encoding/csv unless you benchmark and prove it's at least as fast. DSS is massive, and performance is important.
	obj, ok := iObj.(*[]tc.DeliveryServiceServer)
	if !ok {
		return fmt.Errorf("object is '%T' must be a *[]tc.DeliveryServiceServer\n", iObj)
	}
	for _, dss := range *obj {
		if dss.DeliveryService == nil || dss.Server == nil {
			continue // TODO warn?
		}
		if _, err := io.WriteString(w, strconv.Itoa(*dss.DeliveryService)+`,`+strconv.Itoa(*dss.Server)+"\n"); err != nil {
			return fmt.Errorf("writing object cache file: " + err.Error())
		}
	}
	return nil
}

// DSSDecode is a DecodeFunc optimized for Delivery Service Servers.
// If iObj is not a *[]tc.DeliveryServiceServer it immediately returns an error.
func DSSDecode(r io.Reader, iObj interface{}) error {
	// Please don't change this to use encoding/csv unless you benchmark and prove it's at least as fast. DSS is massive, and performance is important.
	obj, ok := iObj.(*[]tc.DeliveryServiceServer)
	if !ok {
		return fmt.Errorf("object is '%T' must be a *[]tc.DeliveryServiceServer\n", iObj)
	}
	objFileReader := bufio.NewReader(r)
	for {
		dsStr, err := objFileReader.ReadString(',')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return errors.New("malformed object file:" + err.Error())
		}
		dsStr = dsStr[:len(dsStr)-1] // ReadString(',') includes the comma; remove it

		ds, err := strconv.Atoi(dsStr)
		if err != nil {
			return errors.New("malformed object file: first field should be ds id, but is not an integer: '" + dsStr + "'")
		}
		svStr, err := objFileReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return errors.New("malformed object file: final line had ds but not server, or file is missing a trailing newline")
			}
			return errors.New("malformed object file:" + err.Error())
		}
		svStr = svStr[:len(svStr)-1] // ReadString('\n') includes the newline; remove it
		sv, err := strconv.Atoi(svStr)
		if err != nil {
			return errors.New("malformed object file: second field should be server id, but is not an integer: '" + svStr + "'")
		}
		*obj = append(*obj, tc.DeliveryServiceServer{DeliveryService: &ds, Server: &sv})
	}
	return nil
}
