import { browser } from 'protractor'
import { LoginPage } from '../PageObjects/LoginPage.po'
import { ProfilesPage } from '../PageObjects/ProfilesPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';

let fs = require('fs')
let using = require('jasmine-data-provider');

let setupFile = 'Data/Profiles/Setup.json';
let cleanupFile = 'Data/Profiles/Cleanup.json';
let filename = 'Data/Profiles/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));


let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let profilesPage = new ProfilesPage();

describe('Setup API for Profiles', function () {
    it('Setup', async function () {
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})
using(testData.Profiles, async function(profilesData){
    using(profilesData.Login, function(login){
        describe('Traffic Portal - Profiles - ' + login.description, function(){
            it('can login', async function () {
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open profiles page', async function () {
                await profilesPage.OpenConfigureMenu();
                await profilesPage.OpenProfilesPage();
            })
            using(profilesData.Add, function (add) {
                it(add.description, async function () {
                    expect(await profilesPage.CreateProfile(add)).toBeTruthy();
                    await profilesPage.OpenProfilesPage();
                })
            })
            using(profilesData.Update, function (update) {
                it(update.description, async function () {
                    await profilesPage.SearchProfile(update.Name);
                    expect(await profilesPage.UpdateProfile(update)).toBeTruthy();
                    await profilesPage.OpenProfilesPage();
                })
            })
            using(profilesData.Remove, function (remove) {
                it(remove.description, async function () {
                    await profilesPage.SearchProfile(remove.Name);
                    expect(await profilesPage.DeleteProfile(remove)).toBeTruthy();
                    await profilesPage.OpenProfilesPage();
                })
            })

        })
    })
})
describe('Clean up API for Profiles', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})