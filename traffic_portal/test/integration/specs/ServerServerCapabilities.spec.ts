import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po'
import { ServerCapabilitiesPage } from '../PageObjects/ServerCapabilitiesPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { ServersPage } from '../PageObjects/ServersPage.po';
import { API } from '../CommonUtils/API';

let fs = require('fs')
let using = require('jasmine-data-provider');

let setupFile = 'Data/ServerServerCapabilities/Setup.json';
let cleanupFile = 'Data/ServerServerCapabilities/Cleanup.json';
let filename = 'Data/ServerServerCapabilities/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let serverCapabilitiesPage = new ServerCapabilitiesPage();
let serverPage = new  ServersPage();

describe("Setup Server Capabilities and Server for prereq", function(){
    it('Setup', async function(){
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})
using(testData.ServerServerCapabilities, async function(serverServerCapData){
    using(serverServerCapData.Login, function(login){
        describe('Traffic Portal - Server Server Capabilities - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open server page', async function(){
                await serverPage.OpenConfigureMenu();
                await serverPage.OpenServerPage();
            })
            using(serverServerCapData.Link, function(link){
                if(link.description.includes("cannot")){
                    it(link.description, async function(){
                        await serverPage.SearchServer(link.Server);
                        expect(await serverPage.AddServerCapabilitiesToServer(link)).toBeUndefined();
                        await serverPage.OpenServerPage();
                    })
                }else{
                    it(link.description, async function(){
                        await serverPage.SearchServer(link.Server);
                        expect(await serverPage.AddServerCapabilitiesToServer(link)).toBeTruthy();
                        await serverPage.OpenServerPage();
                    })
                }
            })
            using(serverServerCapData.Remove, function(remove){
                it(remove.description, async function(){
                    await serverPage.SearchServer(remove.Server);
                    expect(await serverPage.RemoveServerCapabilitiesFromServer(remove.ServerCapability, remove.validationMessage)).toBeTruthy();
                    await serverPage.OpenServerPage();
                })
            })
            it('can open server capabilities page', async function(){
                await serverCapabilitiesPage.OpenServerCapabilityPage();
            })
            using(serverServerCapData.DeleteServerCapability, function(deleteSC){
                it(deleteSC.description, async function(){
                    await serverCapabilitiesPage.SearchServerCapabilities(deleteSC.ServerCapability);
                    expect(await serverCapabilitiesPage.DeleteServerCapabilities(deleteSC.ServerCapability, deleteSC.validationMessage)).toBeTruthy();
                    await serverCapabilitiesPage.OpenServerCapabilityPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})
describe("Clean up prereq", function(){
    it('Clean up', async function(){
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})
