import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po'
import { CoordinatesPage } from '../PageObjects/CoordinatesPage.po';
import { API } from '../CommonUtils/API';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';

let fs = require('fs')
let using = require('jasmine-data-provider');

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let coordinatesPage = new CoordinatesPage();


let setupFile = 'Data/Coordinates/Setup.json';
let cleanupFile = 'Data/Coordinates/Cleanup.json';
let filename = 'Data/Coordinates/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

describe('Setup API for coordinates test', function () {
    it('Setup', async function () {
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.Coordinates, async function(coordinatesData){
    using(coordinatesData.Login, function(login){
        describe('Traffic Portal - Coordinates - ' + login.description, function(){

            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open coordinates page', async function(){
                await coordinatesPage.OpenTopologyMenu();
                await coordinatesPage.OpenCoordinatesPage();
            })
            using(coordinatesData.Add, function (add) {
                it(add.description, async function () {
                    expect(await coordinatesPage.CreateCoordinates(add)).toBeTruthy();
                    await coordinatesPage.OpenCoordinatesPage();
                })
            })
            using(coordinatesData.Update, function (update) {
                it(update.description, async function () {
                    await coordinatesPage.SearchCoordinates(update.Name);
                    expect(await coordinatesPage.UpdateCoordinates(update)).toBeTruthy();
                    await coordinatesPage.OpenCoordinatesPage();
                })
            })
          
            using(coordinatesData.Remove, function (remove) {
                it(remove.description, async function () {
                    await coordinatesPage.SearchCoordinates(remove.Name);
                    expect(await coordinatesPage.DeleteCoordinates(remove)).toBeTruthy();
                    await coordinatesPage.OpenCoordinatesPage();
                })
            })

            it('can logout', async function(){
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})


describe('Clean up API for coordinates test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})