import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po'
import { PhysLocationsPage } from '../PageObjects/PhysLocationsPage.po';
import { API } from '../CommonUtils/API';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';


let fs = require('fs')
let using = require('jasmine-data-provider');


let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let physlocationsPage = new PhysLocationsPage();

let setupFile = 'Data/PhysLocations/Setup.json';
let cleanupFile = 'Data/PhysLocations/Cleanup.json';
let filename = 'Data/PhysLocations/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

describe('Setup API for physlocation test', function () {
    it('Setup', async function () {
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.PhysLocations, async function(physlocationsData){
    using(physlocationsData.Login, function(login){
        describe('Traffic Portal - PhysLocation - ' + login.description, function(){

            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open parameters page', async function(){
                await physlocationsPage.OpenConfigureMenu();
                await physlocationsPage.OpenPhysLocationPage();
            })
            using(physlocationsData.Add, function (add) {
                it(add.description, async function () {
                    expect(await physlocationsPage.CreatePhysLocation(add)).toBeTruthy();
                    await physlocationsPage.OpenPhysLocationPage();
                })
            })
            using(physlocationsData.Update, function (update) {
                it(update.description, async function () {
                    await physlocationsPage.SearchPhysLocation(update.Name);
                    expect(await physlocationsPage.UpdatePhysLocation(update)).toBeTruthy();
                    await physlocationsPage.OpenPhysLocationPage();
                })
            })
          
            using(physlocationsData.Remove, function (remove) {
                it(remove.description, async function () {
                    await physlocationsPage.SearchPhysLocation(remove.Name);
                    expect(await physlocationsPage.DeletePhysLocation(remove)).toBeTruthy();
                    await physlocationsPage.OpenPhysLocationPage();
                })
            })
        })
    })
})

describe('Clean up API for physlocation test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})