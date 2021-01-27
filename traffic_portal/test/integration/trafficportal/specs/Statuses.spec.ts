import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';
import { StatusesPage } from '../PageObjects/Statuses.po'

let fs = require('fs')
let using = require('jasmine-data-provider');

let setupFile = 'Data/Statuses/Setup.json';
let cleanupFile = 'Data/Statuses/Cleanup.json';
let filename = 'Data/Statuses/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let statusesPage = new StatusesPage();

describe('Setup API for Statuses Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})
using(testData.Statuses, async function(statusesData){
    using(statusesData.Login, function(login){
        describe('Traffic Portal - Statuses - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open statuses page', async function(){
                await statusesPage.OpenConfigureMenu();
                await statusesPage.OpenStatusesPage();
            })

            using(statusesData.Add, function (add) {
                it(add.description, async function () {
                    expect(await statusesPage.CreateStatus(add)).toBeTruthy();
                    await statusesPage.OpenStatusesPage();
                })
            })
            using(statusesData.Update, function (update) {
                it(update.description, async function () {
                    await statusesPage.SearchStatus(update.Name);
                    expect(await statusesPage.UpdateStatus(update)).toBeTruthy();
                    await statusesPage.OpenStatusesPage();
                })
            })
            using(statusesData.Remove, function (remove) {
                it(remove.description, async function () {
                    await statusesPage.SearchStatus(remove.Name);
                    expect(await statusesPage.DeleteStatus(remove)).toBeTruthy();
                    await statusesPage.OpenStatusesPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})
describe('Clean Up API for Statuses Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})