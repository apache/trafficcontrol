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
export const tenants = {
    cleanup: [
        {
            action: "DeleteTenants",
            route: "/tenants",
            method: "delete",
            data: [
                {
                    route: "/tenants/",
                    getRequest: [
                        {
                            route: "/tenants",
                            queryKey: "name",
                            queryValue: "TPTestReadOnly",
                            replace: "route"
                        }
                    ]
                }
            ]
        }
    ],
    setup: [
        {
            action: "CreateTenants",
            route: "/tenants",
            method: "post",
            data: [
                {
                    active: true,
                    name: "TPTestReadOnly",
                    parentId: 1
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
            add: [
                {
                    description: "create a tenant",
                    Name: "TPTenantTest",
                    Active: "true",
                    ParentTenant: "tenantSame",
                    validationMessage: "Tenant created"
                }
            ],
            update: [
                {
                    description: "update a tenant",
                    Name: "TPTenantTest",
                    Active: "false",
                    validationMessage: "Tenant updated"
                }
            ],
            remove: [
                {
                    description: "delete a tenant",
                    Name: "TPTenantTest",
                    validationMessage: "tenant was deleted."
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
                    description: "create a tenant",
                    Name: "TPTenantTest",
                    Active: "true",
                    ParentTenant: "tenantSame",
                    validationMessage: "missing required Permissions: TENANT:CREATE"
                }
            ],
            update: [
                {
                    description: "update a tenant",
                    Name: "tenantChild",
                    Active: "false",
                    validationMessage: "missing required Permissions: TENANT:UPDATE"
                }
            ],
            remove: [
                {
                    description: "delete a tenant",
                    Name: "tenantChild",
                    validationMessage: "missing required Permissions: TENANT:DELETE"
                }
            ]
        }
    ]
}
