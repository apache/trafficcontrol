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
export const parameters = {
	cleanup: [
		{
			action: "DeleteParameters",
			route: "/parameters",
			method: "delete",
			data: [
				{
					route: "/parameters/",
					getRequest: [
						{
							route: "/parameters",
							queryKey: "name",
							queryValue: "TPParamtest2",
							replace: "route"
						}
					]
				}
			]
		}
	],
	setup: [
		{
			action: "CreateParameters",
			route: "/parameters",
			method: "post",
			data: [
				{
					name: "TPParamtest2",
					value: "quest",
					configFile: "records.config",
					secure: true
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
			check: [
				{
					description: "check CSV link from Parameter page",
					Name: "Export as CSV"
				}
			],
			toggle:[
				{
					description: "hide first table column",
					Name: "Config File"
				},
				{
					description: "redisplay first table column",
					Name: "Config File"
				}
			],
			add: [
				{
					description: "create a Parameters",
					Name: "TPParamtest1",
					ConfigFile: "test.config",
					Value: "90",
					Secure: "true",
					validationMessage: "Parameter created"
				}
			],
			update: [
				{
					description: "update parameter configfile",
					Name: "TPParamtest1",
					ConfigFile: "newtest.config",
					validationMessage: "Parameter updated"
				}
			],
			remove: [
				{
					description: "delete a Parameters",
					Name: "TPParamtest1",
					validationMessage: "parameter was deleted."
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
			check: [
				{
					description: "check CSV link from Parameter page",
					Name: "Export as CSV"
				}
			],
			toggle:[
				{
					description: "hide first table column",
					Name: "Config File"
				},
				{
					description: "redisplay first table column",
					Name: "Config File"
				}
			],
			add: [
				{
					description: "create a Parameters",
					Name: "TPParamtest1",
					ConfigFile: "test.config",
					Value: "90",
					Secure: "true",
					validationMessage: "missing required Permissions: PARAMETER:CREATE"
				}
			],
			update: [
				{
					description: "update parameter configfile",
					Name: "TPParamtest2",
					ConfigFile: "newtest.config",
					validationMessage: "missing required Permissions: PARAMETER:UPDATE"
				}
			],
			remove: [
				{
					description: "delete a Parameters",
					Name: "TPParamtest2",
					validationMessage: "missing required Permissions: PARAMETER:DELETE"
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
			check: [
				{
					description: "check CSV link from Parameter page",
					Name: "Export as CSV"
				}
			],
			toggle:[
				{
					description: "hide first table column",
					Name: "Config File"
				},
				{
					description: "redisplay first table column",
					Name: "Config File"
				}
			],
			add: [
				{
					description: "create a Parameters",
					Name: "TPParamtest1",
					ConfigFile: "test.config",
					Value: "90",
					Secure: "true",
					validationMessage: "Parameter created"
				}
			],
			update: [
				{
					description: "update parameter configfile",
					Name: "TPParamtest1",
					ConfigFile: "newtest.config",
					validationMessage: "Parameter updated"
				}
			],
			remove: [
				{
					description: "delete a Parameters",
					Name: "TPParamtest1",
					validationMessage: "parameter was deleted."
				}
			]
		}
	]
};
