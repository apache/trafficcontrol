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

var FooterController = function($scope, propertiesModel) {

    var chunk = function(array, n) {
        var retval = [];
        for (var i = 0, len = array.length; i < len; i += n) {
            retval.push(array.slice(i, i + n));
        }
        return retval;
    };

    $scope.footerChunks = chunk(propertiesModel.properties.footer.links, 4);

};

FooterController.$inject = ['$scope', 'propertiesModel'];
module.exports = FooterController;
