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
export const deliveryservices = {
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
							queryValue: "DSTest",
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
							queryValue: "DSTest",
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
					route: "/regions?name=DSTest"
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
							queryValue: "DSTest",
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
					route: "/server_capabilities?name=DSTestCap"
				}
			]
		},
		{
			action: "DeleteDeliveryServices",
			route: "/deliveryservices",
			method: "delete",
			data: [
				{
					route: "/deliveryservices/",
					getRequest: [
						{
							route: "/deliveryservices",
							queryKey: "xmlId",
							queryValue: "dstestro1",
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
					name: "location",
					value: "/a/b/c/d",
					configFile: "remap.config",
					secure: false
				}
			]
		},
		{
			action: "CreateServerCapabilities",
			route: "/server_capabilities",
			method: "post",
			data: [
				{
					name: "DSTestCap"
				}
			]
		},
		{
			action: "CreateDeliveryServices",
			route: "/deliveryservices",
			method: "post",
			data: [
				{
					active: "PRIMED",
					cdnId: 0,
					displayName: "DSTestReadOnly",
					dscp: 0,
					geoLimit: 0,
					geoProvider: 0,
					initialDispersion: 1,
					ipv6RoutingEnabled: true,
					logsEnabled: false,
					missLat: 41.881944,
					missLong: -87.627778,
					multiSiteOrigin: false,
					orgServerFqdn: "http://origin.infra.ciab.test",
					protocol: 0,
					qstringIgnore: 0,
					rangeRequestHandling: 0,
					regionalGeoBlocking: false,
					requiredCapabilities: [],
					tenantId: 0,
					typeId: 1,
					xmlId: "dstestro1",
					getRequest: [
						{
							route: "/tenants",
							queryKey: "name",
							queryValue: "tenantSame",
							replace: "tenantId"
						},
						{
							route: "/cdns",
							queryKey: "name",
							queryValue: "dummycdn",
							replace: "cdnId"
						}
					]
				}
			]
		},
		{
			action: "CreateDivisions",
			route: "/divisions",
			method: "post",
			data: [
				{
					name: "DSTest"
				}
			]
		},
		{
			action: "CreateRegions",
			route: "/regions",
			method: "post",
			data: [
				{
					name: "DSTest",
					division: "4",
					divisionName: "DSTest",
					getRequest: [
						{
							route: "/divisions",
							queryKey: "name",
							queryValue: "DSTest",
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
					name: "DSTest",
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
							queryValue: "DSTest",
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
					cachegroupID: 0,
					cdnID: 0,
					domainName: "test.net",
					hostName: "DSTest",
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
					typeID: 11,
					getRequest: [
						{
							route: "/phys_locations",
							queryKey: "name",
							queryValue: "DSTest",
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
							replace: "cachegroupID"
						}
					]
				}
			]
		}
	],
	tests: [
		{
			description: "Admin Role Delivery Service actions",
			login: {
				username: "TPAdmin",
				password: "pa$$word"
			},
			add: [
				{
					description: "create ANY_MAP delivery service",
					name: "tpdservice1",
					tenant: "tenantSame",
					type: "ANY_MAP",
					validationMessage: "Delivery Service creation was successful"
				},
				{
					description: "create DNS delivery service",
					name: "tpdservice2",
					tenant: "tenantSame",
					type: "DNS",
					validationMessage: "Delivery Service creation was successful"
				},
				{
					description: "create STEERING delivery service",
					name: "tpdservice3",
					tenant: "tenantSame",
					type: "STEERING",
					validationMessage: "Delivery Service creation was successful"
				}
			],
			update: [
				{
					name: "tpdservice1",
					newName: "TPServiceNew1",
					validationMessage: "Delivery Service update was successful"
				}
			],
			assignServer: [
				{
					serverHostname: "DSTest",
					xmlId: "TPServiceNew1",
					validationMessage: "server assignments complete"
				}
			],
			remove: [
				{
					name: "tpdservice1",
					validationMessage: "ds was deleted."
				},
				{
					name: "tpdservice2",
					validationMessage: "ds was deleted."
				},
				{
					name: "tpdservice3",
					validationMessage: "ds was deleted."
				}
			]
		},
		{
			description: "Read Only Role Delivery Service actions",
			login: {
				username: "TPReadOnly",
				password: "pa$$word"
			},
			add: [
				{
					description: "create ANY_MAP delivery service",
					name: "tpdservice1",
					type: "ANY_MAP",
					tenant: "tenantSame",
					validationMessage: "missing required Permissions: DELIVERY-SERVICE:CREATE"
				}
			],
			update: [
				{
					name: "dstestro1",
					newName: "TPServiceNew1",
					validationMessage: "missing required Permissions: DELIVERY-SERVICE:UPDATE"
				}
			],
			assignServer: [
				{
					serverHostname: "DSTest",
					xmlId: "dstestro1",
					validationMessage: "missing required Permissions: SERVER:UPDATE, DELIVERY-SERVICE:UPDATE"
				}
			],
			remove: [
				{
					name: "dstestro1",
					validationMessage: "missing required Permissions: DELIVERY-SERVICE:DELETE"
				}
			]
		},
		{
			description: "Operation Role Delivery Service actions",
			login: {
				username: "TPOperator",
				password: "pa$$word"
			},
			add: [
				{
					description: "create ANY_MAP delivery service",
					name: "optpdservice1",
					tenant: "tenantSame",
					type: "ANY_MAP",
					validationMessage: "Delivery Service creation was successful"
				},
				{
					description: "create DNS delivery service",
					name: "optpdservice2",
					tenant: "tenantSame",
					type: "DNS",
					validationMessage: "Delivery Service creation was successful"
				},
				{
					description: "create STEERING delivery service",
					name: "optpdservice3",
					tenant: "tenantSame",
					type: "STEERING",
					validationMessage: "Delivery Service creation was successful"
				}
			],
			update: [
				{
					name: "optpdservice1",
					newName: "opTPServiceNew1",
					validationMessage: "Delivery Service update was successful"
				}
			],
			assignServer: [
				{
					serverHostname: "DSTest",
					xmlId: "opTPServiceNew1",
					validationMessage: "server assignments complete"
				}
			],
			remove: [
				{
					name: "optpdservice1",
					validationMessage: "ds was deleted."
				},
				{
					name: "optpdservice2",
					validationMessage: "ds was deleted."
				},
				{
					name: "optpdservice3",
					validationMessage: "ds was deleted."
				}
			]
		}
	]
}
