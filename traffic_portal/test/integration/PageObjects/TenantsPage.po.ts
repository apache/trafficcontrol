import { ElementFinder, browser, by, element } from 'protractor'
import { async, delay } from 'q';
import { BasePage } from './BasePage.po';
import {SideNavigationPage} from './SideNavigationPage.po';
export class TenantsPage extends BasePage {
  
    private btnCreateNewTenant = element(by.xpath("//button[@title='Create New Tenant']"));
    private txtName = element(by.name('name'));
    private txtActive = element(by.name('active'));
    private txtParentTenant = element(by.name('parentId'));
    private txtSearch = element(by.id('tenantsTable_filter')).element(by.css('label input'));
    private mnuTenantTable = element(by.id('tenantsTable'));
    private btnDelete = element(by.buttonText('Delete'));
    private txtConfirmTenantName = element(by.name('confirmWithNameInput'));
    private config = require('../config');
    private randomize = this.config.randomize;

    async OpenTenantPage(){
      let snp = new SideNavigationPage();
      await snp.ClickUserAdminMenu();
      await snp.NavigateToTenantsPage();
    }
    async CreateTenant(tenant){
        let result = false;
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        await this.btnCreateNewTenant.click();
        await this.txtName.sendKeys(tenant.Name+this.randomize);
        await this.txtActive.sendKeys(tenant.Active);
        if(tenant.ParentTenant == '- root'){
          await this.txtParentTenant.sendKeys(tenant.ParentTenant);
        }else{
          await this.txtParentTenant.sendKeys(tenant.ParentTenant+this.randomize);   
        }
        await basePage.ClickCreate();
        if(await basePage.GetOutputMessage() == tenant.existsMessage){
          await snp.NavigateToTenantsPage();
          result = true;
        }else if(await basePage.GetOutputMessage() == tenant.validationMessage){
          result = true;
        }else{
          result = false;
        }
        return result;
    }
    async SearchTenant(name:string){
        let result = false;
        let snp = new SideNavigationPage();
        await snp.NavigateToTenantsPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        await element.all(by.repeater('t in ::tenants')).filter(function(row){
            return row.element(by.name('name')).getText().then(function(val){
              return val === name;
            });
          }).first().click();
    }
    async DeleteTenant(name:string,outputMessage:string){
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmTenantName.sendKeys(name);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function(value){
            if(outputMessage == value){
              return true;
            }else{
              return false;
            }
          })
          return result;
    }
}