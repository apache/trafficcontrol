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
export const cdns = {
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
					description: "create a CDN",
					Name: "TPCDN1",
					Domain: "tptest1",
					DNSSEC: "true",
					validationMessage: "cdn was created."
				},
				{
					description: "create multiple CDN with false DNSSEC",
					Name: "TPCDN2",
					Domain: "tptest2",
					DNSSEC: "false",
					validationMessage: "cdn was created."
				},
				{
					description: "Validate cannot create CDN with invalid characters",
					Name: "TP_CDN1",
					Domain: "tptest1",
					DNSSEC: "true",
					validationMessage: "'name' invalid characters found - Use alphanumeric . or - ."
				}
			],
			update: [
				{
					description: "perform snapshot",
					Name: "TPCDN1",
					NewName: "",
					validationMessage: "Snapshot performed"
				},
				{
					description: "queue CDN updates",
					Name: "TPCDN2",
					NewName: "",
					validationMessage: "Queued CDN server updates"
				},
				{
					description: "clear CDN updates",
					Name: "TPCDN2",
					NewName: "",
					validationMessage: "Cleared CDN server updates"
				},
				{
					description: "update cdn ttl override",
					Name: "TPCDN2",
					NewName: "",
					validationMessage: "cdn was updated."
				},
				{
					description: "update cdn name",
					Name: "TPCDN1",
					NewName: "TPNewCDN1",
					validationMessage: "cdn was updated."
				},
			],
			remove: [
				{
					description: "delete CDN",
					Name: "TPNewCDN1",
					validationMessage: "cdn was deleted."
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
					description: "create a CDN",
					Name: "TPCDN1",
					Domain: "tptest1",
					DNSSEC: "true",
					validationMessage: "missing required Permissions: CDN:CREATE"
				}
			],
			update: [
				{
					description: "perform snapshot",
					Name: "TPCDN2",
					NewName: "",
					validationMessage: "missing required Permissions: CDN-SNAPSHOT:CREATE"
				},
				{
					description: "queue CDN updates",
					Name: "TPCDN2",
					NewName: "",
					validationMessage: "missing required Permissions: SERVER:QUEUE"
				},
				{
					description: "clear CDN updates",
					Name: "TPCDN2",
					NewName: "",
					validationMessage: "missing required Permissions: SERVER:QUEUE"
				},
				{
					description: "update cdn ttl override",
					Name: "TPCDN2",
					NewName: "",
					validationMessage: "missing required Permissions: CDN:UPDATE"
				},
				{
					description: "update cdn name",
					Name: "TPCDN2",
					NewName: "TPNewCDN2",
					validationMessage: "missing required Permissions: CDN:UPDATE"
				},
			],
			remove: [
				{
					description: "delete CDN",
					Name: "TPCDN2",
					validationMessage: "missing required Permissions: CDN:DELETE"
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
					description: "create a CDN",
					Name: "TPCDN3",
					Domain: "tptest3",
					DNSSEC: "true",
					validationMessage: "cdn was created."
				},
				{
					description: "create multiple CDN with false DNSSEC",
					Name: "TPCDN4",
					Domain: "tptest4",
					DNSSEC: "false",
					validationMessage: "cdn was created."
				},
				{
					description: "Validate cannot create CDN with invalid characters",
					Name: "TP_CDN5",
					Domain: "tptest5",
					DNSSEC: "true",
					validationMessage: "'name' invalid characters found - Use alphanumeric . or - ."
				}
			],
			update: [
				{
					description: "perform snapshot",
					Name: "TPCDN3",
					NewName: "",
					validationMessage: "Snapshot performed"
				},
				{
					description: "queue CDN updates",
					Name: "TPCDN4",
					NewName: "",
					validationMessage: "Queued CDN server updates"
				},
				{
					description: "clear CDN updates",
					Name: "TPCDN4",
					NewName: "",
					validationMessage: "Cleared CDN server updates"
				},
				{
					description: "update cdn ttl override",
					Name: "TPCDN4",
					NewName: "",
					validationMessage: "cdn was updated."
				},
				{
					description: "update cdn name",
					Name: "TPCDN3",
					NewName: "TPNewCDN3",
					validationMessage: "cdn was updated."
				},
			],
			remove: [
				{
					description: "delete CDN",
					Name: "TPCDN2",
					validationMessage: "cdn was deleted."
				},
				{
					description: "delete CDN",
					Name: "TPNewCDN3",
					validationMessage: "cdn was deleted."
				},
				{
					description: "delete CDN",
					Name: "TPCDN4",
					validationMessage: "cdn was deleted."
				}
			]
		}
	]
};
