/*


 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

 */

var CustomController = function($scope, $sce, $timeout, $location) {

	var pinIframe = function() {
		var headerHeight = $('#header').css("height");
		$('#customFrameWrapper').css("top", headerHeight);
		$('#customFrameWrapper').css("bottom", 0);
	};

	$scope.url = decodeURIComponent($location.search().url);

	$scope.trustSrc = function(src) {
		return $sce.trustAsResourceUrl(src);
	};

	var init = function () {
		$timeout(function () {
			pinIframe();
		}, 200);
	};
	init();

};

CustomController.$inject = ['$scope', '$sce', '$timeout', '$location'];
module.exports = CustomController;