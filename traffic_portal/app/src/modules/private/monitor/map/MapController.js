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

var MapController = function(cacheGroups, cacheGroupHealth, $scope, locationUtils, NgMap) {

	$scope.map = NgMap.getMap('cgMap');

	$scope.cacheGroups = [];

	$scope.cacheGroupTypes = [];

	$scope.cgTitle = function(cg) {
		return cg.name + ' (' + cg.type + ')';
	};

	$scope.parentCg = function(cg) {
		return cg.parent ? cg.parent : 'None'
	};

	$scope.secondaryParentCg = function(cg) {
		return cg.secondaryParent ? cg.secondaryParent : 'None'
	};

	$scope.icon = function(cg) {
		var properties = {
			path: 'M8 2.1c1.1 0 2.2 0.5 3 1.3 0.8 0.9 1.3 1.9 1.3 3.1s-0.5 2.5-1.3 3.3l-3 3.1-3-3.1c-0.8-0.8-1.3-2-1.3-3.3 0-1.2 0.4-2.2 1.3-3.1 0.8-0.8 1.9-1.3 3-1.3z',
			fillOpacity: 0.8,
			scale: 3,
			strokeColor: 'white',
			strokeWeight: 2
		}
		// color map markers by type UNLESS there are offline caches, then make red and bigger
		if (parseInt(cg.offline) > 0) {
			properties['fillColor'] = 'red';
			properties['scale'] = 5;
		} else {
			properties['fillColor'] = colors[_.indexOf($scope.cacheGroupTypes, cg.type)];
		}
		return properties;
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	var massageCacheGroups = function() {
		var cgHealthCacheGroups = cacheGroupHealth.cachegroups,
			cgHealth;
		var cgTypes = [];
		_.each(cacheGroups, function(cg) {
			cgTypes.push(cg.typeName);
			cgHealth = _.find(cgHealthCacheGroups, function(cghcg){ return cghcg.name == cg.name });
			$scope.cacheGroups.push(
				{
					id: cg.id,
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
		$scope.cacheGroupTypes = _.uniq(cgTypes);
	};

	var colors = [
		'#3F51B5', // blue
		'#00AAA0', // turquoise
		'#FF7A5A', // orangish
		'#FFB85F', // yellowish
		'#462066', // purple
		'#FCF4D9' // whitish
	];

	var init = function() {
		massageCacheGroups();
	};
	init();


};

MapController.$inject = ['cacheGroups', 'cacheGroupHealth', '$scope', 'locationUtils', 'NgMap'];
module.exports = MapController;
