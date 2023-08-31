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
export const physLocations = {
	cleanup: [
		{
			action: "DeleteServers",
			route: "/servers",
			method: "delete",
			data: [
				{
					route: "/servers/",
					getRequest: [
						{
							route: "/servers",
							queryKey: "hostName",
							queryValue: "PhysTest",
							replace: "route"
						}
					]
				}
			]
		},
		{
			action: "DeletePhysLocations",
			route: "/phys_locations",
			method: "delete",
			data: [
				{
					route: "/phys_locations/",
					getRequest: [
						{
							route: "/phys_locations",
							queryKey: "name",
							queryValue: "TPPhysLocation2",
							replace: "route"
						}
					]
				}
			]
		},
		{
			action: "DeleteRegions",
			route: "/regions",
			method: "delete",
			data: [
				{
					route: "/regions?name=PhysTest"
				},
				{
					route: "/regions?name=PhysTest2"
				}
			]
		},
		{
			action: "DeleteDivisions",
			route: "/divisions",
			method: "delete",
			data: [
				{
					route: "/divisions/",
					getRequest: [
						{
							route: "/divisions",
							queryKey: "name",
							queryValue: "PhysTest",
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
					name: "PhysTest"
				}
			]
		},
		{
			action: "CreateRegions",
			route: "/regions",
			method: "post",
			data: [
				{
					name: "PhysTest",
					division: "4",
					divisionName: "PhysTest",
					getRequest: [
						{
							route: "/divisions",
							queryKey: "name",
							queryValue: "PhysTest",
							replace: "division"
						}
					]
				},
				{
					name: "PhysTest2",
					division: "4",
					divisionName: "PhysTest",
					getRequest: [
						{
							route: "/divisions",
							queryKey: "name",
							queryValue: "PhysTest",
							replace: "division"
						}
					]
				}
			]
		},
		{
			action: "CreatePhysLocation",
			route: "/phys_locations",
			method: "post",
			data: [
				{
					address: "Buckingham Palace",
					city: "London",
					comments: "Buckingham Palace",
					email: "steve.kingstone@royal.gsx.gov.uk",
					name: "TPPhysLocation2",
					phone: "0-843-816-6276",
					poc: "Her Majesty The Queen Elizabeth Alexandra Mary Windsor II",
					regionId: 3,
					shortName: "tpphys2",
					state: "NA",
					zip: "99999",
					getRequest: [
						{
							route: "/regions",
							queryKey: "name",
							queryValue: "PhysTest",
							replace: "regionId"
						}
					]
				}
			]
		},
		{
			action: "CreateServers",
			route: "/servers",
			method: "post",
			data: [
				{
					cacheGroupID: 0,
					cdnID: 0,
					domainName: "test.net",
					hostName: "PhysTest",
					httpsPort: 443,
					iloIpAddress: "",
					iloIpGateway: "",
					iloIpNetmask: "",
					iloPassword: "",
					iloUsername: "",
					interfaces: [
						{
							ipAddresses: [
								{
									address: "::1",
									gateway: "::2",
									serviceAddress: true
								}
							],
							maxBandwidth: null,
							monitor: true,
							mtu: 1500,
							name: "eth0"
						}
					],
					interfaceMtu: 1500,
					interfaceName: "eth0",
					ip6Address: "::1",
					ip6Gateway: "::2",
					ipAddress: "0.0.0.1",
					ipGateway: "0.0.0.2",
					ipNetmask: "255.255.255.0",
					mgmtIpAddress: "",
					mgmtIpGateway: "",
					mgmtIpNetmask: "",
					offlineReason: "",
					physicalLocationID: 0,
					profiles: ["testProfile"],
					routerHostName: "",
					routerPortName: "",
					statusId: 3,
					tcpPort: 80,
					typeID: 12,
					updPending: false,
					getRequest: [
						{
							route: "/phys_locations",
							queryKey: "name",
							queryValue: "TPPhysLocation2",
							replace: "physicalLocationID"
						},
						{
							route: "/cdns",
							queryKey: "name",
							queryValue: "dummycdn",
							replace: "cdnID"
						},
						{
							route: "/cachegroups",
							queryKey: "name",
							queryValue: "testCG",
							replace: "cacheGroupID"
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
			check: [
				{
					description: "check CSV link from Physical Location page",
					Name: "Export as CSV"
				}
			],
			add: [
				{
					description: "create a PhysLocation",
					Name: "TPPhysLocation1",
					ShortName: "TPPhys1",
					Address: "address example",
					City: "city example",
					State: "NA",
					Zip: "11111",
					Poc: "test",
					Phone: "111-111-1111",
					Email: "emailtest@gtesting.com",
					Region: "PhysTest",
					Comments: "test",
					validationMessage: "Physical location created"
				}
			],
			update: [
				{
					description: "update physlocation region",
					Name: "TPPhysLocation1",
					Region: "PhysTest2",
					validationMessage: "Physical location updated"
				}
			],
			remove: [
				{
					description: "delete a PhysLocation",
					Name: "TPPhysLocation1",
					validationMessage: "Physical location deleted"
				},
				{
					description: "delete a PhysLocation that currently link with a Server",
					Name: "TPPhysLocation2",
					validationMessage: "can not delete a phys_location"
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
					description: "check CSV link from Physical Location page",
					Name: "Export as CSV"
				}
			],
			add: [
				{
					description: "create a PhysLocation",
					Name: "TPPhysLocation1",
					ShortName: "TPPhys1",
					Address: "address example",
					City: "city example",
					State: "NA",
					Zip: "11111",
					Poc: "test",
					Phone: "111-111-1111",
					Email: "emailtest@gtesting.com",
					Region: "PhysTest",
					Comments: "test",
					validationMessage: "missing required Permissions: PHYSICAL-LOCATION:CREATE"
				}
			],
			update: [
				{
					description: "update physlocation region",
					Name: "TPPhysLocation2",
					Region: "PhysTest2",
					validationMessage: "missing required Permissions: PHYSICAL-LOCATION:UPDATE"
				}
			],
			remove: [
				{
					description: "delete a PhysLocation",
					Name: "TPPhysLocation2",
					validationMessage: "missing required Permissions: PHYSICAL-LOCATION:DELETE"
				}
			]
		}
	]
};
