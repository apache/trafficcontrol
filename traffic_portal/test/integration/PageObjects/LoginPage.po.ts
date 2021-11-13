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
import { browser, by, element} from 'protractor';

import { randomize } from "../config";
import { BasePage } from './BasePage.po'

/**
 * LoginData is all the data needed to authenticate with Traffic Ops (and some
 * that isn't).
 */
interface LoginData {
    /** Optional human-readable description for the login. This is not used.  */
	description?: string;
    /** The password used for authentication. */
	password: string;
    /** The username of the user as whom to authenticate. */
	username: string;
    /**
     * If present, the content of this string is matched against the alert that
     * is showing. A value of `undefined` indicates that there should be no
     * alert.
     */
	validationMessage?: string;
}

export class LoginPage extends BasePage{
    private txtUserName = element(by.id("loginUsername"))
    private txtPassword = element(by.id("loginPass"))
    private btnLogin = element(by.name("loginSubmit"))
    private lnkResetPassword= element (by.xpath("//button[text()='Reset Password']"))
    private lblUserName = element(by.xpath("//span[@id='headerUsername']"))
    private bannerEnvironment = element(by.css('.enviro-banner.prod'));
    private randomize = randomize;


    public async Login(login: LoginData){
        let result = false;
        const basePage = new BasePage();
        if(login.username === 'admin'){
            await this.txtUserName.sendKeys(login.username)
            await this.txtPassword.sendKeys(login.password)
            await browser.actions().mouseMove(this.btnLogin).perform();
            await browser.actions().click(this.btnLogin).perform();
        }else{
            await this.txtUserName.sendKeys(login.username+this.randomize)
            await this.txtPassword.sendKeys(login.password)
            await browser.actions().mouseMove(this.btnLogin).perform();
            await browser.actions().click(this.btnLogin).perform();
        }
        if(await browser.getCurrentUrl() === browser.params.baseUrl + "/#!/login"){
            result = await basePage.GetOutputMessage().then(value => value === login.validationMessage);
        }else{
            result = true;
        }
        return result;
    }

    public async ClickResetPassword(): Promise<void> {
        this.lnkResetPassword.click()
    }

    public async CheckUserName(login: LoginData): Promise<boolean> {
        if(await this.lblUserName.getText() === 'admin' || await this.lblUserName.getText() === login.username+this.randomize){
            return true;
        }else{
            return false;
        }
    }

    public async CheckBanner(): Promise<boolean> {
        if(await this.bannerEnvironment.isPresent()){
            return true;
        }else{
            return false;
        }
    }
};
