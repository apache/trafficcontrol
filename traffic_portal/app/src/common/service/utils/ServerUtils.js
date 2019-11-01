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

var ServerUtils = function($window, propertiesModel, userModel) {

	this.isCache = function(server) {
		return server.type && (server.type.indexOf('EDGE') == 0 || server.type.indexOf('MID') == 0);
	};

	this.isEdge = function(server) {
		return server.type && (server.type.indexOf('EDGE') == 0);
	};

	this.isOffline = function(status) {
		return (status == 'OFFLINE' || status == 'ADMIN_DOWN');
	};

	this.offlineReason = function(server) {
		return (server.offlineReason) ? server.offlineReason : 'None';
	};

	this.ssh = function(ip, $event) {
		if (ip && ip.length > 0) {
			$window.location.href = 'ssh://' + userModel.user.username + '@' + ip;
		}
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
	};

	this.openCharts = function(server, $event) {
		if ($event) {
			$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		}
		$window.open(
			propertiesModel.properties.servers.charts.baseUrl + server.hostName,
			'_blank'
		);
	};

};

ServerUtils.$inject = ['$window', 'propertiesModel', 'userModel'];
module.exports = ServerUtils;
