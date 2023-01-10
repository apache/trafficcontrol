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


interface Removal {
	description: string;
	name: string;
	validationMessage: string;
}

interface InvalidRemoval extends Removal {
	invalid: true;
	invalidName: string;
}

interface ValidRemoval extends Removal {
	invalid: false;
}

type PossiblyInvalidRemoval = InvalidRemoval | ValidRemoval;

export const serverCapabilities = {
	tests: [
		{
			testName: "Admin Role",
			logins: [
				{
					username: "TPAdmin",
					password: "pa$$word"
				}
			],
			check: [
				{
					description: "check CSV link from Server Capabilities page",
					Name: "Export as CSV"
				}
			],
			add: [
				{
					description: "can create a server capability",
					name: "TP_SC",
					capabilityDescription: "Server Capability",
					validationMessage: "server capability was created."
				},
				{
					description: "can create multiple server capabilities",
					name: "TP_SC_2",
					capabilityDescription: "Server Capability 2",
					validationMessage: "server capability was created."
				},
				{
					description: "can create multiple server capabilities",
					name: "TP_SC_3",
					capabilityDescription: "Server Capability 3",
					validationMessage: "server capability was created."
				},
				{
					description: "can handle creating existing server capability",
					name: "TP_SC_2",
					capabilityDescription: "Server Capability 2",
					validationMessage: "server_capability name 'TP_SC_2' already exists."
				},
				{
					description: "can handle invalid period in server capability",
					name: "TP.AUTOMATED",
					capabilityDescription: "Server Capability Automated",
					validationMessage: "Must be alphamumeric with no spaces. Dashes and underscores also allowed."
				},
				{
					description: "can handle invalid space in server capability",
					name: "TP AUTOMATED",
					capabilityDescription: "Server Capability Automated",
					validationMessage: "Must be alphamumeric with no spaces. Dashes and underscores also allowed."
				},
				{
					description: "can handle invalid character in server capability",
					name: "TP#AUTOMATED",
					capabilityDescription: "Server Capability Automated",
					validationMessage: "Must be alphamumeric with no spaces. Dashes and underscores also allowed."
				},
				{
					description: "can handle an empty text field",
					name: "",
					capabilityDescription: "",
					validationMessage: "Required"
				}
			],
			remove: [
				{
					description: "can delete a server capability",
					invalid: false,
					name: "TP_SC",
					validationMessage: "server capability was deleted."
				},
				{
					description: "can handle an invalid delete entry",
					invalid: true,
					name: "TP_SC_2",
					invalidName: "TP_AUTOMATED"
				},
				{
					description: "can delete multiple server capabilities",
					invalid: false,
					name: "TP_SC_2",
					validationMessage: "server capability was deleted."
				}
			] as Array<PossiblyInvalidRemoval>
		},
		{
			testName: "ReadOnly Role",
			logins: [
				{
					username: "TPReadOnly",
					password: "pa$$word"
				}
			],
			check: [
				{
					description: "check CSV link from Server Capabilities page",
					Name: "Export as CSV"
				}
			],
			add: [
				{
					description: "can handle readonly role creating a server capability",
					name: "TP_SC",
					capabilityDescription: "Server Capability",
					validationMessage: "missing required Permissions: SERVER-CAPABILITY:CREATE"
				}
			],
			remove: [
				{
					description: "can handle readonly role deleting a server capability",
					invalid: false,
					name: "TP_SC_3",
					validationMessage: "missing required Permissions: SERVER-CAPABILITY:DELETE"
				}
			] as Array<PossiblyInvalidRemoval>
		},
		{
			testName: "Operation Role",
			logins: [
				{
					username: "TPOperator",
					password: "pa$$word"
				}
			],
			check: [
				{
					description: "check CSV link from Server Capabilities page",
					Name: "Export as CSV"
				}
			],
			add: [
				{
					description: "can create a server capability",
					name: "TP_SC",
					capabilityDescription: "Server Capability",
					validationMessage: "server capability was created."
				},
				{
					description: "can create multiple server capabilities",
					name: "TP_SC_2",
					capabilityDescription: "Server Capability 2",
					validationMessage: "server capability was created."
				},
				{
					description: "can handle creating existing server capability",
					name: "TP_SC_2",
					capabilityDescription: "Server Capability 2",
					validationMessage: "server_capability name 'TP_SC_2' already exists."
				},
				{
					description: "can handle invalid period in server capability",
					name: "TP.AUTOMATED",
					capabilityDescription: "Server Capability Automated",
					validationMessage: "Must be alphamumeric with no spaces. Dashes and underscores also allowed."
				},
				{
					description: "can handle invalid space in server capability",
					name: "TP AUTOMATED",
					capabilityDescription: "Server Capability Automated",
					validationMessage: "Must be alphamumeric with no spaces. Dashes and underscores also allowed."
				},
				{
					description: "can handle invalid character in server capability",
					name: "TP#AUTOMATED",
					capabilityDescription: "Server Capability Automated",
					validationMessage: "Must be alphamumeric with no spaces. Dashes and underscores also allowed."
				},
				{
					description: "can handle an empty text field",
					name: "",
					capabilityDescription: "",
					validationMessage: "Required"
				}
			],
			remove: [
				{
					description: "can delete a server capability",
					name: "TP_SC",
					invalid: false,
					validationMessage: "server capability was deleted."
				},
				{
					description: "can handle an invalid delete entry",
					name: "TP_SC_2",
					invalid: true,
					invalidName: "TP_AUTOMATED"
				},
				{
					description: "can delete multiple server capabilities",
					name: "TP_SC_2",
					invalid: false,
					validationMessage: "server capability was deleted."
				},
				{
					description: "can delete multiple server capabilities",
					name: "TP_SC_3",
					invalid: false,
					validationMessage: "server capability was deleted."
				}
			] as Array<PossiblyInvalidRemoval>
		}
	]
}
