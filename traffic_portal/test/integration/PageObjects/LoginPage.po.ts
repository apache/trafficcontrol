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

export class LoginPage extends BasePage {
    private readonly txtUserName = element(by.id("loginUsername"))
    private readonly txtPassword = element(by.id("loginPass"))
    private readonly btnLogin = element(by.name("loginSubmit"))
    private readonly lnkResetPassword= element(by.buttonText("Reset Password"))
    private readonly lblUserName = element(by.id("headerUsername"))
    private readonly bannerEnvironment = element(by.className('enviro-banner.prod'));


    public async Login(login: LoginData){
        let username = login.username;
        if (login.username !== this.login.username) {
            username += this.randomize;
        }

        await this.txtUserName.sendKeys(username)
        await this.txtPassword.sendKeys(login.password)
        await browser.actions().mouseMove(this.btnLogin).perform();
        await browser.actions().click(this.btnLogin).perform();

        const val = await this.GetOutputMessage();
        if(await browser.getCurrentUrl() === browser.params.baseUrl + "/#!/login"){
            return val === login.validationMessage;
        }
        return true;
    }

    public async ClickResetPassword(): Promise<void> {
        this.lnkResetPassword.click()
    }

    public async CheckUserName(login: LoginData): Promise<boolean> {
        const txt = await this.lblUserName.getText();
        return txt === 'admin' || txt === login.username+this.randomize;
    }

    public async CheckBanner(): Promise<boolean> {
        return this.bannerEnvironment.isPresent();
    }
};
