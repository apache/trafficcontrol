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

const DropQStringDotConfigFileName = "drop_qstring.config"
const DropQStringDotConfigParamName = "content"

func MakeDropQStringDotConfig(
	profileName string,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	dropQStringVal *string, // value of the parameter name "content" configFile "drop_qstring.config"; nil if it doesn't exist
) string {
	text := GenericHeaderComment(profileName, toToolName, toURL)
	if dropQStringVal != nil {
		text += *dropQStringVal + "\n"
	} else {
		text += `/([^?]+) $s://$t/$1` + "\n"
	}
	return text
}
