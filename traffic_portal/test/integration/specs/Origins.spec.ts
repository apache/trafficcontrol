import { browser } from 'protractor'
import { LoginPage } from '../PageObjects/LoginPage.po'
import { OriginsPage } from '../PageObjects/OriginsPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';

let fs = require('fs')
let using = require('jasmine-data-provider');

let setupFile = 'Data/Origins/Setup.json';
let cleanupFile = 'Data/Origins/Cleanup.json';
let filename = 'Data/Origins/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let originsPage = new OriginsPage();

describe('Setup Origin Delivery Service', function () {
    it('Setup', async function () {
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})
using(testData.Origins, async function (originsData) {
    using(originsData.Login, function (login) {
        describe('Traffic Portal - Origins - ' + login.description, function () {
            it('can login', async function () {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open origins page', async function () {
                await originsPage.OpenConfigureMenu();
                await originsPage.OpenOriginsPage();
            })
            using(originsData.Add, function (add) {
                it(add.description, async function () {
                    expect(await originsPage.CreateOrigins(add)).toBeTruthy();
                    await originsPage.OpenOriginsPage();
                })
            })
            using(originsData.Update, function (update) {
                if (update.validationMessage == undefined) {
                    it(update.description, async function () {
                        await originsPage.SearchOrigins(update.Name);
                        expect(await originsPage.UpdateOrigins(update)).toBeUndefined();
                        await originsPage.OpenOriginsPage();
                    })
                } else {
                    it(update.description, async function () {
                        await originsPage.SearchOrigins(update.Name);
                        expect(await originsPage.UpdateOrigins(update)).toBeTruthy();
                        await originsPage.OpenOriginsPage();
                    })
                }
            })
            using(originsData.Remove, function (remove) {
                it(remove.description, async function () {
                    await originsPage.SearchOrigins(remove.Name);
                    expect(await originsPage.DeleteOrigins(remove)).toBeTruthy();
                    await originsPage.OpenOriginsPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean up Origin Delivery Service', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})