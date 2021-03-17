import { browser, by, element } from 'protractor'
import { BasePage } from './BasePage.po';
import { SideNavigationPage } from './SideNavigationPage.po';

interface TopologiesData {
    description: string;
    Name:string;
    DescriptionData: string;
    Type: string;
    CacheGroup: string;
    validationMessage: string;
}
export class TopologiesPage extends BasePage {
    private btnCreateNewTopologies = element(by.xpath("//button[@title='Create Topology']"));
    private txtCacheGroupType = element(by.name('selectFormDropdown'));
    private txtSearch = element(by.id('topologiesTable_filter')).element(by.css('label input'));
    private txtName = element(by.name('name'));
    private txtDescription = element(by.id('description'));
    private btnAddCacheGroup = element(by.xpath("//a[@title='Add child cache groups to TOPOLOGY']"));
    private txtSearchCacheGroup = element(by.id('availableCacheGroupsTable_filter')).element(by.css('label input'));
    private btnDelete = element(by.xpath("//button[text()='Delete']"));
    private txtConfirmName = element(by.name('confirmWithNameInput'));
    private config = require('../config');
    private randomize = this.config.randomize;
    async OpenTopologiesPage(){
        let snp = new SideNavigationPage();
        await snp.NavigateToTopologiesPage();
    }
    async OpenTopologyMenu(){
        let snp = new SideNavigationPage();
        await snp.ClickTopologyMenu();
    }
    async CreateTopologies(topologies:TopologiesData){
        let result = false;
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        await snp.NavigateToTopologiesPage();
        //click '+'
        await this.btnCreateNewTopologies.click();
        await this.txtName.sendKeys(topologies.Name + this.randomize)
        await this.txtDescription.sendKeys(topologies.DescriptionData + this.randomize)
        //click add cache group +
        await this.btnAddCacheGroup.click();
        //choose type
        await this.txtCacheGroupType.sendKeys(topologies.Type);
        await basePage.ClickSubmit();
        //choose Cachegroup
        await this.txtSearchCacheGroup.sendKeys(topologies.CacheGroup + this.randomize)
        if(await browser.isElementPresent(by.xpath("//td[@data-search='^" + topologies.CacheGroup + this.randomize + "$']")) === true){
            await element(by.xpath("//td[@data-search='^" + topologies.CacheGroup + this.randomize + "$']")).click();
        }else{
            throw new Error('Test '+topologies.description+' Failed because '+topologies.CacheGroup+' does not display to assign');
        }
        
        if(await basePage.ClickSubmit() == false){
            result = undefined;
            await snp.NavigateToTopologiesPage();
            throw new Error('Test '+topologies.description+' Failed because cannot click on Submit button');
        }
        if(result != undefined){
            if(await basePage.ClickCreate() == false){
                result = undefined;
                await snp.NavigateToTopologiesPage();
                throw new Error('Test '+topologies.description+' Failed because cannot click on Create button');
            }else{
                result = await basePage.GetOutputMessage().then(value => value === topologies.validationMessage);
                if (topologies.description == 'create a Topologies with empty cachegroup (no server)'){
                    result = await basePage.GetOutputMessage().then(value => value.includes(topologies.validationMessage));
                }
            }
            return result;
        }
    }
       
