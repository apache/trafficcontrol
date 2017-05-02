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

var DashboardsMapController = function(cacheGroups, cacheGroupHealth, $scope, NgMap) {

	$scope.map = NgMap.getMap('cgMap');

	$scope.cacheGroups = [];

	$scope.parentCg = function(cg) {
		return cg.parent ? cg.parent : 'None'
	};

	$scope.secondaryParentCg = function(cg) {
		return cg.secondaryParent ? cg.secondaryParent : 'None'
	};

	var massageCacheGroups = function() {
		var cgHealthCacheGroups = cacheGroupHealth.cachegroups,
			cgHealth;
		_.each(cacheGroups, function(cg) {
			cgHealth = _.find(cgHealthCacheGroups, function(cghcg){ return cghcg.name == cg.name });
			$scope.cacheGroups.push(
				{
					name: cg.name,
					parent: cg.parentCachegroupName,
					secondaryParent: cg.secondaryParentCachegroupName,
					pos: [ cg.latitude, cg.longitude ],
					type: cg.typeName,
					offline: cgHealth ? cgHealth.offline : '-',
					online: cgHealth ? cgHealth.online : '-'
				}
			);
		});
	};

	var init = function() {
		massageCacheGroups();
	};
	init();


};

DashboardsMapController.$inject = ['cacheGroups', 'cacheGroupHealth', '$scope', 'NgMap'];
module.exports = DashboardsMapController;
