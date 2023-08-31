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
export const coordinates = {
	cleanup: [
		{
			action: "DeleteCoordinates",
			route: "/coordinates",
			method: "delete",
			data: [
				{
					route: "/coordinates?id=",
					id: 0,
					getRequest: [
						{
							route: "/coordinates",
							queryKey: "name",
							queryValue: "TPCoordinates2",
							replace: "id"
						}
					]
				}
			]
		}
	],
	setup: [
		{
			action: "CreateCoordinates",
			route: "/coordinates",
			method: "post",
			data: [
				{
					name: "TPCoordinates2",
					latitude: 0,
					longitude: 0
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
					description: "Operator Role",
					username: "TPOperator",
					password: "pa$$word"
				}
			],
			add: [
				{
					description: "create a Coordinates",
					Name: "TPCoordinates1",
					Latitude: 0,
					Longitude: 0,
					validationMessage: "created"
				}
			],
			update: [
				{
					description: "update coordinates latitude",
					Name: "TPCoordinates1",
					Latitude: 1,
					validationMessage: "updated"
				}
			],
			remove: [
				{
					description: "delete a Coordinates",
					Name: "TPCoordinates1",
					validationMessage: "deleted"
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
					description: "create a Coordinates",
					Name: "TPCoordinates1",
					Latitude: 0,
					Longitude: 0,
					validationMessage: "missing required Permissions: COORDINATE:CREATE"
				}
			],
			update: [
				{
					description: "update coordinates latitude",
					Name: "TPCoordinates2",
					Latitude: 1,
					validationMessage: "missing required Permissions: COORDINATE:UPDATE"
				}
			],
			remove: [
				{
					description: "delete a Coordinates",
					Name: "TPCoordinates2",
					validationMessage: "missing required Permissions: COORDINATE:DELETE"
				}
			]
		}
	]
}
