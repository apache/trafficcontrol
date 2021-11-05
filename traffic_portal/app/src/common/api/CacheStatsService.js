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

var CacheStatsService = function($http, ENV, messageModel) {

	this.getBandwidth = function(cdnName, start, end) {
		const url = ENV.api.unstable + "cache_stats";
		const params = {
			cdnName: cdnName,
			metricType: 'bandwidth',
			startDate: start.seconds(0).format(),
			endDate: end.seconds(0).format()
		};

		return $http.get(url, { params: params }).then(
			function(result) {
				return result.data.response;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.getConnections = function(cdnName, start, end) {
		const url = ENV.api.unstable + "cache_stats";
		const params = {
			cdnName: cdnName,
			metricType: 'connections',
			startDate: start.seconds(0).format(),
			endDate: end.seconds(0).format()
		};

		return $http.get(url, { params: params }).then(
			function(result) {
				return result.data.response;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

};

CacheStatsService.$inject = ['$http', 'ENV', 'messageModel'];
module.exports = CacheStatsService;
