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

export const ASNs = {
	cleanup: [
		{
			action: "DeleteCacheGroup",
			route: "/cachegroups",
			method: "delete",
			data: [
				{
					route: "/cachegroups/",
					getRequest: [
						{
							route: "/cachegroups",
							queryKey: "name",
							queryValue: "asntestcg1",
							replace: "route"
						}
					]
				},
				{
					route: "/cachegroups/",
					getRequest: [
						{
							route: "/cachegroups",
							queryKey: "name",
							queryValue: "asntestcg2",
							replace: "route"
						}
					]
				},
				{
					route: "/cachegroups/",
					getRequest: [
						{
							route: "/cachegroups",
							queryKey: "name",
							queryValue: "asntestcg3",
							replace: "route"
						}
					]
				},
				{
					route: "/cachegroups/",
					getRequest: [
						{
							route: "/cachegroups",
							queryKey: "name",
							queryValue: "asntestcg4",
							replace: "route"
						}
					]
				}
			]
		}
	],
	setup: [
		{
			action: "CreateCacheGroups",
			route: "/cachegroups",
			method: "post",
			data: [
				{
					name: "asntestcg1",
					shortName: "asntcg1",
					latitude: 0,
					longitude: 0,
					fallbackToClosest: true,
					localizationMethods: [
						"DEEP_CZ",
						"CZ",
						"GEO"
					],
					typeId: 23
				},
				{
					name: "asntestcg2",
					shortName: "asntcg2",
					latitude: 0,
					longitude: 0,
					fallbackToClosest: true,
					localizationMethods: [
						"DEEP_CZ",
						"CZ",
						"GEO"
					],
					typeId: 23
				},
				{
					name: "asntestcg3",
					shortName: "asntcg3",
					latitude: 0,
					longitude: 0,
					fallbackToClosest: true,
					localizationMethods: [
						"DEEP_CZ",
						"CZ",
						"GEO"
					],
					typeId: 23
				},
				{
					name: "asntestcg4",
					shortName: "asntcg4",
					latitude: 0,
					longitude: 0,
					fallbackToClosest: true,
					localizationMethods: [
						"DEEP_CZ",
						"CZ",
						"GEO"
					],
					typeId: 23
				}
			]
		}
	],
	tests: [
		{
			logins: [
				{
					description: "Admin Role",
					username: "TPAdmin",
					password: "pa$$word"
				}
			],
			add: [
				{
					description: "create an ASN",
					ASNs: "1111",
					CacheGroup: "asntestcg1",
					validationMessage: "ASN created"
				},
				{
					description: "create multiple ASN",
					ASNs: "2222",
					CacheGroup: "asntestcg3",
					validationMessage: "ASN created"
				},
				{
					description: "create an ASN with existence number",
					ASNs: "2222",
					CacheGroup: "asntestcg1",
					validationMessage: "already exists"
				},
				{
					description: "create an ASN with existence number different cachegroup",
					ASNs: "2222",
					CacheGroup: "asntestcg2",
					validationMessage: "already exists"
				}
			],
			update: [
				{
					description: "update an ASN to have unique number",
					ASNs: "1111",
					NewASNs: "3333",
					validationMessage: "ASN updated"
				},
				{
					description: "update an ASN to have existence number",
					ASNs: "3333",
					NewASNs: "2222",
					validationMessage: "already exists"
				},
				{
					description: "update cachegroup of an ASN",
					ASNs: "3333",
					CacheGroup: "asntestcg2",
					validationMessage: "ASN updated"
				}
			],
			remove: [
				{
					description: "delete an ASN",
					ASNs: "3333",
					validationMessage: "asn was deleted."
				}
			]
		},
		{
			logins: [
				{
					description: "Read Only Role",
					username: "TPReadOnly",
					password: "pa$$word"
				}
			],
			add: [
				{
					description: "create an ASN",
					ASNs: "1111",
					CacheGroup: "asntestcg1",
					validationMessage: "missing required Permissions: ASN:CREATE, CACHE-GROUP:UPDATE"
				}
			],
			update: [
				{
					description: "update cachegroup of an ASN",
					ASNs: "2222",
					CacheGroup: "asntestcg2",
					validationMessage: "missing required Permissions: ASN:UPDATE, CACHE-GROUP:UPDATE"
				}
			],
			remove: [
				{
					description: "delete an ASN",
					ASNs: "2222",
					validationMessage: "missing required Permissions: ASN:DELETE, CACHE-GROUP:UPDATE"
				}
			]
		},
		{
			logins: [
				{
					description: "Operation Role",
					username: "TPOperator",
					password: "pa$$word"
				}
			],
			add: [
				{
					description: "create an ASN",
					ASNs: "4444",
					CacheGroup: "asntestcg1",
					validationMessage: "ASN created"
				},
				{
					description: "create multiple ASN",
					ASNs: "5555",
					CacheGroup: "asntestcg4",
					validationMessage: "ASN created"
				},
				{
					description: "create an ASN with existence number",
					ASNs: "4444",
					CacheGroup: "asntestcg1",
					validationMessage: "already exists"
				},
				{
					description: "create an ASN with existence number different cachegroup",
					ASNs: "4444",
					CacheGroup: "asntestcg2",
					validationMessage: "already exists"
				}
			],
			update: [
				{
					description: "update an ASN to have unique number",
					ASNs: "4444",
					NewASNs: "6666",
					validationMessage: "ASN updated"
				},
				{
					description: "update an ASN to have existence number",
					ASNs: "6666",
					NewASNs: "5555",
					validationMessage: "already exists"
				},
				{
					description: "update cachegroup of an ASN",
					ASNs: "6666",
					CacheGroup: "asntestcg2",
					validationMessage: "ASN updated"
				}
			],
			remove: [
				{
					description: "delete an ASN",
					ASNs: "6666",
					validationMessage: "asn was deleted."
				},
				{
					description: "delete an ASN",
					ASNs: "2222",
					validationMessage: "asn was deleted."
				},
				{
					description: "delete an ASN",
					ASNs: "5555",
					validationMessage: "asn was deleted."
				}
			]
		}
	]
}
