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
import { existsSync, readdirSync, unlink } from "fs";

import { browser, by, element, until } from 'protractor';
import { BasePage } from './BasePage.po';

export class TopNavigationPage extends BasePage{

    private btnShot = element(by.css('div[title="Diff CDN Config Snapshot"]'));
    private selectCDN = element(by.name('selectFormDropdown'));
    private btnPerformSnapshot = element(by.buttonText('Perform Snapshot'));
    private btnYesSnapshot = element(by.buttonText('Yes'));
    private btnQueueCDN = element(by.css('div[title="Queue CDN Server Updates"]'));
    private btnDBDump = element(by.css('div[title="DB Dump"]'));
    private lnkUser = element(by.id('headerUsername'));
    private mnuManageUserProfile = element(by.linkText('Manage User Profile'));
    private btnLogout = element(by.xpath("//a[@uib-popover='Logout']"));
    async PerformSnapshot(cdnname:string,message:string){
        let result = false;
        let basePage = new BasePage();
        await this.btnShot.click();
        await this.selectCDN.sendKeys(cdnname);
        await basePage.ClickSubmit();
        await this.btnPerformSnapshot.click();
        await this.btnYesSnapshot.click();
        await basePage.GetOutputMessage().then(function(value){
            if(message == value){
                result = true;
            }else{
                result = false;
            }
        })
        return result
    }

    async QueueServerUpdates(cdnname:string,message:string ){
        let result=false
        let basePage = new BasePage();
        await this.btnQueueCDN.click();
        await this.selectCDN.sendKeys(cdnname);
        await basePage.ClickSubmit();
        await basePage.GetOutputMessage().then(function(value){
            if(message == value){
                result = true;
            }else{
                result = false;
            }
        })
        return result;
    }

    async FileDownloaded(){
        let filename= "";
        let result = false;
        let readme = 'Readme.md';
        const folder = 'Downloads';
        await this.btnDBDump.click();
        await browser.wait(async function(){
            await readdirSync(folder).forEach(file => {
                if (file != readme){
                    filename = file;
                }
            });
        }, 30*1000, 'File has not downloaded within 30 seconds').catch(function(){
            if(existsSync(`Downloads/${filename}`))
        {
            //if file exist result will be true
            result = true;
            //delete the file
            unlink(`Downloads/${filename}`, (err) => {
                if (err) throw err;
            });
        }
        });
        return result;
    }

    async ManageUserProfile(username:string){
        let result = false;
        await this.lnkUser.click();
        await this.mnuManageUserProfile.click();
        await this.GetSubPageTitle().then(function(value){
            if(username == value){
                result = true;
            }else{
                result = false;
            }
        })
        return result;
    }

    async Logout(){
        let result = false;
        await this.btnLogout.click();
        if(await browser.wait(until.urlIs(browser.params.baseUrl + "/#!/login"), 10000) === true){
            result = true;
        }
        return result;
    }
}
