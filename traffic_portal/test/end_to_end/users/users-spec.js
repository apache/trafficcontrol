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

describe('Traffic Portal Users Test Suite', () => {
    const pageData = new pd();
    const commonFunctions = new cfunc();
    const newUser = function() {
        return {
            name: 'User Name',
            username: 'username-' + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz0123456789'),
            email: 'user-' + commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz0123456789') + '@example.com',
            password: commonFunctions.shuffle('abcdefghijklmonpqrstuvwxyz0123456789')
        };
    };
    const myNewUser = newUser();
    const myOtherNewUser = newUser();

    it('should register a new user', async () => {
        await browser.setLocation('users');
        await pageData.registerUserButton.click();
        await pageData.email.sendKeys(myNewUser.email);
        commonFunctions.selectDropdownbyNum(pageData.role, 3);
        commonFunctions.selectDropdownbyNum(pageData.tenant, 1);
        await pageData.sendRegistration.click();
        //expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl) + "#!/users");
    });

    it('should create a new user', async () => {
        browser.setLocation('users');
        await pageData.createUserButton.click();
        await pageData.fullName.sendKeys(myOtherNewUser.name);
        await pageData.username.sendKeys(myOtherNewUser.username);
        await pageData.email.sendKeys(myOtherNewUser.email);
        commonFunctions.selectDropdownbyNum(pageData.role, 3);
        commonFunctions.selectDropdownbyNum(pageData.tenant, 1);
        await pageData.password.sendKeys(myOtherNewUser.password);
        await pageData.confirmPassword.sendKeys(myOtherNewUser.password);
        await pageData.createButton.click();
        expect(element(by.css('.alert-success')).isPresent()).toBe(true);
        expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl) + "#!/users");
    });

    it('should update the username of an existing user', async () => {
        browser.setLocation('users');
        await pageData.searchFilter.clear().sendKeys(myNewUser.email);
        await element.all(by.repeater('u in ::users')).get(0).click();
        await pageData.fullName.clear().sendKeys(myNewUser.name);
        await pageData.username.clear().sendKeys(myNewUser.username);
        await pageData.updateButton.click();
        expect(element(by.css('.alert-success')).isPresent()).toBe(true);
        expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toMatch(commonFunctions.urlPath(browser.baseUrl) + "#!/users/[0-9]+$");
    });

});