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

var DeliveryServiceUtils = function ($window, propertiesModel) {

	this.protocols = {
		0: "HTTP",
		1: "HTTPS",
		2: "HTTP AND HTTPS",
		3: "HTTP TO HTTPS"
	};

	this.qstrings = {
		0: "USE",
		1: "IGNORE",
		2: "DROP"
	};

	this.geoProviders = {
		0: "MaxMind",
		1: "Neustar"
	};

	this.geoLimits = {
		0: "None",
		1: "CZF Only",
		2: "CZF + Country Code(s)",
	};

	this.rrhs = {
		0: "no cache",
		1: "background_fetch",
		2: "cache_range_requests",
		3: "slice"
	};

	this.openCharts = function (ds) {
		$window.open(
			propertiesModel.properties.deliveryServices.charts.customLink.baseUrl + ds.xmlId,
			'_blank'
		);
	};

	this.getSteeringTargetsForDS = function (xmlIds, steeringConfigs) {
		const targetsFor = new Set();
		const dsTargets = {};
		xmlIds.forEach(xmlId => dsTargets[xmlId] = new Set());
		steeringConfigs.forEach(config => {
			config.targets.forEach(target => {
				xmlIds.forEach(xmlId => {
					if (target.deliveryService === xmlId) {
						dsTargets[xmlId].add(config.deliveryService);
					}
				})
			})
		})
		return dsTargets;
	}

};

DeliveryServiceUtils.$inject = ['$window', 'propertiesModel'];
module.exports = DeliveryServiceUtils;
