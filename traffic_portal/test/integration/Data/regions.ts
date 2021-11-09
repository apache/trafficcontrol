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
export const regions = {
	cleanup: [
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
                            queryValue: "TestDivision1",
                            replace: "route"
                        }
                    ]
                },
                {
                    route: "/divisions/",
                    getRequest: [
                        {
                            route: "/divisions",
                            queryKey: "name",
                            queryValue: "TestDivision2",
                            replace: "route"
                        }
                    ]
                }
            ]
        }
    ],
	setup: [
        {
            action: "CreateDivision",
            route: "/divisions",
            method: "post",
            data: [
                {
                    name: "TestDivision1"
                },
                {
                    name: "TestDivision2"
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
            check: [
				{
					description: "check CSV link from Regions page",
					Name: "Export as CSV"
				}
			],
            add: [
                {
                    description: "create a Regions",
                    Name: "TPRegion1",
                    Division: "TestDivision1",
                    validationMessage: "Region created"
                },
                {
                    description: "create multiple Regions",
                    Name: "TPRegion2",
                    Division: "TestDivision2",
                    validationMessage: "Region created"
                }
            ],
            update: [
                {
                    description: "update Region's Division",
                    Name: "TPRegion1",
                    Division: "TestDivision2",
                    validationMessage: "Region updated"
                }
            ],
            remove: [
                {
                    description: "delete a Region",
                    Name: "TPRegion1",
                    validationMessage: "Region deleted"
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
					description: "check CSV link from Regions page",
					Name: "Export as CSV"
				}
			],
            add: [
                {
                    description: "create a Regions",
                    Name: "TPRegion1",
                    Division: "TestDivision1",
                    validationMessage: "missing required Permissions: REGION:CREATE"
                }
            ],
            update: [
                {
                    description: "update Region's Division",
                    Name: "TPRegion2",
                    Division: "TestDivision1",
                    validationMessage: "missing required Permissions: REGION:UPDATE"
                }
            ],
            remove: [
                {
                    description: "delete a Region",
                    Name: "TPRegion2",
                    validationMessage: "missing required Permissions: REGION:DELETE"
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
            check: [
				{
					description: "check CSV link from Regions page",
					Name: "Export as CSV"
				}
			],
            add: [
                {
                    description: "create a Regions",
                    Name: "TPRegion3",
                    Division: "TestDivision1",
                    validationMessage: "Region created"
                },
                {
                    description: "create multiple Regions",
                    Name: "TPRegion4",
                    Division: "TestDivision2",
                    validationMessage: "Region created"
                }
            ],
            update: [
                {
                    description: "update Region's Division",
                    Name: "TPRegion3",
                    Division: "TestDivision2",
                    validationMessage: "Region updated"
                }
            ],
            remove: [
                {
                    description: "delete a Region",
                    Name: "TPRegion2",
                    validationMessage: "Region deleted"
                },
                {
                    description: "delete a Region",
                    Name: "TPRegion3",
                    validationMessage: "Region deleted"
                },
                {
                    description: "delete a Region",
                    Name: "TPRegion4",
                    validationMessage: "Region deleted"
                }
            ]
        }
    ]
};
