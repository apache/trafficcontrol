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

var pd = require('./pageData.js');
var cfunc = require('../common/commonFunctions.js');

describe('Traffic Portal Tenants Test Suite', function() {
    const pageData = new pd();
    const commonFunctions = new cfunc();
    const myNewTenant = {
        name: "tenant-" + commonFunctions.shuffle('abcdefghijklmnopqrstuvwxyz0123456789'),
        active: true,
        parent: "- root"
    };
    const tableRepeater = "t in ::tenants";

    it('should create a new tenant', async () => {
        console.log('Creating new tenant');
        await browser.setLocation('tenants');
        await pageData.createTenantButton.click();
        await pageData.name.sendKeys(myNewTenant.name);
        await commonFunctions.selectDropdownByLabel(pageData.active, myNewTenant.active.toString());
        await commonFunctions.selectDropdownByLabel(pageData.parent, myNewTenant.parent);
        await pageData.createButton.click();
        expect(pageData.successMsg.isPresent()).toBe(true, 'Success alert message should exist');
        expect(pageData.tenantCreatedText.isPresent()).toBe(true, 'Actual message does not match expected message');
        expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl) + "#!/tenants");
    });

    it('should update an existing tenant', async () => {
        console.log('Updating existing tenant');
        await browser.setLocation('tenants');
        await commonFunctions.clickTableEntry(pageData.searchFilter, myNewTenant.name, tableRepeater);
        await commonFunctions.selectDropdownByLabel(pageData.active, !myNewTenant.active.toString());
        await pageData.updateButton.click();
        expect(pageData.successMsg.isPresent()).toBe(true, 'Success alert message should exist');
        expect(pageData.tenantUpdatedText.isPresent()).toBe(true, 'Actual message does not match expected message');
        expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl) + "#!/tenants/[0-9]+$");
    });

    it('should delete an existing tenant', async () => {
        console.log('Deleting an existing tenant');
        await browser.setLocation('tenants');
        await commonFunctions.clickTableEntry(pageData.searchFilter, myNewTenant.name, tableRepeater);
        await pageData.deleteButton.click();
        await pageData.confirmWithNameInput.sendKeys(myNewTenant.name);
        await pageData.deletePermanentlyButton.click();
        expect(pageData.successMsg.isPresent()).toBe(true, 'Success alert message should exist');
        expect(pageData.tenantDeletedText.isPresent()).toBe(true, 'Actual message does not match expected message');
        expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl) + "#!/tenants");
    });

});
