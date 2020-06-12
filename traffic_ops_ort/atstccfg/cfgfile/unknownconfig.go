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
	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/config"
)

func GetConfigFileProfileUnknownConfig(toData *config.TOData, fileName string) (string, string, string, error) {
	inScope := false
	for _, scopeParam := range toData.ScopeParams {
		if scopeParam.ConfigFile != fileName {
			continue
		}
		if scopeParam.Value != "profiles" {
			continue
		}
		inScope = true
		break
	}
	if !inScope {
		return `{"alerts":[{"level":"error","text":"Error - incorrect file scope for route used.  Please use the servers route."}]}`, "", "", config.ErrBadRequest
	}
	params := ParamsToMap(FilterParams(toData.ServerParams, fileName, "", "", "location"))

	commentType := atscfg.GetUnknownConfigCommentType(toData.Server.Profile, params, toData.TOToolName, toData.TOURL)
	txt := atscfg.MakeUnknownConfig(toData.Server.Profile, params, toData.TOToolName, toData.TOURL)
	return txt, atscfg.ContentTypeUnknownConfig, commentType, nil
}
