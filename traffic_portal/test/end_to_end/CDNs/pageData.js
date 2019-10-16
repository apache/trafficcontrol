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
	this.name=element(by.name('name'));
	this.domainName=element(by.name('domainName'));
	this.dnssecEnabled=element(by.name('dnssecEnabled'));
	this.createButton=element(by.buttonText('Create'));
	this.deleteButton=element(by.buttonText('Delete'));
	this.updateButton=element(by.buttonText('Update'));
	this.searchFilter=element(by.id('cdnsTable_filter')).element(by.css('label input'));
	this.deletePermanentlyButton=element(by.buttonText('Delete Permanently'));
	this.moreButton=element(by.buttonText('More'));
	this.manageDnssecKeysButton=element(by.linkText('Manage DNSSEC Keys'));
	this.generateDnssecKeysButton=element(by.buttonText('Generate DNSSEC Keys'));
	this.kskExpirationDays=element(by.name('kskExpirationDays'));
	this.regenerateButton=element(by.buttonText('Regenerate'));
	this.confirmInput=element(by.name('confirmEnterForm')).element(by.tagName('input'));
	this.confirmButton=element(by.buttonText('Confirm'));
	this.expirationDate=element(by.name('expirationDate'));
	this.regenerateDnssecKeysButton=element(by.buttonText('Regenerate DNSSEC Keys'));
	this.regenerateKskButton=element(by.buttonText('Regenerate KSK'));
	this.generateButton=element(by.buttonText('Generate'));
};