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
    tests: [
        {
            logins: [
                {
					description: "Admin Role",
					username: "TPAdmin",
					password: "pa$$word"
				}
            ],
            toggle:[
				{
					description: "hide first table column",
					Name: "1st Parent"
				},
				{
					description: "redisplay first table column",
					Name: "1st Parent"
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
                    validationMessage: "Delivery Service [ cdntesting ] updated"
                }
            ],
            remove: [
                {
                    description: "delete a delivery service",
                    Name: "tpdservice1",
                    validationMessage: "Delivery Service [ tpdservice1 ] deleted"
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
        }
    ]
}