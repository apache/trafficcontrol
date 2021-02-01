import { ElementFinder, ExpectedConditions, browser, by, element } from 'protractor'
import { async, delay } from 'q';
import { BasePage } from './BasePage.po';
import {SideNavigationPage} from '../PageObjects/SideNavigationPage.po';

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
    private config = require('../config');
    private randomize = this.config.randomize;
    
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