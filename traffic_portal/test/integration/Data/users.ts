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

import {randomize} from "../config";

export const users = {
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
                    description: "check CSV link from Users page",
                    Name: "Export as CSV"
                }
            ],
            toggle: [
                {
                    description: "hide first table column",
                    Name: "Email"
                },
                {
                    description: "redisplay first table column",
                    Name: "Email"
                }
            ],
            add: [
                {
                    description: "create a new User",
                    FullName: "TPCreateUser1",
                    Username: "User1",
                    Email: "test@cdn.tc.com",
                    Role: "admin",
                    Tenant: "tenantSame",
                    UCDN: "",
                    Password: "qwe@123#rty",
                    ConfirmPassword: "qwe@123#rty",
                    PublicSSHKey: "",
                    validationMessage: "User created"
                },
            ],
            register: [
                {
                    description: "create a registered User",
                    Email: "test2@cdn.tc.com",
                    Role: "admin",
                    Tenant: "tenantSame",
                    validationMessage: `Sent user registration to {{ ${randomize}test2@cdn.tc.com}} with the following permissions [ role: admin | tenant: tenantSame${randomize} ]`
                }
            ],
            update: [
                {
                    description: "update user's fullname",
                    Username: "User1",
                    NewFullName: "TPUpdatedUser1",
                    validationMessage: "user was updated."
                },
            ],
            updateRegisterUser: [
                {
                    description: "update registered user's fullname",
                    Email: "test2@cdn.tc.com",
                    NewFullName: "TPRegisterUser1",
                    validationMessage: "user was updated."
                }
            ],
        },
    ]
};
