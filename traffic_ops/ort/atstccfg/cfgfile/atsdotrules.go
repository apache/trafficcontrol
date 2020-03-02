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
)

const ATSDotRulesFileName = StorageFileName

func GetConfigFileProfileATSDotRules(toData *TOData) (string, string, error) {
	paramData := map[string]string{}
	// TODO add configFile query param to profile/parameters endpoint, to only get needed data
	for _, param := range toData.ServerParams {
		if param.ConfigFile != ATSDotRulesFileName {
			continue
		}
		if param.Name == "location" {
			continue
		}
		paramData[param.Name] = param.Value
	}
	return atscfg.MakeATSDotRules(toData.Server.Profile, paramData, toData.TOToolName, toData.TOURL), atscfg.ContentTypeATSDotRules, nil
}
