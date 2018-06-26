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

var FormISOController = function(servers, osversions, $scope, $anchorScroll, formUtils, toolsService, messageModel, FileSaver, Blob) {

	$scope.servers = servers;

	$scope.osversions = osversions;

	$scope.selectedServer = {};

	$scope.falseTrue = [
		{ value: 'yes', label: 'yes' },
		{ value: 'no', label: 'no' }
	];

	$scope.iso = {
		dhcp: 'no',
		stream: 'no'
	};

	$scope.isDHCP = function() {
		return $scope.iso.dhcp == 'yes';
	};

	$scope.fqdn = function(server) {
		return server.hostName + '.' + server.domainName;
	};

	$scope.copyServerAttributes = function() {
		$scope.iso = angular.extend($scope.iso, $scope.selectedServer);
	};

	$scope.generate = function(iso) {
		toolsService.generateISO(iso)
			.then(function(result) {
				$anchorScroll(); // scrolls window to top
				if (iso.stream != 'yes') {
                    messageModel.setMessages([{level: 'success', text: 'ISO created at ' + result.isoURL}], false);
                }
                else {
					//var isoStr = result.iso.replace(/\n/g, "");
					//var decodedIso = Base64.atob(result.iso);
					//isoStr += '=';
					//alert(isoStr.length)
                    //var decodedIso = $base64.atob(result.iso);
					var decodedIso = atob(result.iso);
					var newData = new Blob([decodedIso], { type: 'application/x-iso9660-image' } );
					//var encodedIso = new Blob([result.iso], { type: 'application/x-iso9660-image' } );
					//var file = new File([decodedIso], result.name, { type: 'application/x-iso9660-image' } );
					//alert(newData.size + " : " + encodedIso.size);

					FileSaver.saveAs(newData, result.name);
				}
			});
	};

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

    // function b64DecodeUnicode(str) {
    //     return decodeURIComponent(Array.prototype.map.call(atob(str), function(c) {
    //         return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
    //     }).join(''))
    // }

};

FormISOController.$inject = ['servers', 'osversions', '$scope', '$anchorScroll', 'formUtils', 'toolsService', 'messageModel', 'FileSaver', 'Blob'];
module.exports = FormISOController;
