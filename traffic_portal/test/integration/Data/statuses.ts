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
export const statuses = {
	cleanup: [
		{
			action: "DeleteStatuses",
			route: "/statuses",
			method : "delete",
			data: [
				{
					route: "/statuses/",
					getRequest: [
						{
							route: "/statuses",
							queryKey: "name",
							queryValue: "TPStatus2",
							replace: "route"
						}
					]
				}
			]
		}
	],
	setup: [
		{
			action: "CreateStatuses",
			route: "/statuses",
			method: "post",
			data: [
				{
					name: "TPStatus2",
					description: "For readonly"
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
					description: "create a Statuses",
					Name: "TPStatus1",
					DescriptionData: "test",
					validationMessage: "Status created"
				},
				{
					description: "create a Statuses with same name",
					Name: "TPStatus1",
					DescriptionData: "test",
					validationMessage: "already exists."
				}
			],
			update: [
				{
					description: "update Status description",
					Name: "TPStatus1",
					DescriptionData: "update",
					validationMessage: "Status updated"
				}
			],
			remove: [
				{
					description: "delete a Status",
					Name: "TPStatus1",
					validationMessage: "Status deleted"
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
					description: "create a Statuses",
					Name: "TPStatus1",
					DescriptionData: "test",
					validationMessage: "missing required Permissions: STATUS:CREATE"
				}
			],
			update: [
				{
					description: "update Status description",
					Name: "TPStatus2",
					DescriptionData: "update",
					validationMessage: "missing required Permissions: STATUS:UPDATE"
				}
			],
			remove: [
				{
					description: "delete a Status",
					Name: "TPStatus2",
					validationMessage: "missing required Permissions: STATUS:DELETE"
				}
			]
		}
	]
}
