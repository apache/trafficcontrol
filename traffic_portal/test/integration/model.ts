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

/** Represents a TO CDN. */
export interface CDN {
	readonly dnssecEnabled: boolean;
	readonly domainName: string;
	readonly id: number;
	readonly lastUpdated: string;
	readonly name: string;
}

/** Allowed values for a Cache Group's `localizationMethods` entries. */
type LocalizationMethod = "DEEP_CZ" | "CZ" | "GEO";

/**
 * The various array constructions allowed for a Cache Group's
 * `localizationMethods` property - this is as close as we can come with basic
 * typings to ensuring duplicate methods aren't given in requests.
 */
export type LocalizationMethodArray =
	[] |
	[LocalizationMethod] |
	[LocalizationMethod, LocalizationMethod] |
	[LocalizationMethod, LocalizationMethod, LocalizationMethod];

/** Represents a TO Cache Group. */
export interface CacheGroup {
	readonly fallbacks: Array<string>;
	readonly fallbackToClosest: boolean;
	readonly id: number;
	readonly lastUpdated: string;
	readonly latitude: number;
	readonly localizationMethods: LocalizationMethodArray;
	readonly longitude: number;
	readonly name: string;
	readonly parentCacheGroupID: null | number;
	readonly parentCacheGroupName: null | string;
	readonly secondaryParentCacheGroupID: null | number;
	readonly secondaryParentCacheGroupName: null | string;
	readonly shortName: string;
	readonly typeName: string;
	readonly typeId: number;
}

/** The allowable values for a Profile's `type`. */
export type ProfileType = "ATS_PROFILE" |
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

/** Represents a TO Profile */
export interface Profile {
	readonly cdn: number;
	readonly cdnName: string;
	readonly description: string;
	readonly id: number;
	readonly lastUpdated: string;
	readonly name: string;
	readonly routingDisabled: boolean;
	readonly type: ProfileType;
}

/** Represents a TO Role. */
export interface Role {
	readonly capabilities: Array<string> | null;
	readonly description: string;
	readonly id: number;
	readonly lastUpdated: string;
	readonly name: string;
	readonly privLevel: number;
}

/** Fields common to all Tenants. */
interface TenantBase {
	readonly id: number;
	readonly lastUpdated: string;
	readonly name: string;
	readonly parentId: number | null;
}

/** A model of a non-root Tenant. */
interface RegularTenant extends TenantBase {
	readonly parentId: number;
	readonly parentName: string;
}

/** A model of the root Tenant. */
interface RootTenant extends TenantBase {
	readonly name: "root";
	readonly parentId: null;
}

/** Represents a TO Tenant. **/
export type Tenant = RegularTenant | RootTenant;

/** Represents a TO Type. **/
export interface Type {
	readonly description: string;
	readonly id: number;
	readonly lastUpdated: string;
	readonly name: string;
	readonly useInTable: string;
}


/** Represents a TO User. **/
export interface User {
	readonly addressLine1: string | null;
	readonly addressLine2: string | null;
	readonly city: string | null;
	readonly company: string | null;
	readonly country: string | null;
	readonly email: string;
	readonly fullName: string;
	readonly gid: null;
	readonly id: number;
	/**
	 * In practice this should never be null, but it could be if the user was
	 * created before Traffic Ops version 6 and hasn't logged in since.
	 */
	readonly lastAuthenticated: string | null;
	readonly lastUpdated: string;
	readonly newUser: boolean;
	readonly phoneNumber: string | null;
	readonly postalCode: string | null;
	readonly publicSshKey: string | null;
	readonly registrationSent: string | null;
	readonly role: number;
	readonly rolename: string;
	readonly stateOrProvince: null;
	readonly tenant: string;
	readonly tenantId: number;
	readonly uid: null;
	readonly username: string;
}
