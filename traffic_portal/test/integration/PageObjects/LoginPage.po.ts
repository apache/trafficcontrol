import { ElementFinder, browser, by, element, ExpectedConditions } from 'protractor'
import { BasePage } from './BasePage.po'
import { timeout } from 'q'
export class LoginPage extends BasePage{
    private txtUserName = element(by.id("loginUsername"))
    private txtPassword = element(by.id("loginPass"))
    private btnLogin = element(by.name("loginSubmit"))
    private lnkResetPassword= element (by.xpath("//button[text()='Reset Password']"))
    private lblUserName = element(by.xpath("//span[@id='headerUsername']"))

    private config = require('../config');
    private randomize = this.config.randomize;

    async Login(userName: string, password: string ){
        if(userName == 'admin'){
            await this.txtUserName.sendKeys(userName)
            await this.txtPassword.sendKeys(password)
            await browser.actions().mouseMove(this.btnLogin).perform();
            await browser.actions().click(this.btnLogin).perform();    
        }else{
            await this.txtUserName.sendKeys(userName+this.randomize)
            await this.txtPassword.sendKeys(password)
            await browser.actions().mouseMove(this.btnLogin).perform();
            await browser.actions().click(this.btnLogin).perform();    
        }
    }
    ClickResetPassword(){
        this.lnkResetPassword.click()
    }
    async CheckUserName(userName: string) {
        if(await this.lblUserName.getText() == 'admin' || await this.lblUserName.getText() == userName+this.randomize){
            return true;
        }else{
            return false;   
        }
    }
};
