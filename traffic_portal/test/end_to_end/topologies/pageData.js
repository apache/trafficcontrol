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
	this.topLink=element(by.name('topLink'));
	this.searchFilter=element(by.id('topologiesTable_filter')).element(by.css('label input'));
	this.name=element(by.name('name'));
	this.description=element(by.id('description'));
	this.addChildCacheGroupBtn=element(by.css('.add-child-cg-btn'));
	this.selectFormSubmitButton=element(by.buttonText('Submit'));
	this.selectAllCB=element(by.id('selectAllCB'));
	this.updateButton=element(by.buttonText('Update'));
	this.createButton=element(by.buttonText('Create'));
	this.viewDeliveryServicesMenuItem=element(by.name('viewDeliveryServicesMenuItem'));
	this.viewCacheGroupsMenuItem=element(by.name('viewCacheGroupsMenuItem'));
	this.viewServersMenuItem=element(by.name('viewServersMenuItem'));
};
