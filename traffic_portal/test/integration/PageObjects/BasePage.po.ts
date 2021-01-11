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
  private btnCreate= element(by.xpath("//button[text()='Create']"));
  private btnDeletePermanently = element(by.buttonText('Delete Permanently'));
  private btnCancel =  element(by.className('close')).element(by.xpath("//span[text()='×']"));
  private btnUpdate = element(by.xpath("//button[text()='Update']"))
  private btnSubmit = element(by.xpath("//button[text()='Submit']"));
  private btnNo = element(by.xpath("//button[text()='No']"));
  
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
  async ClickUpdate(){
    if(await this.btnUpdate.isEnabled() == true){
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
  async ClickCreate(){
    if(await this.btnCreate.isEnabled() == true){
      await this.btnCreate.click();
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