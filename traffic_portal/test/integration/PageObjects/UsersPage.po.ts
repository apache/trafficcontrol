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
import {SideNavigationPage} from '../PageObjects/SideNavigationPage.po';
import { config, randomize } from '../config';

export class UsersPage extends BasePage {

    private btnCreateNewUser = element(by.css('[title="Create New User"]'));
    private txtFullName = element(by.name('fullName'));
    private txtUserName = element(by.name('uName'));
    private txtEmail = element(by.name('email'));
    private txtRole = element(by.name('role'));
    private txtTenant = element(by.name('tenantId'));
    private txtPassword = element(by.name('uPass'));
    private txtConfirmPassword = element(by.name('confirmPassword'));
    private txtPublicSSHKey = element(by.name('publicSshKey'));
    private readonly config = config;
    private randomize = randomize;

    async OpenUserPage(){
      let snp = new SideNavigationPage();
      await snp.ClickUserAdminMenu();
      await snp.NavigateToUsersPage();
     }

    async CreateUser(user) {
      let result = false;
      let basePage = new BasePage();
      let snp = new SideNavigationPage();
      await this.btnCreateNewUser.click();
      await this.txtFullName.sendKeys(user.FullName + this.randomize);
      await this.txtUserName.sendKeys(user.Username + this.randomize);
      await this.txtEmail.sendKeys(user.FullName + this.randomize + user.Email);
      await this.txtRole.sendKeys(user.Role);
      await this.txtTenant.sendKeys(user.Tenant+this.randomize);
      await this.txtPassword.sendKeys(user.Password);
      await this.txtConfirmPassword.sendKeys(user.ConfirmPassword);
      await this.txtPublicSSHKey.sendKeys(user.PublicSSHKey);
      await basePage.ClickCreate();
      if(await basePage.GetOutputMessage() == user.existsMessage){
        await snp.NavigateToUsersPage();
        result = true;
      }else if(await basePage.GetOutputMessage() == user.validationMessage){
        result = true;
      }else{
        result = false;
      }
      return result;
    }

  }
