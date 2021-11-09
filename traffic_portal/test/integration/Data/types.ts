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
export const types = {
	cleanup: [
		{
			action: "DeleteTypes",
			route: "/types",
			method: "delete",
			data: [
				{
					route: "/types/",
					getRequest: [
						{
							route: "/types",
							queryKey: "name",
							queryValue: "TPType3",
							replace: "route"
						}
					]
				}
			]
		}
	],
	setup: [
		{
			action: "CreateTypes",
			route: "/types",
			method: "post",
			data: [
				{
					name: "TPType3",
					description: "For readonly",
					useInTable: "server"
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
			toggle:[
				{
					description: "hide first table column",
					Name: "description"
				},
				{
					description: "redisplay first table column",
					Name: "description"
				}
			],
			check: [
				{
					description: "check CSV link from Type page",
					Name: "Export as CSV"
				}
			],
			add: [
				{
					description: "create a Type",
					Name: "TPType1",
					DescriptionData: "This is a test",
					validationMessage: "Type created"
				},
				{
					description: "create a Type without description",
					Name: "TPType2",
					DescriptionData: "",
					validationMessage: "'description' cannot be blank"
				}
			],
			update: [
				{
					description: "update description type",
					Name: "TPType1",
					DescriptionData: "Change description",
					validationMessage: "Type updated"
				}
			],
			remove: [
				{
					description: "delete a type",
					Name: "TPType1",
					validationMessage: "Type deleted"
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
			toggle:[
				{
					description: "hide first table column",
					Name: "description"
				},
				{
					description: "redisplay first table column",
					Name: "description"
				}
			],
			check: [
				{
					description: "check CSV link from Type page",
					Name: "Export as CSV"
				}
			],
			add: [
				{
					description: "create a Type",
					Name: "TPType1",
					DescriptionData: "This is a test",
					validationMessage: "missing required Permissions: TYPE:CREATE"
				}
			],
			update: [
				{
					description: "update description type",
					Name: "TPType3",
					DescriptionData: "Change description",
					validationMessage: "missing required Permissions: TYPE:UPDATE"
				}
			],
			remove: [
				{
					description: "delete a type",
					Name: "TPType3",
					validationMessage: "missing required Permissions: TYPE:DELETE"
				}
			]
		}
	]
}
