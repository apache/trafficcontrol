import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po'
import { ServersPage } from '../PageObjects/ServersPage.po';
import { API } from '../CommonUtils/API';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';

let fs = require('fs')
let using = require('jasmine-data-provider');

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let serversPage = new ServersPage();

let setupFile = 'Data/Servers/Setup.json';
let cleanupFile = 'Data/Servers/Cleanup.json';
let filename = 'Data/Servers/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

describe('Setup API call for Servers Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.Servers, async function(serversData){
    using(serversData.Login, function(login){
        describe('Traffic Portal - Servers - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open servers page', async function(){
                await serversPage.OpenConfigureMenu();
                await serversPage.OpenServerPage();
            })
            using(serversData.Add, function (add) {
                it(add.description, async function () {
                    expect(await serversPage.CreateServer(add)).toBeTruthy();
                    await serversPage.OpenServerPage();
                })
            })
            using(serversData.Update, function (update) {
                it(update.description, async function () {
                    await serversPage.SearchServer(update.Name);
                    expect(await serversPage.UpdateServer(update)).toBeTruthy();
                    await serversPage.OpenServerPage();
                })
            })
            using(serversData.Remove, function (remove) {
                it(remove.description, async function () {
                    await serversPage.SearchServer(remove.Name);
                    expect(await serversPage.DeleteServer(remove)).toBeTruthy();
                    await serversPage.OpenServerPage();
                })
            })
            it('can logout', async function(){
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('API Clean Up for Servers Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})