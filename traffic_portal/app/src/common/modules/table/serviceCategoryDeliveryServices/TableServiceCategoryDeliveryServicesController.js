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

const TableServiceCategoryDeliveryServicesController = function(serviceCategory, deliveryServices, steeringTargets, filter, $controller, $scope) {

	// extends the TableDeliveryServicesController to inherit common methods
	angular.extend(this, $controller('TableDeliveryServicesController', { tableName: 'scDS', deliveryServices: deliveryServices, steeringTargets: steeringTargets, filter: filter, $scope: $scope }));

	$scope.serviceCategory = serviceCategory;
	$scope.breadCrumbs = [
		{
			text: "Service Categories",
			href: "#!/service-categories"
		},
		{
			getText: () => serviceCategory.name,
			getHref: () => `#!/service-categories/edit?name=${serviceCategory.name}`
		},
		{
			text: "Delivery Services"
		}
	]
};

TableServiceCategoryDeliveryServicesController.$inject = ['serviceCategory', 'deliveryServices', 'steeringTargets', 'filter', '$controller', '$scope'];
module.exports = TableServiceCategoryDeliveryServicesController;
