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
					Description: "create a EDGE_LOC cachegroup with FailOver CacheGroup Field",
					Name: "TP_Cache1",
					ShortName: "TPC1",
					Type: "EDGE_LOC",
					Latitude: "0",
					Longitude: "0",
					ParentCacheGroup: "infrastructure",
					SecondaryParentCG: "infrastructure",
					FailoverCG: "albany-ny-usa",
					validationMessage: "cache group was created."
				},
				{
					Description: "create multiple EDGE_LOC cachegroup",
					Name: "TP_Cache2",
					ShortName: "TPC2",
					Type: "EDGE_LOC",
					Latitude: "0",
					Longitude: "0",
					ParentCacheGroup: "infrastructure",
					SecondaryParentCG: "infrastructure",
					FailoverCG: "",
					validationMessage: "cache group was created."
				},
				{
					Description: "create a MID_LOC cachegroup",
					Name: "TP_Cache3",
					ShortName: "TPC3",
					Type: "MID_LOC",
					Latitude: "0",
					Longitude: "0",
					ParentCacheGroup: "infrastructure",
					SecondaryParentCG: "infrastructure",
					validationMessage: "cache group was created."
				}
			],
			update: [
				{
					Description: "add more Failover Cache Groups to EDGE_LOC Type cachegroup",
					Name: "TP_Cache1",
					Type: "EDGE_LOC",
					FailoverCG: "TP_Cache2",
					validationMessage: "cache group was updated"
				},
				{
					Description: "Validate cannot add cache group fallback if the cache group fall back is a different Type",
					Name:"TP_Cache1",
					Type:"EDGE_LOC",
					FailoverCG: "TP_Cache3"
				},
				{
					Description: "Validate cannot add an empty cache group fall back",
					Name:"TP_Cache1",
					Type:"EDGE_LOC",
					FailoverCG: " "
				},
				{
					Description: "change Type of the Cache Groups",
					Name: "TP_Cache1",
					Type: "MID_LOC",
					validationMessage: "cache group was updated"
				}
			],
			remove: [
				{
					Description: "delete a cachegroup",
					Name: "TP_Cache1",
					validationMessage: "cache group was deleted."
				},
				{
					Description: "delete a cachegroup",
					Name: "TP_Cache3",
					validationMessage: "cache group was deleted."
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
					Description: "create a CacheGroup",
					Name: "TP_Cache1",
					ShortName: "TPC1",
					Type: "EDGE_LOC",
					Latitude: "0",
					Longitude: "0",
					ParentCacheGroup: "infrastructure",
					SecondaryParentCG: "infrastructure",
					FailoverCG: "albany-ny-usa",
					validationMessage: "missing required Permissions: CACHE-GROUP:CREATE"
				}
			],
			update: [
				{
					Description: "update CacheGroup",
					Name: "TP_Cache2",
					Type: "MID_LOC",
					validationMessage: "missing required Permissions: CACHE-GROUP:UPDATE"
				}
			],
			remove: [
				{
					Description: "delete a cachegroup",
					Name: "TP_Cache2",
					validationMessage: "missing required Permissions: CACHE-GROUP:DELETE"
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
					Description: "create a EDGE_LOC cachegroup with FailOver CacheGroup Field",
					Name: "TP_Cache4",
					ShortName: "TPC4",
					Type: "EDGE_LOC",
					Latitude: "0",
					Longitude: "0",
					ParentCacheGroup: "infrastructure",
					SecondaryParentCG: "infrastructure",
					FailoverCG: "albany-ny-usa",
					validationMessage: "cache group was created."
				},
				{
					Description: "create multiple EDGE_LOC cachegroup",
					Name: "TP_Cache5",
					ShortName: "TPC5",
					Type: "EDGE_LOC",
					Latitude: "0",
					Longitude: "0",
					ParentCacheGroup: "infrastructure",
					SecondaryParentCG: "infrastructure",
					FailoverCG: "",
					validationMessage: "cache group was created."
				},
				{
					Description: "create a MID_LOC cachegroup",
					Name: "TP_Cache6",
					ShortName: "TPC6",
					Type: "MID_LOC",
					Latitude: "0",
					Longitude: "0",
					ParentCacheGroup: "infrastructure",
					SecondaryParentCG: "infrastructure",
					validationMessage: "cache group was created."
				}
			],
			update: [
				{
					Description: "add more Failover Cache Groups to EDGE_LOC Type cachegroup",
					Name: "TP_Cache4",
					Type: "EDGE_LOC",
					FailoverCG: "TP_Cache5",
					validationMessage: "cache group was updated"
				},
				{
					Description: "Validate cannot add cache group fallback if the cache group fall back is a different Type",
					Name:"TP_Cache4",
					Type:"EDGE_LOC",
					FailoverCG: "TP_Cache6"
				},
				{
					Description: "Validate cannot add an empty cache group fall back",
					Name:"TP_Cache4",
					Type:"EDGE_LOC",
					FailoverCG: " "
				},
				{
					Description: "change Type of the Cache Groups",
					Name: "TP_Cache4",
					Type: "MID_LOC",
					validationMessage: "cache group was updated"
				}
			],
			remove: [
				{
					Description: "delete a cachegroup",
					Name: "TP_Cache2",
					validationMessage: "cache group was deleted."
				},
				{
					Description: "delete a cachegroup",
					Name: "TP_Cache4",
					validationMessage: "cache group was deleted."
				},
				{
					Description: "delete a cachegroup",
					Name: "TP_Cache5",
					validationMessage: "cache group was deleted."
				},
				{
					Description: "delete a cachegroup",
					Name: "TP_Cache6",
					validationMessage: "cache group was deleted."
				}
			]
		}
	]
};
