import { browser } from 'protractor';
import { API } from '../CommonUtils/API';
import { LoginPage } from '../PageObjects/LoginPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { DeliveryServicesPage } from '../PageObjects/DeliveryServicesPage.po';

let fs = require('fs')
let using = require('jasmine-data-provider');

let filename = 'Data/TenancyAccess/TestCases.json';
let setupFile = 'Data/TenancyAccess/Setup.json';
let cleanupFile = 'Data/TenancyAccess/Cleanup.json';

let testData = JSON.parse(fs.readFileSync(filename));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let deliveryServicesPage = new DeliveryServicesPage();

using(testData.TenancyAccess, async function(data) {
    using(data.Login, function(login) {
        describe('Traffic Portal - Tenancy Access - ' + login.description,  function(){
            it('Setup', async function() {
                let setupData = JSON.parse(fs.readFileSync(setupFile));
                let output = await api.UseAPI(setupData);
                expect(output).toBeNull();
            })
            it(login.description, async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            using(data.DeliveryService, function(test){
                it('can open delivery service page', async function(){
                    await deliveryServicesPage.OpenServicesMenu();
                    await deliveryServicesPage.OpenDeliveryServicePage();
                })
                it(test.description, async function(){
                    expect(await deliveryServicesPage.SearchDeliveryService(test.DeliveryService)).toBeUndefined();
                })
            })
            it('can logout', async function(){
                expect(await topNavigation.Logout()).toBeTruthy();
            })
            it('Cleanup', async function() {
                let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
                let output = await api.UseAPI(cleanupData);
                expect(output).toBeNull();
            })
        })
    })
})