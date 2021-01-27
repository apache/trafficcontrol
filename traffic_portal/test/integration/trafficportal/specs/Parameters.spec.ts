import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po'
import { ParametersPage } from '../PageObjects/ParametersPage.po';
import { API } from '../CommonUtils/API';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';

let fs = require('fs')
let using = require('jasmine-data-provider');

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let parametersPage = new ParametersPage();


let setupFile = 'Data/Parameters/Setup.json';
let cleanupFile = 'Data/Parameters/Cleanup.json';
let filename = 'Data/Parameters/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

describe('Setup API for parameter test', function () {
    it('Setup', async function () {
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.Parameters, async function(parametersData){
    using(parametersData.Login, function(login){
        describe('Traffic Portal - Parameters - ' + login.description, function(){

            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open parameters page', async function(){
                await parametersPage.OpenConfigureMenu();
                await parametersPage.OpenParametersPage();
            })
            using(parametersData.Add, function (add) {
                it(add.description, async function () {
                    expect(await parametersPage.CreateParameter(add)).toBeTruthy();
                    await parametersPage.OpenParametersPage();
                })
            })
            using(parametersData.Update, function (update) {
                it(update.description, async function () {
                    await parametersPage.SearchParameter(update.Name);
                    expect(await parametersPage.UpdateParameter(update)).toBeTruthy();
                    await parametersPage.OpenParametersPage();
                })
            })
          
            using(parametersData.Remove, function (remove) {
                it(remove.description, async function () {
                    await parametersPage.SearchParameter(remove.Name);
                    expect(await parametersPage.DeleteParameter(remove)).toBeTruthy();
                    await parametersPage.OpenParametersPage();
                })
            })

            it('can logout', async function(){
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean up API for parameter test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})