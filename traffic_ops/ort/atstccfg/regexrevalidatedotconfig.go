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

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func GetConfigFileCDNRegexRevalidateDotConfig(cfg TCCfg, cdnNameOrID string) (string, error) {
	cdnName, err := GetCDNNameFromCDNNameOrID(cfg, cdnNameOrID)
	if err != nil {
		return "", errors.New("getting CDN name from '" + cdnNameOrID + "': " + err.Error())
	}

	toToolName, toURL, err := GetTOToolNameAndURLFromTO(cfg)
	if err != nil {
		return "", errors.New("getting global parameters: " + err.Error())
	}

	fileParamsWithoutProfiles, err := GetConfigFileParameters(cfg, atscfg.RegexRevalidateFileName)
	if err != nil {
		return "", errors.New("getting regexreval parameters: " + err.Error())
	}

	fileParams, err := TCParamsToParamsWithProfiles(fileParamsWithoutProfiles)
	if err != nil {
		return "", errors.New("unmarshalling regexreval parameters profiles: " + err.Error())
	}

	params := map[string][]string{}
	for _, param := range fileParams {
		if !util.StrInArray(param.ProfileNames, tc.GlobalProfileName) {
			continue // TODO add profile query params to TO endpoint
		}
		params[param.Name] = append(params[param.Name], param.Value)
	}

	allJobs, err := GetJobs(cfg) // TODO add cdn query param to jobs endpoint
	if err != nil {
		return "", errors.New("unmarshalling regexreval parameters profiles: " + err.Error())
	}

	cdn, err := GetCDN(cfg, cdnName)
	if err != nil {
		return "", errors.New("getting cdn '" + string(cdnName) + "': " + err.Error())
	}

	dses, err := GetCDNDeliveryServices(cfg, cdn.ID)
	if err != nil {
		return "", errors.New("getting delivery services: " + err.Error())
	}

	dsNames := map[string]struct{}{}
	for _, ds := range dses {
		if ds.XMLID == nil {
			log.Errorln("Regex Revalidate got Delivery Service from Traffic Ops with a nil xmlId! Skipping!")
			continue
		}
		dsNames[*ds.XMLID] = struct{}{}
	}

	jobs := []tc.Job{}
	for _, job := range allJobs {
		if _, ok := dsNames[job.DeliveryService]; !ok {
			continue
		}
		jobs = append(jobs, job)
	}

	txt := atscfg.MakeRegexRevalidateDotConfig(cdnName, params, toToolName, toURL, jobs)
	return txt, nil // TODO implement
}
