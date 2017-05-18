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

var WidgetRoutingController = function($scope, cdnService) {

	var getRoutingMethods = function() {
		cdnService.getRoutingMethods()
			.then(function(response) {
				$scope.native = response.cz;
				$scope.thirdParty = response.geo;
				$scope.federated = response.fed;
				$scope.miss = response.miss;
				$scope.static = response.staticRoute;
				$scope.dsr = response.dsr;
				$scope.error = response.err;
				$scope.regionalAlternate = response.regionalAlternate;
				$scope.regionalDenied = response.regionalDenied;
			});
	};

	var init = function() {
		getRoutingMethods();
	};
	init();

};

WidgetRoutingController.$inject = ['$scope', 'cdnService'];
module.exports = WidgetRoutingController;
