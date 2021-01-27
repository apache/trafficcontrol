import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po'
import { ServerCapabilitiesPage } from '../PageObjects/ServerCapabilitiesPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';

let fs = require('fs')
let using = require('jasmine-data-provider');

let filename = 'Data/ServerCapabilities/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let serverCapabilitiesPage = new ServerCapabilitiesPage();

using(testData.ServerCapabilities, function(serverCapabilitiesData) {
    describe('Traffic Portal - Server Capabilities - '+ serverCapabilitiesData.TestName,  function(){
        using(serverCapabilitiesData.Login, function(login) {
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open server capability page', async function() {
                await serverCapabilitiesPage.OpenConfigureMenu();
                await serverCapabilitiesPage.OpenServerCapabilityPage();
            })
            using(serverCapabilitiesData.Add, function(add) {
                it(add.description, async function(){
                    expect(await serverCapabilitiesPage.CreateServerCapabilities(add.Name, add.validationMessage)).toBeTruthy();
                    await serverCapabilitiesPage.OpenServerCapabilityPage();
                })
            })
            using(serverCapabilitiesData.Delete, function(remove) {
                if(remove.description.includes("invalid")){
                    it(remove.description, async function(){
                        await serverCapabilitiesPage.SearchServerCapabilities(remove.Name)
                        expect(await serverCapabilitiesPage.DeleteServerCapabilities(remove.InvalidName, remove.validationMessage)).toBeFalsy();  
                        await serverCapabilitiesPage.OpenServerCapabilityPage();
                    })
                } else {
                    it(remove.description, async function(){
                        await serverCapabilitiesPage.SearchServerCapabilities(remove.Name)
                        expect(await serverCapabilitiesPage.DeleteServerCapabilities(remove.Name, remove.validationMessage)).toBeTruthy();
                        await serverCapabilitiesPage.OpenServerCapabilityPage();
                    })
                }
            })
            it('can logout', async function(){
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})