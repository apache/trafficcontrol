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

func GetConfigFileProfileURISigningConfig(toData *TOData, fileName string) (string, string, error) {
	dsName := GetDSFromURISigningConfigFileName(fileName)
	if dsName == "" {
		// extra safety, this should never happen, the routing shouldn't get here
		return "", "", errors.New("getting ds name: malformed config file '" + fileName + "'")
	}

	uriSigningKeys, ok := toData.URISigningKeys[tc.DeliveryServiceName(dsName)]
	if !ok {
		return "", "", errors.New("no keys fetched for ds '" + dsName + "!")
	}

	return atscfg.MakeURISigningConfig(uriSigningKeys), atscfg.ContentTypeURISigningDotConfig, nil
}

// GetDSFromURISigningConfigFileName returns the DS of a URI Signing config file name.
// For example, "uri_signing_foobar.config" returns "foobar".
// If the given string is shorter than len("uri_signing_a.config"), the empty string is returned.
func GetDSFromURISigningConfigFileName(fileName string) string {
	if !strings.HasPrefix(fileName, "uri_signing_") || !strings.HasSuffix(fileName, ".config") || len(fileName) <= len("uri_signing_")+len(".config") {
		return ""
	}
	fileName = fileName[len("uri_signing_"):]
	fileName = fileName[:len(fileName)-len(".config")]
	return fileName
}
