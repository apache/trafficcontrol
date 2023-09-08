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
export const topologies = {
    cleanup: [
        {
            action: "DeleteTopologies",
            route: "/topologies",
            method: "delete",
            data: [
                {
                    route: "/topologies?name=TPTopTest2"
                }
            ]
        },
        {
            action: "DeleteServers",
            route : "/servers",
            method : "delete",
            data: [
                {
                    route: "/servers/",
                    getRequest: [
                        {
                            route: "/servers",
                            queryKey: "hostName",
                            queryValue: "topologieserver1",
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
                            queryValue: "topologieserver3",
                            replace: "route"
                        }
                    ]
                }
            ]
        },
        {
            action: "DeleteCacheGroup",
            route: "/cachegroups",
            method: "delete",
            data: [
                {
                    route: "/cachegroups/",
                    getRequest: [
                        {
                            route: "/cachegroups",
                            queryKey: "name",
                            queryValue: "TopoTestCGE1",
                            replace: "route"
                        }
                    ]
                },
                {
                    route: "/cachegroups/",
                    getRequest: [
                        {
                            route: "/cachegroups",
                            queryKey: "name",
                            queryValue: "TopoTestCGE2",
                            replace: "route"
                        }
                    ]
                }
            ]
        },
		{
			action: "DeleteProfile",
			route: "/profiles",
			method: "delete",
			data: [
				{
					route: "/profiles/",
					getRequest: [
						{
							route: "/profiles",
							queryKey: "name",
							queryValue: "TopTestPf",
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
							queryValue: "TopTestPhys",
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
					route: "/regions?name=TopTestReg"
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
							queryValue: "TopTestDiv",
							replace: "route"
						}
					]
				}
			]
		},
		{
			action: "DeleteCDN",
			route: "/cdns",
			method: "delete",
			data: [
				{
					route: "/cdns/",
					getRequest: [
						{
							route: "/cdns",
							queryKey: "name",
							queryValue: "TopTestCDN",
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
					name: "TopTestDiv"
				}
			]
		},
		{
			action: "CreateRegions",
			route: "/regions",
			method: "post",
			data: [
				{
					name: "TopTestReg",
					division: "4",
					divisionName: "TopTestDiv",
					getRequest: [
						{
							route: "/divisions",
							queryKey: "name",
							queryValue: "TopTestDiv",
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
					name: "TopTestPhys",
					phone: "0-843-816-6276",
					poc: "Her Majesty The Queen Elizabeth Alexandra Mary Windsor II",
					regionId: 3,
					shortName: "ttphys",
					state: "NA",
					zip: "99999",
					getRequest: [
						{
							route: "/regions",
							queryKey: "name",
							queryValue: "TopTestReg",
							replace: "regionId"
						}
					]
				}
			]
		},
        {
			action: "CreateCDN",
			route: "/cdns",
			method: "post",
			data: [
				{
					name: "TopTestCDN",
					domainName: "ttcdn",
					dnssecEnabled: false
				}
			]
		},
		{
			action: "CreateProfile",
			route: "/profiles",
			method: "post",
			data: [
				{
					name: "TopTestPf",
					description: "A test profile for API examples",
					cdn: 2,
					type: "UNK_PROFILE",
					routingDisabled: true,
					getRequest: [
						{
							route: "/cdns",
							queryKey: "name",
							queryValue: "TopTestCDN",
							replace: "cdn"
						}
					]
				}
			]
		},
        {
            action: "CreateCacheGroup",
            route: "/cachegroups",
            method: "post",
            data: [
                {
                    name: "TopoTestCGE1",
                    shortName: "CGEdge",
                    latitude: 0,
                    longitude: 0,
                    fallbackToClosest: true,
                    localizationMethods: [
                        "DEEP_CZ",
                        "CZ",
                        "GEO"
                    ],
                    typeId: 23
                },
                {
                    name: "TopoTestCGE2",
                    shortName: "CGEdge2",
                    latitude: 0,
                    longitude: 0,
                    fallbackToClosest: true,
                    localizationMethods: [
                        "DEEP_CZ",
                        "CZ",
                        "GEO"
                    ],
                    typeId: 23
                },
                {
                    name: "TopoTestCGE3",
                    shortName: "CGEdge3",
                    latitude: 0,
                    longitude: 0,
                    fallbackToClosest: true,
                    localizationMethods: [
                        "DEEP_CZ",
                        "CZ",
                        "GEO"
                    ],
                    typeId: 23
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
                    hostName: "topologieserver1",
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
                    profiles: ["TopTestPf"],
                    routerHostName: "",
                    routerPortName: "",
                    statusId: 3,
                    tcpPort: 80,
                    typeID: 12,
                    getRequest: [
                        {
                            route: "/cachegroups",
                            queryKey: "name",
                            queryValue: "TopoTestCGE2",
                            replace: "cacheGroupID"
                        },
                        {
							route: "/phys_locations",
							queryKey: "name",
							queryValue: "TopTestPhys",
							replace: "physicalLocationID"
						},
                        {
							route: "/cdns",
							queryKey: "name",
							queryValue: "TopTestCDN",
							replace: "cdnID"
						}
                    ],
                    updPending: false
                },
                {
                    cacheGroupID: 0,
                    cdnID: 0,
                    domainName: "test.net",
                    hostName: "topologieserver3",
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
                    profiles: ["TopTestPf"],
                    routerHostName: "",
                    routerPortName: "",
                    statusId: 3,
                    tcpPort: 80,
                    typeID: 12,
                    getRequest: [
                        {
                            route: "/cachegroups",
                            queryKey: "name",
                            queryValue: "TopoTestCGE3",
                            replace: "cacheGroupID"
                        },
                        {
							route: "/phys_locations",
							queryKey: "name",
							queryValue: "TopTestPhys",
							replace: "physicalLocationID"
						},
                        {
							route: "/cdns",
							queryKey: "name",
							queryValue: "TopTestCDN",
							replace: "cdnID"
						}
                    ],
                    updPending: false
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
                    description: "create a Topologies with empty cachegroup (no server)",
                    Name: "TPTopTest1",
                    DescriptionData: "test",
                    Type: "EDGE_LOC",
                    CacheGroup: "TopoTestCGE1",
                    validationMessage: "'empty cachegroups' cachegroups with no servers in them: TopoTestCGE1"
                },
                {
                    description: "create a Topologies with cachegroup has server in it",
                    Name: "TPTopTest2",
                    DescriptionData: "test",
                    Type: "EDGE_LOC",
                    CacheGroup: "TopoTestCGE2",
                    validationMessage: "topology was created."
                },
                {
                    description: "create a Topologies with no cache group in it",
                    Name: "TPTopTest3",
                    DescriptionData: "test",
                    Type: "EDGE_LOC",
                    CacheGroup: "wrong",
                    validationMessage: "'length' must provide 1 or more node, 0 found"
                }
            ]
        }
    ]
}
