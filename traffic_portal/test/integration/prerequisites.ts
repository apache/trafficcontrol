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

/** These are the Tenants the tests need. Parents need to be replaced with IDs. */
export const tenants = [
	{
		name: "tenantParent",
		parent: "root"
	},
	{
		name: "tenantSame",
		parent: "tenantParent"
	},
	{
		name: "tenantChild",
		parent: "tenantSame"
	},
	{
		name: "tenantDifferent",
		parent: "root"
	}
];

/** These are the Roles the tests need. */
export const roles = [
	{
		name: "admin",
		privLevel: 30
	},
	{
		name: "operations",
		privLevel: 20
	},
	{
		name: "read-only",
		privLevel: 10
	}
];

/**
 * These are the users needed by the tests.
 *
 * Their roles must be replaced with Role IDs, and their `tenant` property is
 * the name of a Tenant whose ID must be given in the POST that creates them in
 * the field `tenantId`. They are also missing the following properties which
 * can just be generated from the information given here:
 *
 * - `fullName`
 * - `email`
 * - `localPasswd`
 * - `confirmLocalPasswd`
 *
 * Password fields can be arbitrary as long as they are valid and match per
 * object.
 */
export const users = [
	{
		username: "TPAdmin",
		role: "admin",
		tenant: "tenantSame",
	},
	{
		username: "TPOperator",
		role: "operations",
		tenant: "tenantSame",
	},
	{
		username: "TPReadOnly",
		role: "read-only",
		tenant: "tenantSame",
	},
	{
		username: "TPAdminDiff",
		role: "admin",
		tenant: "tenantDifferent",
	},
	{
		username: "TPOperatorDiff",
		role: "operations",
		tenant: "tenantDifferent",
	},
	{
		username: "TPReadOnlyDiff",
		role: "read-only",
		tenant: "tenantDifferent",
	},
	{
		username: "TPAdminParent",
		role: "admin",
		tenant: "tenantParent",
	},
	{
		username: "TPOperatorParent",
		role: "operations",
		tenant: "tenantParent",
	},
	{
		username: "TPReadOnlyParent",
		role: "read-only",
		tenant: "tenantParent",
	},
	{
		username: "TPAdminChild",
		role: "admin",
		tenant: "tenantChild",
	},
	{
		username: "TPOperatorChild",
		role: "operations",
		tenant: "tenantChild",
	},
	{
		username: "TPReadOnlyChild",
		role: "read-only",
		tenant: "tenantChild",
	}
];

/** These are the CDNs the tests need. */
export const CDNs = [
	{
		name: "dummyCDN",
		domainName: "cdnp3"
	}
];

/** Allowed values for a Cache Group's `localizationMethods` entries. */
type LocalizationMethod = "DEEP_CZ" | "CZ" | "GEO";

/**
 * Represents a TO Cache Group. This exists to enforce the `type` and
 * `localizationMethods` properties being valid on members of the exported
 * `cacheGroups` array.
 */
interface CacheGroup {
	fallbackToClosest: boolean;
	latitude: number;
	localizationMethods:
		null |
		[] |
		[LocalizationMethod] |
		[LocalizationMethod, LocalizationMethod] |
		[LocalizationMethod, LocalizationMethod, LocalizationMethod];
	longitude: number;
	name: string;
	shortName: string;
	/**
	 * Because only server Types can be created, this Type must be one of the
	 * ones known to exist in new TO installations (though not guaranteed to
	 * exist in any TO installation first created before 4.0).
	 */
	type: "EDGE_LOC" | "MID_LOC" | "ORG_LOC" | "TR_LOC" | "TC_LOC";
}

/**
 * These are the Cache Groups the tests need.
 *
 * Types need to be replaced by IDs. Also, since non-server Type creation is not
 * allowed, the tests can't just create them if they don't exist. Any Types used
 * here MUST exist in Traffic Ops. Since new
 */
export const cacheGroups: Array<CacheGroup> = [
	{
		name: "testCG",
		shortName: "testCG",
		latitude: 0,
		longitude: 0,
		fallbackToClosest: true,
		localizationMethods: [
			"DEEP_CZ",
			"CZ",
			"GEO"
		],
		type: "EDGE_LOC"
	}
];

/**
 * Represents a TO Profile. This exists to enforce the `type` being valid on
 * members of the exported `profiles` array.
 */
interface Profile {
	name: string;
	description: string;
	cdn: string;
	type: "ATS_PROFILE" |
	"TR_PROFILE" |
	"TM_PROFILE" |
	"TS_PROFILE" |
	"TP_PROFILE" |
	"INFLUXDB_PROFILE" |
	"RIAK_PROFILE" |
	"SPLUNK_PROFILE" |
	"DS_PROFILE" |
	"ORG_PROFILE" |
	"KAFKA_PROFILE" |
	"LOGSTASH_PROFILE" |
	"ES_PROFILE" |
	"UNK_PROFILE" |
	"GROVE_PROFILE";
	routingDisabled: boolean;
}

/** These are the Profiles the tests need. CDNs need to be replaced with IDs. */
export const profiles: Array<Profile> = [
	{
		name: "testProfile",
		description: "A Profile used in testing",
		cdn: "dummyCDN",
		type: "ATS_PROFILE",
		routingDisabled: false
	}
];
