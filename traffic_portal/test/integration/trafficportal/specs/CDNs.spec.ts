import { browser } from 'protractor'
import { LoginPage } from '../PageObjects/LoginPage.po'
import { CDNPage } from '../PageObjects/CDNPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';


let fs = require('fs')
let using = require('jasmine-data-provider');

let filename = 'Data/CDN/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let cdnsPage = new CDNPage();

using(testData.CDN, async function(cdnsData){
    using(cdnsData.Login, function(login){
        describe('Traffic Portal - CDN - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open CDN page', async function(){
                await cdnsPage.OpenCDNsPage();
            })

            using(cdnsData.Add, function (add){
                it(add.description, async function(){
                    expect(await cdnsPage.CreateCDN(add)).toBeTruthy();
                    await cdnsPage.OpenCDNsPage();
                })
            })
            using(cdnsData.Update, function(update){
                it(update.description, async function(){
                    await cdnsPage.SearchCDN(update.Name);
                    expect(await cdnsPage.UpdateCDN(update)).toBeTruthy();
                })
            
            })
            using(cdnsData.Remove, function(remove){
                it(remove.description, async function(){
                    await cdnsPage.SearchCDN(remove.Name);
                    expect(await cdnsPage.DeleteCDN(remove)).toBeTruthy();

                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })

        })
    })
})