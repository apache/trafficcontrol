import { ElementFinder, browser, by, element, ExpectedConditions, protractor } from 'protractor';
import { async, delay } from 'q';
import { BasePage } from './BasePage.po';
import { fstatSync } from 'fs';



export class TopNavigationPage extends BasePage{
    
    private lnkToggleLeftNavigationView = element(by.id('menu_toggle'));
    private btnShot = element(by.css('div[title="Diff CDN Config Snapshot"]'));
    private selectCDN = element(by.name('selectFormDropdown'));
    private btnCancelSnapshot = element(by.buttonText('Cancel'));
    private btnPerformSnapshot = element(by.buttonText('Perform Snapshot'));
    private btnYesSnapshot = element(by.buttonText('Yes'));
    private btnQueueCDN = element(by.css('div[title="Queue CDN Server Updates"]'));
    private btnDBDump = element(by.css('div[title="DB Dump"]'));
    private btnChangeLog = element(by.css('div[title="Change Logs"]'));
    private lnkUser = element(by.id('headerUsername'));
    private mnuManageUserProfile = element(by.linkText('Manage User Profile'));
    private txtEmail = element(by.name('email'));
    private mnuLogout = element(by.xpath("//li[@ng-if='userLoaded']")).element(by.linkText('Logout'));
    private btnLogout = element(by.xpath("//a[@uib-popover='Logout']"));
    private lnkAllUserPage = element(by.linkText('Users'));
    private bxLoginContainer = element(by.id("loginContainer"));
    private txtUserName = element(by.id("loginUsername"))
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
        const fs = require('fs');
        const folder = 'Downloads';
        await this.btnDBDump.click();
        await browser.wait(async function(){
            await fs.readdirSync(folder).forEach(file => {
                if (file != readme){
                    filename = file;
                }
            });
        }, 30*1000, 'File has not downloaded within 30 seconds').catch(function(){
            if(fs.existsSync(`Downloads/${filename}`))
        {
            //if file exist result will be true
            result = true;
            //delete the file
            fs.unlink(`Downloads/${filename}`, (err) => {
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
        if(await browser.wait(ExpectedConditions.visibilityOf(this.txtUserName), 20000) == true){
            result = true;
        }
        return result;
    }
}