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

exports.config = {
	framework: 'jasmine',
	seleniumAddress: 'http://localhost:4444/wd/hub',
	baseUrl: 'https://localhost:4443',
	getPageTimeout: 30000,

	capabilities: {
		'browserName': 'chrome'
	},
	params: {
		adminUser: 'admin',
		adminPassword: 'twelve'
	},
	jasmineNodeOpts: {defaultTimeoutInterval: 600000},

	suites: {
		loginTests: 'login/login-spec.js',
		allTests: [
			'login/login-spec.js',
			'CDNs/cdns-spec.js',
			'cacheGroups/cache-groups-spec.js',
			'profiles/profiles-spec.js',
			'divisions/divisions-spec.js',
			'regions/regions-spec.js',
			'physLocations/phys-locations-spec.js',
			'serverCapabilities/server-capabilities-spec.js',
			'deliveryServices/delivery-services-spec.js',
			'servers/servers-spec.js'
		]
	}
};
