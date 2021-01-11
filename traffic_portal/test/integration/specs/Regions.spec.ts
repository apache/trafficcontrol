import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';
import { RegionsPage } from '../PageObjects/RegionsPage.po';

let fs = require('fs')
let using = require('jasmine-data-provider');

let setupFile = 'Data/Regions/Setup.json';
let cleanupFile = 'Data/Regions/Cleanup.json';
let filename = 'Data/Regions/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let regionsPage = new RegionsPage();

describe('Setup Divisions for Regions Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.Regions, async function(regionsData){
    using(regionsData.Login, function(login){
        describe('Traffic Portal - Regions - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open regions page', async function(){
                await regionsPage.OpenTopologyMenu();
                await regionsPage.OpenRegionsPage();
            })

            using(regionsData.Add, function (add) {
                it(add.description, async function () {
                    expect(await regionsPage.CreateRegions(add)).toBeTruthy();
                    await regionsPage.OpenRegionsPage();
                })
            })
            using(regionsData.Update, function (update) {
                it(update.description, async function () {
                    await regionsPage.SearchRegions(update.Name);
                    expect(await regionsPage.UpdateRegions(update)).toBeTruthy();
                    await regionsPage.OpenRegionsPage();
                })
            })
            using(regionsData.Remove, function (remove) {
                it(remove.description, async function () {
                    await regionsPage.SearchRegions(remove.Name);
                    expect(await regionsPage.DeleteRegions(remove)).toBeTruthy();
                    await regionsPage.OpenRegionsPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean Up Divisions for Regions Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})