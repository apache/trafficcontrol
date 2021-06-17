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
			action: "CreateDeliveryServices",
			route: "/deliveryservices",
			method: "post",
			data: [
				{
					active: true,
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
					cachegroupId: 0,
					cdnId: 0,
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
					physLocationId: 0,
					profileId: 0,
					routerHostName: "",
					routerPortName: "",
					statusId: 3,
					tcpPort: 80,
					typeId: 11,
					updPending: false,
					getRequest: [
						{
							route: "/phys_locations",
							queryKey: "name",
							queryValue: "DSTest",
							replace: "physLocationId"
						},
						{
							route: "/cdns",
							queryKey: "name",
							queryValue: "dummycdn",
							replace: "cdnId"
						},
						{
							route: "/cachegroups",
							queryKey: "name",
							queryValue: "testCG",
							replace: "cachegroupId"
						},
						{
							route: "/profiles",
							queryKey: "name",
							queryValue: "testProfile",
							replace: "profileId"
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
				}
            ],
            add: [
                {
                    description: "create ANY_MAP delivery service",
                    Name: "tpdservice1",
                    Type: "ANY_MAP",
                    validationMessage: "Delivery Service [ tpdservice1 ] created"
                },
                {
                    description: "create DNS delivery service",
                    Name: "tpdservice2",
                    Type: "DNS",
                    validationMessage: "Delivery Service [ tpdservice2 ] created"
                },
                {
                    description: "create STEERING delivery service",
                    Name: "tpdservice3",
                    Type: "STEERING",
                    validationMessage: "Delivery Service [ tpdservice3 ] created"
                }
            ],
            update: [
                {
                    description: "update delivery service display name",
                    Name: "tpdservice1",
                    NewName: "TPServiceNew1",
                    validationMessage: "Delivery Service [ tpdservice1 ] updated"
                }
            ],
            assignserver: [
                {
                    description: "assign server to delivery service",
                    ServerName: "DSTest",
                    DSName: "TPServiceNew1",
                    validationMessage: "server assignments complete"
                }
            ],
            assignrequiredcapabilities: [
                {
                    description: "assign required capabilities to delivery service",
                    RCName: "DSTestCap",
                    DSName: "tpdservice2",
                    validationMessage: "deliveryservice.RequiredCapability was created."
                }
            ],
            remove: [
                {
                    description: "delete a delivery service",
                    Name: "tpdservice1",
                    validationMessage: "Delivery service [ tpdservice1 ] deleted"
                },
                {
                    description: "delete a delivery service",
                    Name: "tpdservice2",
                    validationMessage: "Delivery service [ tpdservice2 ] deleted"
                },
                {
                    description: "delete a delivery service",
                    Name: "tpdservice3",
                    validationMessage: "Delivery service [ tpdservice3 ] deleted"
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
                    description: "create ANY_MAP delivery service",
                    Name: "tpdservice1",
                    Type: "ANY_MAP",
                    validationMessage: "Forbidden."
                }
            ],
            update: [
                {
                    description: "update delivery service display name",
                    Name: "dstestro1",
                    NewName: "TPServiceNew1",
                    validationMessage: "Forbidden."
                }
            ],
            assignserver: [
                {
                    description: "assign server to delivery service",
                    ServerName: "DSTest",
                    DSName: "dstestro1",
                    validationMessage: "Forbidden."
                }
            ],
            assignrequiredcapabilities: [
                {
                    description: "assign required capabilities to delivery service",
                    RCName: "DSTestCap",
                    DSName: "dstestro1",
                    validationMessage: "Forbidden."
                }
            ],
            remove: [
                {
                    description: "delete a delivery service",
                    Name: "dstestro1",
                    validationMessage: "Forbidden."
                }
            ]
        }
    ]
}
