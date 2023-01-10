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
import { ExpectedConditions, browser, by, element } from 'protractor';

import { BasePage } from './BasePage.po';
import {SideNavigationPage} from '../PageObjects/SideNavigationPage.po';
import { randomize } from '../config';

export class ServerCapabilitiesPage extends BasePage{

     private btnCreateServerCapabilities = element(by.name('createServerCapabilityButton'));
     private txtSCName = element(by.id("name"))
     private txtSCDescription = element(by.id("description"))
     private searchFilter = element(by.id('serverCapabilitiesTable_filter')).element(by.css('label input'));
     private btnDelete = element(by.buttonText('Delete'))
     private txtConfirmCapabilitiesName = element(by.name('confirmWithNameInput'));
     private randomize = randomize;


     public async OpenServerCapabilityPage(){
      let snp = new SideNavigationPage();
      await snp.NavigateToServerCapabilitiesPage();
     }
     public async OpenConfigureMenu(){
      let snp = new SideNavigationPage();
      await snp.ClickConfigureMenu();
     }

      public async CreateServerCapabilities(nameSC: string, descriptionSC: string, outputMessage:string){
        let result = false
        let basePage = new BasePage();
        let snp= new SideNavigationPage();
        let name = nameSC+this.randomize;
        await this.btnCreateServerCapabilities.click();
        if(name != this.randomize){
          await this.txtSCName.sendKeys(name);
        }
        await this.txtSCDescription.sendKeys(descriptionSC);
        if(outputMessage == await(basePage.GetBlankErrorMessage()) || outputMessage == await(basePage.GetSyntaxErrorMessage())) {
          await snp.NavigateToServerCapabilitiesPage();
          result = true;
        }else{
          await basePage.ClickCreate();
          if ((await basePage.GetOutputMessage()) == outputMessage) {
            result = true;
          }else if((await basePage.GetOutputMessage()).includes(' already exists.') || (await basePage.GetOutputMessage()).includes('Forbidden')){
            await snp.NavigateToServerCapabilitiesPage();
            result = true;
          }else{
            result = false;
          }
        }
        return result;
      }


    public async SearchServerCapabilities(nameSC:string){
      let name = nameSC+this.randomize;
      await this.searchFilter.clear();
      await this.searchFilter.sendKeys(name);
      await element.all(by.repeater('sc in ::serverCapabilities')).filter(function(row){
        return row.element(by.name('name')).getText().then(function(val){
          return val === name;
        });
      }).first().click();
    }

     public async DeleteServerCapabilities(nameSC:string, outputMessage:string){
      let result = false;
      let basePage = new BasePage();
      let name = nameSC+this.randomize;
      await this.btnDelete.click();
      await browser.wait(ExpectedConditions.visibilityOf(this.txtConfirmCapabilitiesName), 1000);
      await this.txtConfirmCapabilitiesName.sendKeys(name);
      if(await basePage.ClickDeletePermanently() == true){
        result = await basePage.GetOutputMessage().then(function(value){
          if(outputMessage == value){
            return true;
          }else{
            return false;
          }
        })
      }else{
        await basePage.ClickCancel();
      }
      return result;
     }

     public async CheckCSV(name: string): Promise<boolean> {
      return element(by.cssContainingText("span", name)).isPresent();
  }
}
