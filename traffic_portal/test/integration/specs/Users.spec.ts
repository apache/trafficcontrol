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

import { browser } from 'protractor';

import { LoginPage } from '../PageObjects/LoginPage.po';
import {UsersPage} from '../PageObjects/UsersPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
// import { API } from '../CommonUtils/API';
import { users } from "../Data/users";


// const api = new API();
const loginPage = new LoginPage();
const topNavigation = new TopNavigationPage();
const usersPage = new UsersPage();

users.tests.forEach(async usersData =>{
    usersData.logins.forEach(login => {
        describe('Traffic Portal - Users - ' + login.description, function(){
            it('can login', async () => {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login);
                expect(await loginPage.CheckUserName(login)).toBeTruthy();
            });
            it('can open users page', async () => {
                await usersPage.OpenUserPage();
            });
            usersData.check.forEach(check => {
                it(check.description, async () => {
                    expect(await usersPage.CheckCSV(check.Name)).toBe(true);
                    await usersPage.OpenUserPage();
                });
            });
            usersData.add.forEach(add => {
                it(add.description, async () => {
                    expect(await usersPage.CreateUser(add)).toBeTruthy();
                    await usersPage.OpenUserPage();
                });
            });
            usersData.update.forEach(update => {
                it(update.description, async () => {
                    await usersPage.SearchUser(update.Username);
                    expect(await usersPage.UpdateUser(update)).toBeTruthy();
                });
            });
            it('can logout', async () => {
                expect(await topNavigation.Logout()).toBeTruthy();
            });
        });
    });
});