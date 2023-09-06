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
export const divisions = {
	cleanup: [
		{
			action: "DeleteDivisions",
			route : "/divisions",
			method : "delete",
			data: [
				{
					route: "/divisions/",
					getRequest: [
						{
							route: "/divisions",
							queryKey: "name",
							queryValue: "TPDivision2",
							replace: "route"
						}
					]
				}
			]
		}
	],
	setup: [
		{
			action: "CreateDivisions",
			route: "/divisions",
			method: "post",
			data: [
				{
					name: "TPDivision2"
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
			check: [
				{
					description: "check CSV link from Division page",
					Name: "Export as CSV"
				}
			],
			add: [
				{
					description: "create a Divisions",
					Name: "TPDivision1",
					validationMessage: "division was created."
				}
			],
			update: [
				{
					description: "update Division's name",
					Name: "TPDivision1",
					NewName: "NewDivision1",
					validationMessage: "division was updated"
				}
			],
			remove: [
				{
					description: "delete Division",
					Name: "NewDivision1",
					validationMessage: "division was deleted."
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
					description: "check CSV link from Division page",
					Name: "Export as CSV"
				}
			],
			add: [
				{
					description: "create a Divisions",
					Name: "TPDivision1",
					validationMessage: "missing required Permissions: DIVISION:CREATE"
				}
			],
			update: [
				{
					description: "update Division's name",
					Name: "TPDivision2",
					NewName: "NewDivision2",
					validationMessage: "missing required Permissions: DIVISION:UPDATE"
				}
			],
			remove: [
				{
					description: "delete Division",
					Name: "TPDivision2",
					validationMessage: "missing required Permissions: DIVISION:DELETE"
				}
			]
		}
	]
};
