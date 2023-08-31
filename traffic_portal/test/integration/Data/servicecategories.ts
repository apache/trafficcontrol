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
export const serviceCategories = {
	cleanup: [
		{
			action: "DeleteServiceCategories",
			route: "/service_categories",
			method: "delete",
			data: [
				{
					route: "/service_categories/TPTest2"
				}
			]
		}
	],
	setup: [
		{
			action: "CreateServiceCategories",
			route: "/service_categories",
			method: "post",
			data: [
				{
					name: "TPTest2"
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
					description: "create a Service Categories",
					Name: "TPServiceCategories1",
					validationMessage: "was created"
				}
			],
			update: [
				{
					description: "update service categories name",
					Name: "TPServiceCategories1",
					NewName: "TPSCNew1",
					validationMessage: "was updated"
				}
			],
			remove: [
				{
					description: "delete a service categories",
					Name: "TPSCNew1",
					validationMessage: "was deleted"
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
					description: "create a Service Categories",
					Name: "TPServiceCategories1",
					validationMessage: "missing required Permissions: SERVICE-CATEGORY:CREATE"
				}
			],
			update: [
				{
					description: "update service categories name",
					Name: "TPTest2",
					NewName: "TPSCNew1",
					validationMessage: "missing required Permissions: SERVICE-CATEGORY:UPDATE"
				}
			],
			remove: [
				{
					description: "delete a service categories",
					Name: "TPTest2",
					validationMessage: "missing required Permissions: SERVICE-CATEGORY:DELETE"
				}
			]
		}
	]
}
