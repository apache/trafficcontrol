/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

/**
 *
 */
export const enum ProfileType {
	/**
	 *
	 */
	ATS_PROFILE = "ATS_PROFILE",
	/**
	 *
	 */
	TR_PROFILE = "TR_PROFILE",
	/**
	 *
	 */
	TM_PROFILE = "TM_PROFILE",
	/**
	 *
	 */
	TS_PROFILE = "TS_PROFILE",
	/**
	 *
	 */
	TP_PROFILE = "TP_PROFILE",
	/**
	 *
	 */
	INFLUXDB_PROFILE = "INFLUXDB_PROFILE",
	/**
	 *
	 */
	RIAK_PROFILE = "RIAK_PROFILE",
	/**
	 *
	 */
	SPLUNK_PROFILE = "SPLUNK_PROFILE",
	/**
	 *
	 */
	DS_PROFILE = "DS_PROFILE",
	/**
	 *
	 */
	ORG_PROFILE = "ORG_PROFILE",
	/**
	 *
	 */
	KAFKA_PROFILE = "KAFKA_PROFILE",
	/**
	 *
	 */
	LOGSTASH_PROFILE = "LOGSTASH_PROFILE",
	/**
	 *
	 */
	ES_PROFILE = "ES_PROFILE",
	/**
	 *
	 */
	UNK_PROFILE = "UNK_PROFILE",
	/**
	 *
	 */
	GROVE_PROFILE = "GROVE_PROFILE"
}

/**
 *
 */
export interface Parameter {
	configFile: string;
	id?: number;
	lastUpdated?: Date | null;
	name: string;
	profiles: null | Array<string>;
	secure: boolean;
	value: string;
}

/**
 *
 */
export interface Profile {
	id?: number;
	lastUpdated?: Date;
	name: string;
	cdnName: string;
	cdn: number;
	routingDisabled: boolean;
	type: ProfileType;
	params?: Array<Parameter>;
}
