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
	this.viewCapabilitiesMenuItem=element(by.css('a[ng-click*=viewCapabilities]'));
	this.addCapabilityBtn=element(by.name('addCapabilityBtn'));
	this.selectFormDropdown=element(by.name('selectFormDropdown'));
	this.status=element(by.name('status'));
	this.hostName=element(by.name('hostName'));
	this.domainName=element(by.name('domainName'));
	this.cdn=element(by.name('cdn'));
	this.cachegroup=element(by.name('cachegroup'));
	this.type=element(by.name('type'));
	this.profile=element(by.name('profile'));
	this.interfaceName=element(by.name('interfaceName'));
	this.ipAddress=element(by.name('ipAddress'));
	this.ipNetmask=element(by.name('ipNetmask'));
	this.ipGateway=element(by.name('ipGateway'));
	this.interfaceMtu=element(by.name('interfaceMtu'));
	this.physLocation=element(by.name('physLocation'));
	this.createButton=element(by.buttonText('Create'));
	this.deleteButton=element(by.buttonText('Delete'));
	this.updateButton=element(by.buttonText('Update'));
	this.submitButton=element(by.buttonText('Submit'));
	this.searchFilter=element(by.id('serversTable_filter')).element(by.css('label input'));
	this.confirmWithNameInput=element(by.name('confirmWithNameInput'));
	this.deletePermanentlyButton=element(by.buttonText('Delete Permanently'));
};