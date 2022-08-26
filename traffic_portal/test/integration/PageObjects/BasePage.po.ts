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
import { browser, element, by, ExpectedConditions } from 'protractor';
/**
 * Class representing generic page.
 * Methods/properties for global elements should go here.
 *
 * @export
 * @class BasePage
 */
export class BasePage {
  constructor() {}

  protected lblPageTitle = element(by.xpath("//ol[@class='breadcrumb pull-left']//li[@class='active']"))
  protected lblSubPageTitle = element(by.xpath("//li[@class='active ng-binding']"))
  private lblPopupPageTitle = element(by.xpath("//h4[@class='modal-title ng-binding']"));
  private lbOutputMessage = element((by.xpath("(//div[contains(@class,'alert alert-dismissable')]//div)[1]")));
  private lbOutputWarning = element(by.xpath("(//div[contains(@class,'alert alert-dismissable alert-warning')]//div)[1]"));
  private lbBlankError = element(by.xpath("//small[text()='Required']"));
  private lbSyntaxError = element(by.xpath("//small[text()='Must be alphamumeric with no spaces. Dashes and underscores also allowed.']"));
  private btnCreate= element(by.buttonText('Create'));
  private btnDeletePermanently = element(by.buttonText('Delete Permanently'));
  private btnCancel =  element(by.className('close')).element(by.xpath("//span[text()='×']"));
  private btnUpdate = element(by.buttonText('Update'));
  private btnSubmit = element(by.buttonText("Submit"));
  private btnRegister = element(by.buttonText('Send Registration'));
  private btnNo = element(by.buttonText("No"));

  async ClickNo(){
    await this.btnNo.click();
  }
  async ClickSubmit(){
    if(await this.btnSubmit.isEnabled() == true){
      await this.btnSubmit.click();
      return true;
    }else{
      return false;
    }

  }
  public async ClickUpdate(): Promise<boolean>{
    if(await this.btnUpdate.isEnabled()){
      await this.btnUpdate.click();
      return true;
    }else{
      return false;
    }
  }
  async ClickDeletePermanently(){
    if(await this.btnDeletePermanently.isEnabled() == true){
      await this.btnDeletePermanently.click();
      return true;
    }else{
      return false;
    }
  }
  public async ClickCreate(): Promise<boolean> {
    if(await this.btnCreate.isEnabled()){
      await this.btnCreate.click();
      return true;
    }else{
      return false;
    }
  }
  public async ClickRegister(): Promise<boolean> {
    if(await this.btnRegister.isEnabled()){
      await this.btnRegister.click();
      return true;
    }else{
      return false;
    }
  }
  async ClickCancel(){
    await this.btnCancel.click();
  }

  GetOutputMessage(){
    browser.wait(ExpectedConditions.visibilityOf(this.lbOutputMessage), 2000);
    return this.lbOutputMessage.getText();
  }
  GetOutputWarning(){
    browser.wait(ExpectedConditions.visibilityOf(this.lbOutputWarning), 2000);
    return this.lbOutputWarning.getText();
  }
  GetBlankErrorMessage(){
    return this.lbBlankError.getText();
  }
  GetSyntaxErrorMessage(){
      return this.lbSyntaxError.getText();
  }
  GetPageTitle(){
    return this.lblPageTitle.getText()
  }
  GetSubPageTitle(){
    return this.lblSubPageTitle.getText()
  }
  GetPopupTitle(){
    return this.lblPopupPageTitle.getText();
  }
}
