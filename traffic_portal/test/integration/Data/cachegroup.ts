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

export const cachegroups = {
	tests: [
		{
			testName: "Admin Role",
			logins: [
				{
					"username": "TPAdmin",
					"password": "pa$$word"
				}
			],
			create: [
				{
					description: "create a EDGE_LOC cachegroup with FailOver CacheGroup Field",
					name: "TP_Cache1",
					shortName: "TPC1",
					type: "EDGE_LOC",
					latitude: 0,
					longitude: 0,
					parentCacheGroup: "infrastructure",
					secondaryParentCG: "infrastructure",
					failoverCG: "albany-ny-usa",
					validationMessage: "cachegroup was created."
				},
				{
					description: "create multiple EDGE_LOC cachegroup",
					name: "TP_Cache2",
					shortName: "TPC2",
					type: "EDGE_LOC",
					latitude: "0",
					longitude: "0",
					parentCacheGroup: "infrastructure",
					secondaryParentCG: "infrastructure",
					failoverCG: "",
					validationMessage: "cachegroup was created."
				},
				{
					description: "create a MID_LOC cachegroup",
					name: "TP_Cache3",
					shortName: "TPC3",
					type: "MID_LOC",
					latitude: "0",
					longitude: "0",
					parentCacheGroup: "infrastructure",
					secondaryParentCG: "infrastructure",
					validationMessage: "cachegroup was created."
				}
			],
			update: [
				{
					description: "add more Failover Cache Groups to EDGE_LOC type cachegroup",
					name: "TP_Cache1",
					type: "EDGE_LOC",
					failoverCG: "TP_Cache2",
					validationMessage: "cachegroup was updated."
				},
				{
					description: "Validate cannot add cache group fallback if the cache group fall back is a different type",
					name:"TP_Cache1",
					type:"EDGE_LOC",
					failoverCG: "TP_Cache3"
				},
				{
					description: "Validate cannot add an empty cache group fall back",
					name:"TP_Cache1",
					type:"EDGE_LOC",
					failoverCG: " "
				},
				{
					description: "change type of the Cache Groups",
					name: "TP_Cache1",
					type: "MID_LOC",
					validationMessage: "cachegroup was updated."
				}
			],
			remove: [
				{
					description: "delete a cachegroup",
					name: "TP_Cache1",
					validationMessage: "cachegroup was deleted."
				},
				{
					description: "delete a cachegroup",
					name: "TP_Cache3",
					validationMessage: "cachegroup was deleted."
				}
			]
		},
		{
			testName: "ReadOnly Role",
			logins: [
				{
					"username": "TPReadOnly",
					"password": "pa$$word"
				}
			],
			create: [
				{
					description: "create a CacheGroup",
					name: "TP_Cache1",
					shortName: "TPC1",
					type: "EDGE_LOC",
					latitude: "0",
					longitude: "0",
					parentCacheGroup: "infrastructure",
					secondaryParentCG: "infrastructure",
					failoverCG: "albany-ny-usa",
					validationMessage: "Forbidden."
				}
			],
			update: [
				{
					description: "update CacheGroup",
					name: "TP_Cache2",
					type: "MID_LOC",
					validationMessage: "Forbidden."
				}
			],
			remove: [
				{
					description: "delete a cachegroup",
					name: "TP_Cache2",
					validationMessage: "Forbidden."
				}
			]
		},
		{
			testName: "Operation Role",
			logins: [
				{
					"username": "TPOperator",
					"password": "pa$$word"
				}
			],
			create: [
				{
					description: "create a EDGE_LOC cachegroup with FailOver CacheGroup Field",
					name: "TP_Cache4",
					shortName: "TPC4",
					type: "EDGE_LOC",
					latitude: "0",
					longitude: "0",
					parentCacheGroup: "infrastructure",
					secondaryParentCG: "infrastructure",
					failoverCG: "albany-ny-usa",
					validationMessage: "cachegroup was created."
				},
				{
					description: "create multiple EDGE_LOC cachegroup",
					name: "TP_Cache5",
					shortName: "TPC5",
					type: "EDGE_LOC",
					latitude: "0",
					longitude: "0",
					parentCacheGroup: "infrastructure",
					secondaryParentCG: "infrastructure",
					failoverCG: "",
					validationMessage: "cachegroup was created."
				},
				{
					description: "create a MID_LOC cachegroup",
					name: "TP_Cache6",
					shortName: "TPC6",
					type: "MID_LOC",
					latitude: "0",
					longitude: "0",
					parentCacheGroup: "infrastructure",
					secondaryParentCG: "infrastructure",
					validationMessage: "cachegroup was created."
				}
			],
			update: [
				{
					description: "add more Failover Cache Groups to EDGE_LOC type cachegroup",
					name: "TP_Cache4",
					type: "EDGE_LOC",
					failoverCG: "TP_Cache5",
					validationMessage: "cachegroup was updated."
				},
				{
					description: "Validate cannot add cache group fallback if the cache group fall back is a different type",
					name:"TP_Cache4",
					type:"EDGE_LOC",
					failoverCG: "TP_Cache6"
				},
				{
					description: "Validate cannot add an empty cache group fall back",
					name:"TP_Cache4",
					type:"EDGE_LOC",
					failoverCG: " "
				},
				{
					description: "change type of the Cache Groups",
					name: "TP_Cache4",
					type: "MID_LOC",
					validationMessage: "cachegroup was updated."
				}
			],
			remove: [
				{
					description: "delete a cachegroup",
					name: "TP_Cache2",
					validationMessage: "cachegroup was deleted."
				},
				{
					description: "delete a cachegroup",
					name: "TP_Cache4",
					validationMessage: "cachegroup was deleted."
				},
				{
					description: "delete a cachegroup",
					name: "TP_Cache5",
					validationMessage: "cachegroup was deleted."
				},
				{
					description: "delete a cachegroup",
					name: "TP_Cache6",
					validationMessage: "cachegroup was deleted."
				}
			]
		}
	]
};
