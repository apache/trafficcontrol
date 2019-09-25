package atscfg

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

const AstatsSeparator = "="
const AstatsFileName = "astats.config"

func MakeAStatsDotConfig(
	profileName string,
	paramData map[string]string, // GetProfileParamData(tx, profile.ID, AstatsFileName)
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
) string {
	hdr := GenericHeaderComment(profileName, toToolName, toURL)
	txt := GenericProfileConfig(paramData, AstatsSeparator)
	if txt == "" {
		txt = "\n" // If no params exist, don't send "not found," but an empty file. We know the profile exists.
	}
	txt = hdr + txt
	return txt
}
