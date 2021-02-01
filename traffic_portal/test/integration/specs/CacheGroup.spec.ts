import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po'
import { CacheGroupPage } from '../PageObjects/CacheGroup.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';


let fs = require('fs')
let using = require('jasmine-data-provider');

let filename = 'Data/CacheGroup/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let cacheGroupPage = new CacheGroupPage();

using(testData.CacheGroup, function (cacheGroupData) {
    describe('Traffic Portal - CacheGroup - ' + cacheGroupData.TestName, function () {
        using(cacheGroupData.Login, function (login) {
            it('can login', async function () {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open cache group page', async function () {
                await cacheGroupPage.OpenTopologyMenu();
                await cacheGroupPage.OpenCacheGroupsPage();
            })
            using(cacheGroupData.Create, function (create) {
                it(create.description, async function () {
                    expect(await cacheGroupPage.CreateCacheGroups(create, create.validationMessage)).toBeTruthy();
                    await cacheGroupPage.OpenCacheGroupsPage();
                })
            })
            using(cacheGroupData.Update, function (update) {
                if(update.description.includes("cannot")){
                    it(update.description, async function () {
                        await cacheGroupPage.SearchCacheGroups(update.Name)
                        expect(await cacheGroupPage.UpdateCacheGroups(update, update.validationMessage)).toBeUndefined();
                        await cacheGroupPage.OpenCacheGroupsPage();
                    }) 
                }else{
                    it(update.description, async function () {
                        await cacheGroupPage.SearchCacheGroups(update.Name)
                        expect(await cacheGroupPage.UpdateCacheGroups(update, update.validationMessage)).toBeTruthy();
                        await cacheGroupPage.OpenCacheGroupsPage();
                    }) 
                }
                
            })
            using(cacheGroupData.Remove, function (remove) {
                it(remove.description, async function () {
                    await cacheGroupPage.SearchCacheGroups(remove.Name)
                    expect(await cacheGroupPage.DeleteCacheGroups(remove.Name, remove.validationMessage)).toBeTruthy();
                }) 
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})