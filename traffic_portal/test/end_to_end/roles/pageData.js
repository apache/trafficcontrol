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
    this.createRoleButton = element(by.css('button[title="Create Role"]'));
    this.name = element(by.name("name"));
    this.privLevel = element(by.name("privLevel"));
    this.description = element(by.css('textarea[name="description"]'));
    this.createButton = element(by.buttonText('Create'));
    this.successMsg = element(by.css('.alert-success'));
    this.roleCreatedText = element(by.cssContainingText('div', 'Role created')); // just a guess, keep getting Internal Server Error when running locally
    this.searchFilter=element(by.id('rolesTable_filter')).element(by.css('label')).element(by.css('input'));
    this.updateButton = element(by.buttonText('Update'));
    this.confirmUpdateButton = element(by.buttonText('Yes'));
    this.roleUpdatedText = element(by.cssContainingText('div', 'role was updated.'));
    this.deleteButton = element(by.buttonText('Delete'));
    this.confirmWithNameInput = element(by.name('confirmWithNameInput'));
    this.deletePermanentlyButton = element(by.buttonText('Delete Permanently'));
    this.roleDeletedText = element(by.cssContainingText('div', 'role was deleted.'));
};
