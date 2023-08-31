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
export const serverServerCapabilities = {
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
							queryValue: "testserver1",
							replace: "route"
						}
					]
				},
				{
					route: "/servers/",
					getRequest: [
						{
							route: "/servers",
							queryKey: "hostName",
							queryValue: "testserver2",
							replace: "route"
						}
					]
				},
				{
					route: "/servers/",
					getRequest: [
						{
							route: "/servers",
							queryKey: "hostName",
							queryValue: "testserver3",
							replace: "route"
						}
					]
				},
				{
					route: "/servers/",
					getRequest: [
						{
							route: "/servers",
							queryKey: "hostName",
							queryValue: "testserver4",
							replace: "route"
						}
					]
				},
				{
					route: "/servers/",
					getRequest: [
						{
							route: "/servers",
							queryKey: "hostName",
							queryValue: "testserver5",
							replace: "route"
						}
					]
				},
				{
					route: "/servers/",
					getRequest: [
						{
							route: "/servers",
							queryKey: "hostName",
							queryValue: "testserver6",
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
		},
		{
			action: "DeleteServerCapabilities",
			route: "/server_capabilities",
			method: "delete",
			data: [
				{
					route: "/server_capabilities?name=servercap1"
				},
				{
					route: "/server_capabilities?name=servercap2"
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
					hostName: "testserver1",
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
				},
				{
					cacheGroupID: 8,
					cdnID: 2,
					domainName: "test.net",
					hostName: "testserver2",
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
					physicalLocationID: 2,
					profiles: ["testProfile"],
					routerHostName: "",
					routerPortName: "",
					statusId: 3,
					tcpPort: 80,
					typeID: 12,
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
				},
				{
					cacheGroupID: 8,
					cdnID: 2,
					domainName: "test.net",
					hostName: "testserver3",
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
					physicalLocationID: 2,
					profiles: ["testProfile"],
					routerHostName: "",
					routerPortName: "",
					statusId: 3,
					tcpPort: 80,
					typeID: 13,
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
				},
				{
					cacheGroupID: 8,
					cdnID: 2,
					domainName: "test.net",
					hostName: "testserver4",
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
					physicalLocationID: 2,
					profiles: ["testProfile"],
					routerHostName: "",
					routerPortName: "",
					statusId: 3,
					tcpPort: 80,
					typeID: 12,
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
				},
				{
					cacheGroupID: 8,
					cdnID: 2,
					domainName: "test.net",
					hostName: "testserver5",
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
					physicalLocationID: 2,
					profiles: ["testProfile"],
					routerHostName: "",
					routerPortName: "",
					statusId: 3,
					tcpPort: 80,
					typeID: 12,
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
				},
				{
					cacheGroupID: 8,
					cdnID: 2,
					domainName: "test.net",
					hostName: "testserver6",
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
					physicalLocationID: 2,
					profiles: ["testProfile"],
					routerHostName: "",
					routerPortName: "",
					statusId: 3,
					tcpPort: 80,
					typeID: 13,
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
		},
		{
			action: "CreateServerCapabilities",
			route: "/server_capabilities",
			method: "post",
			data: [
				{
					name: "servercap1"
				},
				{
					name: "servercap2"
				},
				{
					name: "servercap3"
				},
				{
					name: "servercap4"
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
			link: [
				{
					description: "link server capability to server",
					Server: "testserver1",
					ServerCapability: "servercap1",
					validationMessage: "server server_capability was created."
				},
				{
					description: "link multiple server capabilities to server",
					Server: "testserver1",
					ServerCapability: "servercap2",
					validationMessage:"server server_capability was created."
				},
				{
					description: "link server capability to multiple servers",
					Server: "testserver2",
					ServerCapability: "servercap2",
					validationMessage: "server server_capability was created."
				},
				{
					description: "Validate cannot link server capabilities to server other than MID or EDGE",
					Server: "testserver3",
					ServerCapability: "servercap2"
				},
				{
					description: "link same server capability to servers",
					Server: "testserver2",
					ServerCapability: "servercap2",
					validationMessage: "already exists."
				}
			],
			remove: [
				{
					description: "remove server capability from server",
					Server: "testserver1",
					ServerCapability: "servercap1",
					validationMessage: "server server_capability was deleted."
				}
			],
			deleteServerCapability: [
				{
					description: "delete server capability linked with one or more servers",
					ServerCapability: "servercap2",
					validationMessage: "can not delete a server capability with 2 assigned servers"
				},
				{
					description: "delete server capabilities that is not link to any server",
					ServerCapability: "servercap3",
					validationMessage: "server capability was deleted."
				}
			]
		},
		{
			logins: [
				{
					description: "ReadOnly Role",
					username: "TPReadOnly",
					password: "pa$$word"
				}
			],
			link: [
				{
					description: "link server capability to server",
					Server: "testserver1",
					ServerCapability: "servercap1",
					validationMessage: "missing required Permissions: SERVER:UPDATE"
				}
			],
			remove: [
				{
					description: "remove server capability from server",
					Server: "testserver1",
					ServerCapability: "servercap2",
					validationMessage: "missing required Permissions: SERVER:UPDATE"
				}
			],
			deleteServerCapability: [
				{
					description: "delete server capability linked with one or more servers",
					ServerCapability: "servercap2",
					validationMessage: "missing required Permissions: SERVER-CAPABILITY:DELETE"
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
			link: [
				{
					description: "link server capability to server",
					Server: "testserver4",
					ServerCapability: "servercap1",
					validationMessage: "server server_capability was created."
				},
				{
					description: "link multiple server capabilities to server",
					Server: "testserver4",
					ServerCapability: "servercap2",
					validationMessage:"server server_capability was created."
				},
				{
					description: "link server capability to multiple servers",
					Server: "testserver5",
					ServerCapability: "servercap2",
					validationMessage: "server server_capability was created."
				},
				{
					description: "Validate cannot link server capabilities to server other than MID or EDGE",
					Server: "testserver6",
					ServerCapability: "servercap2"
				},
				{
					description: "link same server capability to servers",
					Server: "testserver5",
					ServerCapability: "servercap2",
					validationMessage: "already exists."
				}
			],
			remove: [
				{
					description: "remove server capability from server",
					Server: "testserver4",
					ServerCapability: "servercap1",
					validationMessage: "server server_capability was deleted."
				}
			],
			deleteServerCapability: [
				{
					description: "delete server capability linked with one or more servers",
					ServerCapability: "servercap2",
					validationMessage: "can not delete a server capability with 4 assigned servers"
				},
				{
					description: "delete server capabilities that is not link to any server",
					ServerCapability: "servercap4",
					validationMessage: "server capability was deleted."
				}
			]
		}
	]
};
