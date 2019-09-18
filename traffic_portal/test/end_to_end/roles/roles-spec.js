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

describe('Traffic Portal Roles Test Suite', function() {
    const pageData = new pd();
    const commonFunctions = new cfunc();
    const myNewRole = {
        name: "role-" + commonFunctions.shuffle('abcdefghijklmnopqrstuvwxyz0123456789'),
        privLevel: 30,
        description: "This is my new role"
    };
    const tableRepeater = "r in ::roles";

    it('should create a new role', async () => {
        console.log('Creating new role');
        await browser.setLocation('roles');
        await pageData.createRoleButton.click();
        await pageData.name.sendKeys(myNewRole.name);
        await pageData.privLevel.sendKeys(myNewRole.privLevel);
        await pageData.description.click().sendKeys(myNewRole.description);
        await pageData.createButton.click();
        expect(pageData.successMsg.isPresent()).toBe(true, 'Success alert message should exist');
        expect(pageData.roleCreatedText.isPresent()).toBe(true, 'Actual message does not match expected message');
        expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl) + "#!/roles");
    });

    it('should update an existing role', async () => {
        console.log('Updating existing role');
        await browser.setLocation('roles');
        await commonFunctions.clickTableEntry(pageData.searchFilter, myNewRole.name, tableRepeater);
        await pageData.privLevel.clear().sendKeys((myNewRole.privLevel / 2).toString());
        await pageData.updateButton.click();
        await pageData.confirmUpdateButton.click();
        expect(pageData.successMsg.isPresent()).toBe(true, 'Success alert message should exist');
        expect(pageData.roleUpdatedText.isPresent()).toBe(true, 'Actual message does not match expected message');
        expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl) + "#!/roles/[0-9]+$");
    });

    it('should delete an existing role', async () => {
        console.log('Deleting an existing role');
        await browser.setLocation('roles');
        await commonFunctions.clickTableEntry(pageData.searchFilter, myNewRole.name, tableRepeater);
        await pageData.deleteButton.click();
        await pageData.confirmWithNameInput.sendKeys(myNewRole.name);
        await pageData.deletePermanentlyButton.click();
        expect(pageData.successMsg.isPresent()).toBe(true, 'Success alert message should exist');
        expect(pageData.roleDeletedText.isPresent()).toBe(true, 'Actual message does not match expected message');
        expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl) + "#!/roles");
    });

});
