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

describe('Traffic Portal Origins Test Suite', function() {
    const pageData = new pd();
    const commonFunctions = new cfunc();
    const myNewOrigin = {
        name: "origin-" + commonFunctions.shuffle('abcdefghijklmnopqrstuvwxyz0123456789'),
        fdqn: "fake.origin.example.com",
        tenant: "- root",
        protocol: "http"
    };
    const tableRepeater = "o in ::origins";

    it('should create a new origin', async () => {
        console.log('Creating new origin');
        await browser.setLocation('origins');
        await pageData.createOriginButton.click();
        await pageData.name.sendKeys(myNewOrigin.name);
        await commonFunctions.selectDropdownByLabel(pageData.tenant, myNewOrigin.tenant);
        await pageData.fqdn.sendKeys(myNewOrigin.fdqn);
        await commonFunctions.selectDropdownByLabel(pageData.protocol, myNewOrigin.protocol);
        await commonFunctions.selectDropdownByNum(pageData.ds, 1);
        await pageData.createButton.click();
        expect(pageData.successMsg.isPresent()).toBe(true, 'Success alert message should exist');
        expect(pageData.originCreatedText.isPresent()).toBe(true, 'Actual message does not match expected message');
        expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl) + "#!/origins");
    });

    it('should update an existing origin', async () => {
        console.log('Updating existing origin');
        await browser.setLocation('origins');
        await commonFunctions.clickTableEntry(pageData.searchFilter, myNewOrigin.name, tableRepeater);
        await pageData.fqdn.clear().sendKeys('updated.' + myNewOrigin.fdqn);
        await pageData.updateButton.click();
        expect(pageData.successMsg.isPresent()).toBe(true, 'Success alert message should exist');
        expect(pageData.originUpdatedText.isPresent()).toBe(true, 'Actual message does not match expected message');
        expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl) + "#!/origins/[0-9]+$");
    });

    it('should delete an existing origin', async () => {
        console.log('Deleting an existing origin');
        await browser.setLocation('origins');
        await commonFunctions.clickTableEntry(pageData.searchFilter, myNewOrigin.name, tableRepeater);
        await pageData.deleteButton.click();
        await pageData.confirmWithNameInput.sendKeys(myNewOrigin.name);
        await pageData.deletePermanentlyButton.click();
        expect(pageData.successMsg.isPresent()).toBe(true, 'Success alert message should exist');
        expect(pageData.originDeletedText.isPresent()).toBe(true, 'Actual message does not match expected message');
        expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl) + "#!/origins");
    });

});
