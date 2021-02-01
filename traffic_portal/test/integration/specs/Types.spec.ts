import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';
import { TypesPage } from '../PageObjects/Types.po'

let fs = require('fs')
let using = require('jasmine-data-provider');

let setupFile = 'Data/Types/Setup.json';
let cleanupFile = 'Data/Types/Cleanup.json';
let filename = 'Data/Types/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let typesPage = new TypesPage();

describe('Setup API for Types Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})
using(testData.Types, async function(typesData){
    using(typesData.Login, function(login){
        describe('Traffic Portal - Types - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open types page', async function(){
                await typesPage.OpenConfigureMenu();
                await typesPage.OpenTypesPage();
            })

            using(typesData.Add, function (add) {
                it(add.description, async function () {
                    expect(await typesPage.CreateType(add)).toBeTruthy();
                    await typesPage.OpenTypesPage();
                })
            })
            using(typesData.Update, function (update) {
                it(update.description, async function () {
                    await typesPage.SearchType(update.Name);
                    expect(await typesPage.UpdateType(update)).toBeTruthy();
                    await typesPage.OpenTypesPage();
                })
            })
            using(typesData.Remove, function (remove) {
                it(remove.description, async function () {
                    await typesPage.SearchType(remove.Name);
                    expect(await typesPage.DeleteTypes(remove)).toBeTruthy();
                    await typesPage.OpenTypesPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})
describe('Clean Up API for Types Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})