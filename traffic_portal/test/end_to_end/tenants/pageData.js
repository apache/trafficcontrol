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
    this.createTenantButton = element(by.css('button[ng-click="createTenant()"]')); // this feels fragile, but there don't seem to be other good ways to select this button
    this.name = element(by.name("name"));
    this.active = element(by.name("active"));
    this.parent = element(by.name("parentId"));
    this.createButton = element(by.buttonText('Create'));
    this.successMsg = element(by.css('.alert-success'));
    this.tenantCreatedText = element(by.cssContainingText('div', 'Tenant created'));
    this.searchFilter=element(by.id('tenantsTable_filter')).element(by.css('label')).element(by.css('input'));
    this.updateButton = element(by.buttonText('Update'));
    this.tenantUpdatedText = element(by.cssContainingText('div', 'Tenant updated'));
    this.deleteButton = element(by.buttonText('Delete'));
    this.confirmWithNameInput = element(by.name('confirmWithNameInput'));
    this.deletePermanentlyButton = element(by.buttonText('Delete Permanently'));
    this.tenantDeletedText = element(by.cssContainingText('div', 'Tenant deleted'));
};
