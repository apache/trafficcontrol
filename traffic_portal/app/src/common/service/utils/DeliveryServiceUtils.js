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

/**
 * DeliveryServiceUtils provides utilities for dealing Delivery Services.
 */
class DeliveryServiceUtils {
	protocols = {
		0: "HTTP",
		1: "HTTPS",
		2: "HTTP AND HTTPS",
		3: "HTTP TO HTTPS"
	};

	qstrings = {
		0: "USE",
		1: "IGNORE",
		2: "DROP"
	};

	geoProviders = {
		0: "MaxMind",
		1: "Neustar"
	};

	geoLimits = {
		0: "None",
		1: "CZF Only",
		2: "CZF + Country Code(s)",
	};

	rrhs = {
		0: "no cache",
		1: "background_fetch",
		2: "cache_range_requests",
		3: "slice"
	};


	/**
	 * @param {import("angular").IWindowService} $window
	 * @param {import("../../models/PropertiesModel")} propertiesModel
	 */
	constructor($window, propertiesModel) {
		this.$window = $window;
		this.propertiesModel = propertiesModel;
	}

	/**
	 * Opens a new browsing context for the URL for the "charts" for a given
	 * Delivery Service.
	 *
	 * Note that this will break in half if the custom charts base URL is not
	 * configured, and callers accept that responsibility.
	 *
	 * @deprecated Links should be links; use ng-href with the appropriate link
	 * on an Anchor element instead of utilizing this method in a click handler.
	 *
	 * @param {{xmlId: string}} ds
	 */
	openCharts(ds) {
		this.$window.open(
			this.propertiesModel.properties.deliveryServices.charts.customLink.baseUrl + ds.xmlId,
			"_blank"
		);
	}

	/**
	 * Maps the targets of Steering Delivery Services to those Delivery
	 * Services.
	 *
	 * @param {string[]} xmlIds
	 * @param {{deliveryService: string; targets: {deliveryService: string}[]}[]} steeringConfigs
	 * @returns {Record<string, Set<string>>}
	 */
	getSteeringTargetsForDS(xmlIds, steeringConfigs) {
		const dsTargets = Object.fromEntries(xmlIds.map(x => [x, new Set()]));
		for (const config of steeringConfigs) {
			for (const target of config.targets) {
				for (const xmlID of xmlIds) {
					if (target.deliveryService === xmlID) {
						dsTargets[xmlID].add(config.deliveryService);
					}
				}
			}
		}
		return dsTargets;
	}
}

DeliveryServiceUtils.$inject = ["$window", "propertiesModel"];
module.exports = DeliveryServiceUtils;
