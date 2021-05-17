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
                    description: "check CSV link from CDN page",
                    Name: "Export as CSV"
                }
            ],
            add: [
                {
                    description: "create a new User",
                    FullName: "TPCreateUser1",
                    Username: "User1",
                    Email: "test@cdn.tc.com",
                    Role: "admin",
                    Tenant: "- root",
                    Password: "qwe@123#rty",
                    ConfirmPassword: "qwe@123#rty",
                    PublicSSHKey: "",
                    validationMessage: "User created"
                },
                // {
                //     description: "create a registered User",
                //     Email: "test1@cdn.tc.com",
                //     Role: 1,
                //     Tenant: 1,
                //     validationMessage: "Registered User created."
                // }
            ],
            update: [
                {
                    description: "update the new User",
                    Username: "User1",
                    NewFullName: "TPUpdatedUser1`",
                    validationMessage: "user was updated."
                },
                // {
                //     description: "update the registered User",
                //     Email: "test1@cdn.tc.com",
                //     FullName: "TPCreateUser2",
                //     validationMessage: "Registered User updated."
                // }
            ],
        },
    ]
};