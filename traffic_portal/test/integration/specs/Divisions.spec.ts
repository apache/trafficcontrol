import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';
import { DivisionsPage } from '../PageObjects/Divisions.po';

let fs = require('fs')
let using = require('jasmine-data-provider');

let setupFile = 'Data/Divisions/Setup.json';
let cleanupFile = 'Data/Divisions/Cleanup.json';
let filename = 'Data/Divisions/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let divisionsPage = new DivisionsPage();

describe('Setup API for Divisions Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.Divisions, async function(divisionsData){
    using(divisionsData.Login, function(login){
        describe('Traffic Portal - Divisions - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open divisions page', async function(){
                await divisionsPage.OpenTopologyMenu();
                await divisionsPage.OpenDivisionsPage();
            })

            using(divisionsData.Add, function (add) {
                it(add.description, async function () {
                    expect(await divisionsPage.CreateDivisions(add)).toBeTruthy();
                    await divisionsPage.OpenDivisionsPage();
                })
            })
            using(divisionsData.Update, function (update) {
                it(update.description, async function () {
                    await divisionsPage.SearchDivisions(update.Name);
                    expect(await divisionsPage.UpdateDivisions(update)).toBeTruthy();
                    await divisionsPage.OpenDivisionsPage();
                })
            })
            using(divisionsData.Remove, function (remove) {
                it(remove.description, async function () {
                    await divisionsPage.SearchDivisions(remove.Name);
                    expect(await divisionsPage.DeleteDivisions(remove)).toBeTruthy();
                    await divisionsPage.OpenDivisionsPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean Up API for Divisions Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})