import { browser } from 'protractor';
import { LoginPage } from '../PageObjects/LoginPage.po';
import { TopNavigationPage } from '../PageObjects/TopNavigationPage.po';
import { API } from '../CommonUtils/API';
import { TopologiesPage } from '../PageObjects/Topologies.po';
import { CacheGroupPage } from '../PageObjects/CacheGroup.po';

let fs = require('fs')
let using = require('jasmine-data-provider');

let setupFile = 'Data/Topology/Setup.json';
let cleanupFile = 'Data/Topology/Cleanup.json';
let filename = 'Data/Topology/TestCases.json';
let testData = JSON.parse(fs.readFileSync(filename));

let api = new API();
let loginPage = new LoginPage();
let topNavigation = new TopNavigationPage();
let topologiesPage = new TopologiesPage();
let cacheGroupPage = new CacheGroupPage();

describe('Setup CacheGroup for Topology Test', function(){
    it('Setup', async function(){
        let setupData = JSON.parse(fs.readFileSync(setupFile));
        let output = await api.UseAPI(setupData);
        expect(output).toBeNull();
    })
})

using(testData.Topologies, async function(topologiesData){
    using(topologiesData.Login, async function(login){
        describe('Traffic Portal -  Topologies - ' + login.description, function(){
            it('can login', async function(){
                browser.get(browser.params.baseUrl);
                await loginPage.Login(login.username, login.password);
                expect(await loginPage.CheckUserName(login.username)).toBeTruthy();
            })
            it('can open topologies page', async function(){
                await topologiesPage.OpenTopologyMenu();
                await topologiesPage.OpenTopologiesPage();
            })
            using(topologiesData.Add, function (add) {
                if(add.description.includes("cannot")){
                    it(add.description, async function () {
                        expect(await topologiesPage.CreateTopologies(add)).toBeUndefined();
                        await topologiesPage.OpenTopologiesPage();
                    })
                }else{
                    it(add.description, async function () {
                        expect(await topologiesPage.CreateTopologies(add)).toBeTruthy();
                        await topologiesPage.OpenTopologiesPage();
                    })
                }
            })
            using(topologiesData.Update, function(update){
                if(update.description.includes("cannot")){
                    it(update.description, async function () {
                        await topologiesPage.SearchTopologies(update.Name);
                        expect(await topologiesPage.UpdateTopologies(update)).toBeUndefined();
                        await topologiesPage.OpenTopologiesPage();
                    })
                }else{
                    it(update.description, async function () {
                        await topologiesPage.SearchTopologies(update.Name);
                        expect(await topologiesPage.UpdateTopologies(update)).toBeTruthy();
                        await topologiesPage.OpenTopologiesPage();
                    })
                }
            })
            using(topologiesData.Remove, function (remove) {
                it(remove.description, async function () {
                    await topologiesPage.SearchTopologies(remove.Name);
                    expect(await topologiesPage.DeleteTopologies(remove)).toBeTruthy();
                    await topologiesPage.OpenTopologiesPage();
                })
            })
            it('can open cachegroup page', async function(){
                await cacheGroupPage.OpenCacheGroupsPage();
            })
            using(topologiesData.RemoveCG, function (removeCacheGroup) {
                it(removeCacheGroup.description, async function () {
                    await cacheGroupPage.SearchCacheGroups(removeCacheGroup.Name)
                    expect(await cacheGroupPage.DeleteCacheGroups(removeCacheGroup.Name, removeCacheGroup.validationMessage)).toBeTruthy();
                }) 
            })
            it('can logout', async function () {
                expect(await topNavigation.Logout()).toBeTruthy();
            })
        })
    })
})

describe('Clean Up CacheGroup for Topology Test', function () {
    it('Cleanup', async function () {
        let cleanupData = JSON.parse(fs.readFileSync(cleanupFile));
        let output = await api.UseAPI(cleanupData);
        expect(output).toBeNull();
    })
})