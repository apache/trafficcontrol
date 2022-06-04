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
                    description: "create a delivery service request then fullfill and complete the request",
                    XmlId: "cdntesting",
                    DisplayName: "testingoverload",
                    Active: "Active",
                    ContentRoutingType: "ANY_MAP",
                    Tenant: "tenantSame",
                    CDN: "dummycdn",
                    RawText: "test",
                    validationMessage: "Created request to create the cdntesting delivery service",
                    FullfillMessage: "Delivery Service creation was successful",
                    CompleteMessage: "Delivery service request status was updated"
                }
            ],
            remove: [
                {
                    description: "create a delivery service request then delete the request",
                    XmlId: "cdntesting2",
                    DisplayName: "testingoverload2",
                    Active: "Active",
                    ContentRoutingType: "ANY_MAP",
                    Tenant: "tenantSame",
                    CDN: "dummycdn",
                    RawText: "test",
                    validationMessage: "Created request to create the cdntesting2 delivery service",
                    DeleteMessage: "Delivery service request was deleted"
                }
            ],
            update: [
                {
                    description: "create a delivery service request then update the request then fullfill and complete the request",
                    XmlId: "cdntesting3",
                    DisplayName: "testingoverload2",
                    Active: "Active",
                    ContentRoutingType: "ANY_MAP",
                    Tenant: "tenantSame",
                    CDN: "dummycdn",
                    RawText: "test",
                    validationMessage: "Created request to create the cdntesting3 delivery service",
                    UpdateMessage: "Updated delivery service request for cdntesting3 and set status to submitted",
                    FullfillMessage: "Delivery Service creation was successful",
                    CompleteMessage: "Delivery service request status was updated"
                }
            ]
        }
    ]
}
