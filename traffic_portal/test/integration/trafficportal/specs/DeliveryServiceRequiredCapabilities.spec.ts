import { browser } from 'protractor';
import { API } from '../CommonUtils/API';
import { LoginPage } from '../PageObjects/LoginPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { DeliveryServicesPage } from '../PageObjects/DeliveryServicesPage.po';

let fs = require('fs')
let using = require('jasmine-data-provider');

let testFile = 'Data/DeliveryServiceRequiredCapabilities/TestCases.json';
let setupFile = 'Data/DeliveryServiceRequiredCapabilities/Setup.json';
let cleanupFile = 'Data/DeliveryServiceRequiredCapabilities/Cleanup.json';

let testData = JSON.parse(fs.readFileSync(testFile));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let deliveryServicesPage = new DeliveryServicesPage();

using(testData.DeliveryServiceRequiredCapabilities, async function(data) {
    using(data.Login, function(login) {
        describe('Traffic Portal - Delivery Service Required Capabilities - ' + login.description,  function(){
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
            it('can open delivery service page', async function(){
                await deliveryServicesPage.OpenServicesMenu();
                await deliveryServicesPage.OpenDeliveryServicePage();
            })
            using(data.Add, function(test){
                it(test.description, async function(){
                    await deliveryServicesPage.SearchDeliveryService(test.DeliveryService);
                    expect(await deliveryServicesPage.AddRequiredServerCapabilities(test.ServerCapability, test.validationMessage)).toBeTruthy();
                    await deliveryServicesPage.OpenDeliveryServicePage();
                })
            })
            using(data.Remove, function(test){
                it(test.description, async function(){
                    await deliveryServicesPage.SearchDeliveryService(test.DeliveryService);
                    expect(await deliveryServicesPage.RemoveRequiredServerCapabilities(test.ServerCapability, test.validationMessage)).toBeTruthy();
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