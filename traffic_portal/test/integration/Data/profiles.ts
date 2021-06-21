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
export const profiles = {
	cleanup: [
		{
			action: "DeleteProfile",
			route: "/profiles",
			method: "delete",
			data: [
				{
					route: "/profiles/",
					getRequest: [
						{
							route: "/profiles",
							queryKey: "name",
							queryValue: "TPProfiles2",
							replace: "route"
						}
					]
				}
			]
		},
		{
			action: "DeleteCDN",
			route: "/cdns",
			method: "delete",
			data: [
				{
					route: "/cdns/",
					getRequest: [
						{
							route: "/cdns",
							queryKey: "name",
							queryValue: "cdnprofile1",
							replace: "route"
						}
					]

				}
			]
		}
	],
	setup: [
		{
			action: "CreateCDN",
			route: "/cdns",
			method: "post",
			data: [
				{
					name: "cdnprofile1",
					domainName: "cdnp1",
					dnssecEnabled: false
				}

			]
		},
		{
			action: "CreateProfile",
			route: "/profiles",
			method: "post",
			data: [
				{
					name: "TPProfiles2",
					description: "A test profile for API examples",
					cdn: 2,
					type: "UNK_PROFILE",
					routingDisabled: true,
					getRequest: [
						{
							route: "/cdns",
							queryKey: "name",
							queryValue: "cdnprofile1",
							replace: "cdn"
						}
					]
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
				},
				{
					description: "Operation Role",
					username: "TPOperator",
					password: "pa$$word"
				}
			],
			add: [
				{
					description: "create a Profiles",
					Name: "TPProfiles1",
					CDN: "CDN-in-a-Box",
					Type: "ATS_PROFILE",
					RoutingDisable: "true",
					Description: "testing",
					validationMessage: "Profile created"
				}
			],
			update: [
				{
					description: "update profile type",
					Name: "TPProfiles1",
					Type: "RIAK_PROFILE",
					validationMessage: "Profile updated"
				}
			],
			remove: [
				{
					description: "delete a Profile",
					Name: "TPProfiles1",
					validationMessage: "profile was deleted."
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
					description: "create a Profiles",
					Name: "TPProfiles1",
					CDN: "CDN-in-a-Box",
					Type: "ATS_PROFILE",
					RoutingDisable: "true",
					Description: "testing",
					validationMessage: "Forbidden."
				}
			],
			update: [
				{
					description: "update profile type",
					Name: "TPProfiles2",
					Type: "RIAK_PROFILE",
					validationMessage: "Forbidden."
				}
			],
			remove: [
				{
					description: "delete a Profile",
					Name: "TPProfiles2",
					validationMessage: "Forbidden."
				}
			]
		}
	]
};
