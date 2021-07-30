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
import { randomize } from "../config";
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

interface DeleteParameter {
  Name: string;
  validationMessage?: string;
}

interface Parameter extends DeleteParameter {
  ConfigFile: string;
  Secure: string;
  Value: string;
}

interface UpdateParameter {
  description: string;
  ConfigFile: string;
  validationMessage?: string;
}

export class ParametersPage extends BasePage {

  private btnCreateNewParameter = element(by.name('createParameterButton'));
  private txtName = element(by.name('name'));
  private txtConfigFile = element(by.name('configFile'));
  private txtValue = element((by.name("value")));
  private txtSecure = element(by.name('secure'));
  private txtSearch = element(by.id('parametersTable_filter')).element(by.css('label input'));
  private btnDelete = element(by.buttonText('Delete'));
  private btnYes = element(by.buttonText('Yes'));
  private txtConfirmName = element(by.name('confirmWithNameInput'));
  private btnTableColumn = element(by.className("caret"))
  private randomize = randomize;

  public async OpenParametersPage() {
    const snp = new SideNavigationPage();
    await snp.NavigateToParametersPage();
  }

  public async OpenConfigureMenu() {
    const snp = new SideNavigationPage();
    await snp.ClickConfigureMenu();
  }

  public async CreateParameter(parameter: Parameter): Promise<boolean> {
    let result = false;
    const basePage = new BasePage();
    await this.btnCreateNewParameter.click();
    await this.txtName.sendKeys(parameter.Name + this.randomize);
    await this.txtConfigFile.sendKeys(parameter.ConfigFile);
    await this.txtValue.sendKeys(parameter.Value)
    await this.txtSecure.sendKeys(parameter.Secure);
    await basePage.ClickCreate();
    result = await basePage.GetOutputMessage().then(function (value) {
      if (parameter.validationMessage === value) {
        return true;
      } else {
        return false;
      }
    })
    return result;
  }

  public async SearchParameter(nameParameter: string) {
    let name = nameParameter + this.randomize;
    await this.txtSearch.clear();
    await this.txtSearch.sendKeys(name);
    await element.all(by.repeater('p in ::parameters')).filter(function (row) {
      return row.element(by.name('name')).getText().then(function (val) {
        return val === name;
      });
    }).first().click();
  }

  public async UpdateParameter(parameter: UpdateParameter): Promise<boolean> {
    const basePage = new BasePage();
    switch (parameter.description) {
      case "update parameter configfile":
        await this.txtConfigFile.clear();
        await this.txtConfigFile.sendKeys(parameter.ConfigFile);
        await basePage.ClickUpdate();
        await this.btnYes.click();
        break;
      default:
        return false;
    }
    return await basePage.GetOutputMessage().then(value => parameter.validationMessage === value);
  }

  async DeleteParameter(parameter: DeleteParameter): Promise<boolean> {
    let result = false;
    const basePage = new BasePage();
    await this.btnDelete.click();
    await this.btnYes.click();
    await this.txtConfirmName.sendKeys(parameter.Name + this.randomize);
    await basePage.ClickDeletePermanently();
    result = await basePage.GetOutputMessage().then(function (value) {
      if (parameter.validationMessage === value) {
        return true;
      } else {
        return false;
      }
    })
    return result;
  }

  public async CheckCSV(name: string): Promise<boolean> {
    return element(by.cssContainingText("span", name)).isPresent();
  }
  
  public async ToggleTableColumn(name: string): Promise<boolean> {
    await this.btnTableColumn.click();
    const result = await element(by.cssContainingText("th", name)).isPresent();
    await element(by.cssContainingText("label", name)).click();
    await this.btnTableColumn.click();
    return !result;
  }
}
