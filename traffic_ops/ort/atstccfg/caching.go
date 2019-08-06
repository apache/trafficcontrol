package main

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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
)

// GetCachedJSON attempts to get the given object from tempDir/cacheFileName.
// If the cache file doesn't exist, is too old, or is malformed, it uses getter to get the object, and stores it in cacheFileName.
// The object is placed in obj (which must be a pointer to the type of object to decode from JSON), and the error from getter is returned.
func GetCachedJSON(cfg TCCfg, cacheFileName string, obj interface{}, getter func(obj interface{}) error) error {
	err := GetJSONObjFromFile(cfg.TempDir, cacheFileName, cfg.CacheFileMaxAge, obj)
	if err == nil {
		return nil
	}

	log.Infoln("GetCachedJSON failed to get object from '" + cfg.TempDir + "/" + cacheFileName + "', calling getter: " + err.Error())

	currentRetry := 0
	for {
		err := getter(obj)
		if err == nil {
			break
		}

		if currentRetry == cfg.NumRetries {
			return errors.New("getting uncached: " + err.Error())
		}

		sleepSeconds := RetryBackoffSeconds(currentRetry)
		log.Warnf("getting '%v', sleeping for %v seconds: %v\n", cacheFileName, sleepSeconds, err)
		currentRetry++
		time.Sleep(time.Second * time.Duration(sleepSeconds)) // TODO make backoff configurable?
	}

	WriteCacheJSON(cfg.TempDir, cacheFileName, obj)
	return nil
}

// WriteCacheJSON attempts to write obj to tempDir/cacheFileName.
// If there is an error, it is written to stderr but not returned.
func WriteCacheJSON(tempDir string, cacheFileName string, obj interface{}) {
	objBts, err := json.Marshal(obj)
	if err != nil {
		log.Errorln("serializing object '" + cacheFileName + "' to JSON: " + err.Error())
		return
	}

	objPath := filepath.Join(tempDir, cacheFileName)
	// Use os.OpenFile, not os.Create, in order to set perms to 0600 - cookies allow login, therefore the file MUST only allow access by the current user, for security reasons
	objFile, err := os.OpenFile(objPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Errorln("creating object cache file '" + objPath + "': " + err.Error())
		return
	}
	defer objFile.Close()

	if _, err := objFile.Write(objBts); err != nil {
		log.Errorln("writing object cache file '" + objPath + "': " + err.Error())
		return
	}
}

// GetJSONObjFromFile gets obj from tempDir/cacheFileName, if it exists and isn't older than CacheFileMaxAge.
// Just like with json.Unmarshal, obj must be a non-nil pointer to the object to decode into.
func GetJSONObjFromFile(tempDir string, cacheFileName string, cacheFileMaxAge time.Duration, obj interface{}) error {
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

	bts, err := ioutil.ReadAll(objFile)
	if err != nil {
		return errors.New("reading object from file '" + objPath + "':" + err.Error())
	}

	if err := json.Unmarshal(bts, obj); err != nil {
		return errors.New("unmarshalling object from file '" + objPath + "':" + err.Error())
	}

	return nil
}

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
	cookiePath := filepath.Join(tempDir, TempCookieFileName)
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
	cookiePath := filepath.Join(tempDir, TempCookieFileName)

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