    async SearchTopologies(nameTopologies:string){
        let name = nameTopologies + this.randomize;
        let result = false;
        let snp = new SideNavigationPage();
        await snp.NavigateToTopologiesPage();
        await this.txtSearch.clear();
        await this.txtSearch.sendKeys(name);
        if (await browser.isElementPresent(element(by.xpath("//td[@data-search='^" + name + "$']"))) == true) {
            await element(by.xpath("//td[@data-search='^" + name + "$']")).click();
            result = true;
        } else {
            result = undefined;
        }
        return result;
    }
    async UpdateTopologies(topologies){
        let result = false;
        let basePage = new BasePage();
        let snp = new SideNavigationPage();
        if(topologies.Type != undefined){
             //click add cache group +
            await this.btnAddCacheGroup.click();
            //choose type
            await this.txtCacheGroupType.sendKeys(topologies.Type);
            await basePage.ClickSubmit();
            await this.txtSearchCacheGroup.sendKeys(topologies.CacheGroup + this.randomize)
            await element(by.xpath("//td[@data-search='^" + topologies.CacheGroup + this.randomize + "$']")).click();
            if(await basePage.ClickSubmit() == false){
                result = undefined;
            }
        }
        //if TypeChild is not empty
        if(topologies.TypeChild != undefined){
            //add cachegroup child to cachegroup
            await element(by.xpath("//a[@title='Add child cache groups to " + topologies.CacheGroup + this.randomize + "']//i[1]")).click();
            //choose type
            await this.txtCacheGroupType.sendKeys(topologies.TypeChild);
            //if cannot click submit return undefined
            if(await basePage.ClickSubmit() == false ){
                result = undefined;
            }else{
                //search and send in cachegroup child then click submit
                await this.txtSearchCacheGroup.sendKeys(topologies.CacheGroupChild + this.randomize)
                await element(by.xpath("//td[@data-search='^" + topologies.CacheGroupChild + this.randomize + "$']")).click();
                await basePage.ClickSubmit(); 
                if(topologies.TypeGrandChild != undefined){
                    //add grandchild to cachegroup
                    await element(by.xpath("//a[@title='Add child cache groups to " + topologies.CacheGroupChild + this.randomize + "']//i[1]")).click();
                    //choose type
                    await this.txtCacheGroupType.sendKeys(topologies.TypeGrandChild);
                      //if cannot click submit return undefined
                    if(await basePage.ClickSubmit() == false ){
                        result = undefined;
                    }else{
                        //search and send in cachegroup grandchild then click submit
                        await this.txtSearchCacheGroup.sendKeys(topologies.CacheGroupGrandChild + this.randomize)
                        await element(by.xpath("//td[@data-search='^" + topologies.CacheGroupGrandChild + this.randomize + "$']")).click();
                        await basePage.ClickSubmit(); 
                        await basePage.ClickNo();
                    }
                }
            }
        }
        if(topologies.CacheGroupNeedSP != undefined){
            if(await element(by.xpath("//a[@title='Set Secondary Parent Cache Group for " + topologies.CacheGroupNeedSP + this.randomize + "']")).isDisplayed() == true){
                await element(by.xpath("//a[@title='Set Secondary Parent Cache Group for " + topologies.CacheGroupNeedSP + this.randomize + "']")).click();
                await this.txtCacheGroupType.sendKeys(topologies.SecondaryParent + this.randomize );
                if(await basePage.ClickSubmit() == false){
                    result = undefined;
                }
            }else{
                result = undefined;
            }
        }
         //message check
        //if result in the beginning is not undefined
        if(result != undefined){
            if(await basePage.ClickUpdate() == false){
                result = undefined;
                await snp.NavigateToTopologiesPage();
            }else{
                result = await basePage.GetOutputMessage().then(function (value) {
                    if (topologies.validationMessage == value) {
                        return true;
                    } else {
                        return false;
                    }
                })
                if(topologies.validationWarning == true){
                    let warningMessage = ""+topologies.Type+"-typed cachegroup "+topologies.CacheGroup+this.randomize+" is a parent of "+topologies.CacheGroupChild+this.randomize+", unexpected behavior may result"
                    result = await basePage.GetOutputWarning().then(function (value) {
                        if ( warningMessage == value) {
                            return true;
                        } else {
                            return false;
                        }
                    })

                }
            }
            return result;
        }
        
    }
    async DeleteTopologies(topologies){
        let name = topologies.Name + this.randomize;
        let result = false;
        let basePage = new BasePage();
        await this.btnDelete.click();
        await this.txtConfirmName.sendKeys(name);
        await basePage.ClickDeletePermanently();
        result = await basePage.GetOutputMessage().then(function (value) {
            if (topologies.validationMessage == value) {
                return true;
            } else {
                return false;
            }
        })
        return result;
    }
}