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
export const prerequisites = [
	{
		action: "CreateTenants",
		route: "/tenants",
		method: "post",
		data: [
			{
				active: true,
				name: "tenantParent",
				parentId: 1
			},
			{
				active: true,
				name: "tenantSame",
				parentId: 1,
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantParent",
						replace: "parentId"
					}
				]
			},
			{
				active: true,
				name: "tenantChild",
				parentId: 1,
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantSame",
						replace: "parentId"
					}
				]
			},
			{
				active: true,
				name: "tenantDifferent",
				parentId: 1
			}
		]
	},
	{
		action: "CreateUsers",
		route: "/users",
		method: "post",
		data: [
			{
				fullName: "TPAdmin",
				username: "TPAdmin",
				email: "@test.com",
				role: "admin",
				tenantId: 1,
				localPasswd: "pa$$word",
				confirmLocalPasswd: "pa$$word",
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantSame",
						replace: "tenantId"
					}
				]
			},
			{
				fullName: "TPOperator",
				username: "TPOperator",
				email: "@test.com",
				role: "operations",
				tenantId: 1,
				localPasswd: "pa$$word",
				confirmLocalPasswd: "pa$$word",
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantSame",
						replace: "tenantId"
					}
				]
			},
			{
				fullName: "TPReadOnly",
				username: "TPReadOnly",
				email: "@test.com",
				role: "read-only",
				tenantId: 1,
				localPasswd: "pa$$word",
				confirmLocalPasswd: "pa$$word",
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantSame",
						replace: "tenantId"
					}
				]
			},
			{
				fullName: "TPAdminDiff",
				username: "TPAdminDiff",
				email: "@test.com",
				role: "admin",
				tenantId: 1,
				localPasswd: "pa$$word",
				confirmLocalPasswd: "pa$$word",
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantDifferent",
						replace: "tenantId"
					}
				]
			},
			{
				fullName: "TPOperatorDiff",
				username: "TPOperatorDiff",
				email: "@test.com",
				role: "operations",
				tenantId: 1,
				localPasswd: "pa$$word",
				confirmLocalPasswd: "pa$$word",
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantDifferent",
						replace: "tenantId"
					}
				]
			},
			{
				fullName: "TPReadOnlyDiff",
				username: "TPReadOnlyDiff",
				email: "@test.com",
				role: "read-only",
				tenantId: 1,
				localPasswd: "pa$$word",
				confirmLocalPasswd: "pa$$word",
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantDifferent",
						replace: "tenantId"
					}
				]
			},
			{
				fullName: "TPAdminParent",
				username: "TPAdminParent",
				email: "@test.com",
				role: "admin",
				tenantId: 1,
				localPasswd: "pa$$word",
				confirmLocalPasswd: "pa$$word",
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantParent",
						replace: "tenantId"
					}
				]
			},
			{
				fullName: "TPOperatorParent",
				username: "TPOperatorParent",
				email: "@test.com",
				role: "operations",
				tenantId: 1,
				localPasswd: "pa$$word",
				confirmLocalPasswd: "pa$$word",
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantParent",
						replace: "tenantId"
					}
				]
			},
			{
				fullName: "TPReadOnlyParent",
				username: "TPReadOnlyParent",
				email: "@test.com",
				role: "read-only",
				tenantId: 1,
				localPasswd: "pa$$word",
				confirmLocalPasswd: "pa$$word",
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantParent",
						replace: "tenantId"
					}
				]
			},
			{
				fullName: "TPAdminChild",
				username: "TPAdminChild",
				email: "@test.com",
				role: "admin",
				tenantId: 1,
				localPasswd: "pa$$word",
				confirmLocalPasswd: "pa$$word",
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantChild",
						replace: "tenantId"
					}
				]
			},
			{
				fullName: "TPOperatorChild",
				username: "TPOperatorChild",
				email: "@test.com",
				role: "operations",
				tenantId: 1,
				localPasswd: "pa$$word",
				confirmLocalPasswd: "pa$$word",
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantChild",
						replace: "tenantId"
					}
				]
			},
			{
				fullName: "TPReadOnlyChild",
				username: "TPReadOnlyChild",
				email: "@test.com",
				role: "read-only",
				tenantId: 1,
				localPasswd: "pa$$word",
				confirmLocalPasswd: "pa$$word",
				getRequest: [
					{
						route: "/tenants",
						queryKey: "name",
						queryValue: "tenantChild",
						replace: "tenantId"
					}
				]
			}
		]
	},
	{
		action: "CreateCDN",
		route: "/cdns",
		method: "post",
		data: [
			{
				name: "dummycdn",
				domainName: "cdnp3",
				dnssecEnabled: false
			}
		]
	},
	{
		action: "CreateCacheGroups",
		route: "/cachegroups",
		method: "post",
		data: [
			{
				name: "testCG",
				shortName: "tCG",
				latitude: 0,
				longitude: 0,
				fallbackToClosest: true,
				localizationmethods: [
					"DEEP_CZ",
					"CZ",
					"GEO"
				],
				typeId: 23
			}
		]
	},
	{
		action: "CreateProfile",
		route: "/profiles",
		method: "post",
		data: [
			{
				name: "testProfile",
				description: "A test profile for API examples",
				cdn: 1,
				type: "ATS_PROFILE",
				routingDisabled: false,
				getRequest: [
					{
						route: "/cdns",
						queryKey: "name",
						queryValue: "dummycdn",
						replace: "cdn"
					}
				]
			}
		]
	}
]
