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
	"errors"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
)

func GetConfigFileProfileURLSigConfig(cfg TCCfg, profileNameOrID string, fileName string) (string, error) {
	profileName, err := GetProfileNameFromProfileNameOrID(cfg, profileNameOrID)
	if err != nil {
		return "", errors.New("getting profile name from '" + profileNameOrID + "': " + err.Error())
	}

	profileParameters, err := GetProfileParameters(cfg, profileName)
	if err != nil {
		return "", errors.New("getting profile '" + profileName + "' parameters: " + err.Error())
	}
	if len(profileParameters) == 0 {
		// The TO endpoint behind toclient.GetParametersByProfileName returns an empty object with a 200, if the Profile doesn't exist.
		// So we act as though we got a 404 if there are no params (and there should always be storage.config params), to make ORT behave correctly.
		return "", ErrNotFound
	}

	toToolName, toURL, err := GetTOToolNameAndURLFromTO(cfg)
	if err != nil {
		return "", errors.New("getting global parameters: " + err.Error())
	}

	paramData := map[string]string{}
	// TODO add configFile query param to profile/parameters endpoint, to only get needed data
	for _, param := range profileParameters {
		if param.ConfigFile != fileName {
			continue
		}
		if param.Name == "location" {
			continue
		}
		paramData[param.Name] = param.Value
	}

	dsName := GetDSFromURLSigConfigFileName(fileName)
	if dsName == "" {
		// extra safety, this should never happen, the routing shouldn't get here
		return "", errors.New("getting ds name: malformed config file '" + fileName + "'")
	}

	urlSigKeys, err := GetURLSigKeys(cfg, dsName)
	if err != nil {
		return "", errors.New("getting url sig keys for ds '" + dsName + "': " + err.Error())
	}

	return atscfg.MakeURLSigConfig(profileName, urlSigKeys, paramData, toToolName, toURL), nil
}

// GetDSFromURLSigConfigFileName returns the DS of a URLSig config file name.
// For example, "url_sig_foobar.config" returns "foobar".
// If the given string is shorter than len("url_sig_a.config"), the empty string is returned.
func GetDSFromURLSigConfigFileName(fileName string) string {
	if !strings.HasPrefix(fileName, "url_sig_") || !strings.HasSuffix(fileName, ".config") || len(fileName) <= len("url_sig_")+len(".config") {
		return ""
	}
	fileName = fileName[len("url_sig_"):]
	fileName = fileName[:len(fileName)-len(".config")]
	return fileName
}
