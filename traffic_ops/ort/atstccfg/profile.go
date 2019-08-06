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
	"strconv"
)

func GetProfileNameFromProfileNameOrID(cfg TCCfg, profileNameOrID string) (string, error) {
	profileName := profileNameOrID
	if profileID, err := strconv.Atoi(profileNameOrID); err == nil {
		profile, err := GetProfile(cfg, profileID)
		if err != nil {
			return "", errors.New("getting profile '" + profileNameOrID + "': " + err.Error())
		}
		if profile.Name == "" {
			return "", errors.New("getting profile '" + profileNameOrID + "': got profile with empty name")
		}
		profileName = profile.Name
	}
	return profileName, nil
}
