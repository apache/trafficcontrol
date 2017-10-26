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

var DSConfigEditController = function(deliveryService, $scope, $uibModalInstance) {

    $scope.deliveryService = deliveryService;

    $scope.close = function () {
        $uibModalInstance.dismiss('close');
    };

    $scope.isHTTP = function(ds) {
        return ds.type.indexOf('HTTP') !== -1;
    };

    $scope.edgeFQDNs = function(ds) {
        var urlString = '';
        if (_.isArray(ds.exampleURLs) && ds.exampleURLs.length > 0) {
            for (var i = 0; i < ds.exampleURLs.length; i++) {
                urlString += ds.exampleURLs[i] + '\n';
            }
        }
        return urlString;
    };

    $scope.rangeRequestHandling = function(ds) {
        var rrh = '';
        if (ds.rangeRequestHandling == '0') {
            rrh = 'Do not cache range requests';
        } else if (ds.rangeRequestHandling == '1') {
            rrh = 'Background fetch';
        } else if (ds.rangeRequestHandling == '2') {
            rrh = 'Cache range requests';
        }
        return rrh;
    };

};

DSConfigEditController.$inject = ['deliveryService', '$scope', '$uibModalInstance'];
module.exports = DSConfigEditController;
