package cfgfile

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
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func GetConfigFileProfileURLSigConfig(toData *TOData, fileName string) (string, error) {
	paramData := map[string]string{}
	// TODO add configFile query param to profile/parameters endpoint, to only get needed data
	for _, param := range toData.ServerParams {
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

	urlSigKeys, ok := toData.URLSigKeys[tc.DeliveryServiceName(dsName)]
	if !ok {
		return "", errors.New("no keys fetched for ds '" + dsName + "!")
	}

	return atscfg.MakeURLSigConfig(toData.Server.Profile, urlSigKeys, paramData, toData.TOToolName, toData.TOURL), nil
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
