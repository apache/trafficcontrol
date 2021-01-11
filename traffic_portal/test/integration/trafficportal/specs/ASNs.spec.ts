import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';
import { ASNsPage } from '../PageObjects/ASNs.po';

let fs = require('fs')
let using = require('jasmine-data-provider');

let setupFile = 'Data/ASNs/Setup.json';
let cleanupFile = 'Data/ASNs/Cleanup.json';
let filename = 'Data/ASNs/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let asnsPage = new ASNsPage();

describe('Setup API for ASNs Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.ASNs, async function(asnsData){
    using(asnsData.Login, function(login){
        describe('Traffic Portal - ASNs - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open asns page', async function(){
                await asnsPage.OpenTopologyMenu();
                await asnsPage.OpenASNsPage();
            })

            using(asnsData.Add, function (add) {
                it(add.description, async function () {
                    expect(await asnsPage.CreateASNs(add)).toBeTruthy();
                    await asnsPage.OpenASNsPage();
                })
            })
            using(asnsData.Update, function (update) {
                it(update.description, async function () {
                    await asnsPage.SearchASNs(update.ASNs);
                    expect(await asnsPage.UpdateASNs(update)).toBeTruthy();
                    await asnsPage.OpenASNsPage();
                })
            })
            using(asnsData.Remove, function (remove) {
                it(remove.description, async function () {
                    await asnsPage.SearchASNs(remove.ASNs);
                    expect(await asnsPage.DeleteASNs(remove)).toBeTruthy();
                    await asnsPage.OpenASNsPage();
                })
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean Up API for ASNs Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})