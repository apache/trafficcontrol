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

var CacheStatsService = function($http, $q, httpService, ENV, messageModel) {

	this.getBandwidth = function(cdnName, start, end) {
		var request = $q.defer();

		var url = ENV.api['root'] + "cache_stats",
			params = { cdnName: cdnName, metricType: 'bandwidth', startDate: start.seconds(00).format(), endDate: end.seconds(00).format()};

		$http.get(url, { params: params })
			.then(
				function(result) {
					request.resolve(result.data.response);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					request.reject();
				}
			);

		return request.promise;
	};

	this.getConnections = function(cdnName, start, end) {
		var request = $q.defer();

		var url = ENV.api['root'] + "cache_stats",
			params = { cdnName: cdnName, metricType: 'connections', startDate: start.seconds(00).format(), endDate: end.seconds(00).format()};

		$http.get(url, { params: params })
			.then(
				function(result) {
					request.resolve(result.data.response);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					request.reject();
				}
			);

		return request.promise;
	};

};

CacheStatsService.$inject = ['$http', '$q', 'httpService', 'ENV', 'messageModel'];
module.exports = CacheStatsService;
