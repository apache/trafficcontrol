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

var ToolsService = function($http, messageModel, ENV) {

	this.getOSVersions = function() {
		return $http.get(ENV.api.unstable + "osversions").then(
				function(result) {
					return result.data.response;
				},
				function(err) {
					throw err;
				}
			);
	};

	this.generateISO = function(iso) {
		respType = 'arraybuffer';

		return $http.post(ENV.api.unstable + "isos", iso, { responseType:respType }).then(
			function(result) {
				const isoName = iso.hostName + "." + iso.domainName + "-" + iso.osversionDir + ".iso";
				download(result.data, isoName);
				return result.data.response;
			},
			function(err) {
				// apparently there are no alerts sent from this endpoint
				messageModel.setMessages([ { level: 'error', text: err.status.toString() + ': ' + err.statusText } ], false);
				throw err;
			}
		);
	};

};

ToolsService.$inject = ['$http', 'messageModel', 'ENV'];
module.exports = ToolsService;
