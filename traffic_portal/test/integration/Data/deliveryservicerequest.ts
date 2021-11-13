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
export const deliveryservicerequest = {
    cleanup: [
       
    ],
    setup: [
        {
            action: "CreateDeliveryServiceRequest",
            route: "/deliveryservice_requests",
            method: "post",
            data: [
                { 
                    changeType: "create", 
                    status: "submitted", 
                    requested: { 
                        dscp: 0, 
                        regionalGeoBlocking: false, 
                        logsEnabled: false, 
                        geoProvider: 0, 
                        geoLimit: 0, 
                        ccrDnsTtl: 30, 
                        anonymousBlockingEnabled: false, 
                        consistentHashQueryParams: [], 
                        xmlId: "test212", 
                        displayName: "testing212", 
                        active: true, 
                        typeId: 8, 
                        tenantId: 1, 
                        cdnId: 1, 
                        remapText: "test", 
                        tlsVersions: null 
                    } 
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
            create: [
                {
                    description: "create a delivery service request",
                    XmlId: "cdntesting",
                    DisplayName: "testingoverload",
                    Active: "Active",
                    ContentRoutingType: "ANY_MAP",
                    Tenant: "-tenantSame",
                    CDN: "dummycdn",
                    RawText: "test",
                    validationMessage: "Created request to create the cdntesting delivery service",
                    FullfillMessage: "Delivery Service [ cdntesting ] created"
                }
            ]

        }
    ]
}