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
 * ProfileType enumerates the allowable Types for Profiles.
 */
export const enum ProfileType {
	/**
	 * The type of a Profile used by a Cache Server (EDGE-tier or MID-tier).
	 */
	ATS_PROFILE = "ATS_PROFILE",
	/**
	 * The type of a Profile used by a Traffic Router.
	 */
	TR_PROFILE = "TR_PROFILE",
	/**
	 * The type of a Profile used by a Traffic Monitor.
	 */
	TM_PROFILE = "TM_PROFILE",
	/**
	 * The type of a Profile used by a Traffic Stats Server.
	 */
	TS_PROFILE = "TS_PROFILE",
	/**
	 * The type of a Profile used by a Traffic Portal.
	 */
	TP_PROFILE = "TP_PROFILE",
	/**
	 * The type of a Profile used by an InfluxDB Server.
	 */
	INFLUXDB_PROFILE = "INFLUXDB_PROFILE",
	/**
	 * The type of a Profile used by a Traffic Vault (called `RIAK` for legacy
	 * reasons).
	 */
	RIAK_PROFILE = "RIAK_PROFILE",
	/**
	 * The type of a Profile used by a Splunk Server.
	 */
	SPLUNK_PROFILE = "SPLUNK_PROFILE",
	/**
	 * The type of a Profile used by a Delivery Service.
	 */
	DS_PROFILE = "DS_PROFILE",
	/**
	 * The type of a Profile used by an Origin and/or Origin Server.
	 */
	ORG_PROFILE = "ORG_PROFILE",
	/**
	 * The type of a Profile used by a Kafka Server.
	 */
	KAFKA_PROFILE = "KAFKA_PROFILE",
	/**
	 * The type of a Profile used by a Logstash Server.
	 */
	LOGSTASH_PROFILE = "LOGSTASH_PROFILE",
	/**
	 * The type of a Profile used by an ElasticSearch Server.
	 */
	ES_PROFILE = "ES_PROFILE",
	/**
	 * The type of a Profile used by any type of Server not covered by some
	 * other Profile Type.
	 */
	UNK_PROFILE = "UNK_PROFILE",
	/**
	 * The type of a Profile used by a Grove Cache Server (EDGE-tier or
	 * MID-tier).
	 */
	GROVE_PROFILE = "GROVE_PROFILE"
}

/**
 * A Parameter is a piece of configuration for a Server or Delivery Service.
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
 * A Profile is a collection of configuration Parameters for a Server or
 * Delivery Service.
 */
export interface Profile {
	id: number;
	lastUpdated?: Date;
	description: string;
	name: string;
	cdnName: string;
	cdn: number;
	routingDisabled: boolean;
	type: ProfileType;
	params?: Array<Parameter>;
}
