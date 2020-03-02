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
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func GetConfigFileCDNRegexRevalidateDotConfig(toData *TOData) (string, string, error) {
	params := map[string][]string{}
	for _, param := range toData.GlobalParams {
		if param.ConfigFile != atscfg.RegexRevalidateFileName {
			continue
		}
		params[param.Name] = append(params[param.Name], param.Value)
	}

	dsNames := map[string]struct{}{}
	for _, ds := range toData.DeliveryServices {
		if ds.XMLID == nil {
			log.Errorln("Regex Revalidate got Delivery Service from Traffic Ops with a nil xmlId! Skipping!")
			continue
		}
		dsNames[*ds.XMLID] = struct{}{}
	}

	jobs := []tc.Job{}
	for _, job := range toData.Jobs {
		if _, ok := dsNames[job.DeliveryService]; !ok {
			continue
		}
		jobs = append(jobs, job)
	}

	return atscfg.MakeRegexRevalidateDotConfig(tc.CDNName(toData.Server.CDNName), params, toData.TOToolName, toData.TOURL, jobs), atscfg.ContentTypeRegexRevalidateDotConfig, nil
}
