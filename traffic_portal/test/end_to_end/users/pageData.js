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

module.exports = function () {
    this.registerUserButton = element(by.buttonText('Register User'));
    this.email = element(by.name("email"));
    this.role = element(by.name("role"));
    this.tenant = element(by.name("tenantId"));
    this.sendRegistration = element(by.buttonText('Send Registration'));
    this.successMsg = element(by.css('.alert-success'));
    this.userRegisteredText = element(by.cssContainingText('div', 'User was registered.')); // not actually sure if that's right, because I can't get a successful user registration to happen
    this.createUserButton = element(by.css('button[title="Create New User"]'));
    this.fullName = element(by.name("fullName"));
    this.username = element(by.name("uName"));
    this.password = element(by.name("uPass"));
    this.confirmPassword = element(by.name("confirmPassword"));
    this.createButton = element(by.buttonText('Create'));
    this.userCreatedText = element(by.cssContainingText('div', 'User created'));
    this.searchFilter=element(by.id('usersTable_filter')).element(by.css('label')).element(by.css('input'));
    this.updateButton = element(by.buttonText('Update'));
    this.userUpdatedText = element(by.cssContainingText('div', 'user was updated.'));
};
