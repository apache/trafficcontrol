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

var AboutController = function($scope, $sce, $timeout, propertiesModel) {

    var pinIframe = function() {
        var headerHeight = $('#header').css("height"),
            footerHeight = $('#footer').css("height");

        $('#aboutFrameWrapper').css("top", headerHeight);
        $('#aboutFrameWrapper').css("bottom", footerHeight);
    };

    $scope.properties = propertiesModel.properties;

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

AboutController.$inject = ['$scope', '$sce', '$timeout', 'propertiesModel'];
module.exports = AboutController;