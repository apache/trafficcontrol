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

import { by, element } from 'protractor';
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';
import { randomize } from '../config';
interface User {
  FullName: string;
  Username: string;
  Email: string;
  Role: string;
  Tenant: string;
  UCDN: string;
  Password: string;
  ConfirmPassword: string;
  PublicSSHKey: string;
  existsMessage?: string;
  validationMessage?: string;
}

interface UpdateUser {
    description: string;
    Username: string;
    NewFullName: string;
    validationMessage?: string;
}

interface RegisterUser {
    Email: string;
    Role: string;
    Tenant: string;
    existsMessage?: string;
    validationMessage?: string;
}

interface UpdateRegisterUser {
    description: string;
    Email: string;
    NewFullName: string;
    validationMessage?: string;
}
export class UsersPage extends BasePage {
    private btnCreateNewUser = element(by.name('createUserButton'));
    private btnRegisterNewUser = element(by.name('createRegisterUserButton'));
    private txtFullName = element(by.name('fullName'));
    private txtUserName = element(by.name('uName'));
    private txtEmail = element(by.name('email'));
    private txtRole = element(by.name('role'));
    private selTenant = element(by.name('tenantId'));
    private txtUCDN = element(by.name('uCDN'));
    private txtPassword = element(by.name('uPass'));
    private txtConfirmPassword = element(by.name('confirmPassword'));
    private txtPublicSSHKey = element(by.name('publicSshKey'));
    private txtSearch = element(by.id('usersTable_filter')).element(by.css('label input'));
    private btnTableColumn = element(by.css('[title="Select Table Columns"]'));
    private randomize = randomize;

    public async OpenUserPage(): Promise<void> {
        const snp = new SideNavigationPage();
        await snp.NavigateToUsersPage();
    }

    public async OpenUserMenu(): Promise<void> {
        const snp = new SideNavigationPage();
        await snp.ClickUserAdminMenu();
    }

    public async CheckCSV(name: string): Promise<boolean> {
        return element(by.cssContainingText("span", name)).isPresent();
    }

    public async CheckToggle(name: string): Promise<boolean> {
        let result = false;
        await this.btnTableColumn.click();
        //if the box is already checked, uncheck it and return false
        if (await element(by.cssContainingText("th", name)).isPresent()) {
            await element(by.cssContainingText("label", name)).click();
            result = false;
        } else {
            //if the box is unchecked, then check it and return true
            await element(by.cssContainingText("label", name)).click();
            result = true;
        }
        await this.btnTableColumn.click();
        return result;
    }

    public async CreateUser(user: User): Promise<boolean> {
      let result = false;
      const basePage = new BasePage();
      const snp = new SideNavigationPage();
      await this.btnCreateNewUser.click();
      await this.txtFullName.sendKeys(user.FullName + this.randomize);
      await this.txtUserName.sendKeys(user.Username + this.randomize);
      await this.txtEmail.sendKeys(this.randomize + user.Email);
      await this.txtRole.sendKeys(user.Role);
      await this.selTenant.click();
      await element(by.name(user.Tenant+this.randomize)).click();
      await this.txtUCDN.sendKeys(user.UCDN);
      await this.txtPassword.sendKeys(user.Password);
      await this.txtConfirmPassword.sendKeys(user.ConfirmPassword);
      await this.txtPublicSSHKey.sendKeys(user.PublicSSHKey);
      await basePage.ClickCreate();
      if(await basePage.GetOutputMessage() === user.existsMessage){
        await snp.NavigateToUsersPage();
        result = true;
      }else if(await basePage.GetOutputMessage() === user.validationMessage){
        result = true;
      }else{
        result = false;
      }
      return result;
    }

    public async SearchUser(nameUser: string): Promise<void> {
        const snp = new SideNavigationPage();
        const name = nameUser + this.randomize;
        await snp.NavigateToUsersPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        await element.all(by.repeater('u in ::users')).filter(function (row) {
            return row.element(by.name('username')).getText().then(function (val) {
                return val === name;
            });
        }).first().click();
    }

    public async SearchEmailUser(nameEmail: string): Promise<void> {
        const snp = new SideNavigationPage();
        const name = this.randomize + nameEmail;
        await snp.NavigateToUsersPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        await element.all(by.repeater('u in ::users')).filter(function (row) {
            return row.element(by.name('email')).getText().then(function (val) {
                return val === name;
            });
        }).first().click();
    }

    public async UpdateUser(user: UpdateUser): Promise<boolean> {
        const basePage = new BasePage();
        switch (user.description) {
            case "update user's fullname":
                await this.txtFullName.clear();
                await this.txtFullName.sendKeys(user.NewFullName);
                await basePage.ClickUpdate();
                break;
            default:
                return false;
        }
        return basePage.GetOutputMessage().then(value => user.validationMessage === value);
    }

    public async RegisterUser(user: RegisterUser): Promise<boolean> {
        let result = false;
        const basePage = new BasePage();
        const snp = new SideNavigationPage();
        await this.btnRegisterNewUser.click();
        await this.txtEmail.sendKeys(this.randomize + user.Email);
        await this.txtRole.sendKeys(user.Role);
        await this.selTenant.click();
        await element(by.name(user.Tenant+this.randomize)).click();
        await basePage.ClickRegister();
        if (await basePage.GetOutputMessage() === user.existsMessage) {
            await snp.NavigateToUsersPage();
            result = true;
        } else if (await basePage.GetOutputMessage() === user.validationMessage) {
            result = true;
        } else {
            result = false;
        }
        return result;
    }

    public async UpdateRegisterUser(user: UpdateRegisterUser): Promise<boolean> {
        const basePage = new BasePage();
        switch (user.description) {
            case "update registered user's fullname":
                await this.txtFullName.sendKeys(user.NewFullName);
                await basePage.ClickUpdate();
                break;
            default:
                return false;
        }
        return basePage.GetOutputMessage().then(value => user.validationMessage === value);
    }
  }
