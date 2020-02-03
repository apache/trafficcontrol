/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

module.exports = function(){
	this.moreBtn=element(by.name('moreBtn'));
	this.compareProfilesMenuItem=element(by.css('a[ng-click*=compareProfiles]'));
	this.compareDropdown1=element(by.name('compareDropdown1'));
	this.compareDropdown2=element(by.name('compareDropdown2'));
	this.compareSubmit=element(by.name('compareSubmit'));
	this.description=element(by.id('description'));
	this.createProfileButton=element(by.name('createProfileButton'));
	this.name=element(by.name('name'));
	this.cdn=element(by.name('cdn'));
	this.type=element(by.name('type'));
	this.routingDisabled=element(by.name('routingDisabled'));
	this.description=element(by.id('description'));
	this.createButton=element(by.buttonText('Create'));
	this.deleteButton=element(by.buttonText('Delete'));
	this.updateButton=element(by.buttonText('Update'));
	this.searchFilter=element(by.id('profilesTable_filter')).element(by.css('label input'));
	this.confirmWithNameInput=element(by.name('confirmWithNameInput'));
	this.deletePermanentlyButton=element(by.buttonText('Delete Permanently'));
};
