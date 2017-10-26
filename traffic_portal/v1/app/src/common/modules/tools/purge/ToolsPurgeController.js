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

var ToolsPurgeController = function($scope, $uibModal, $stateParams, messageModel, deliveryServicesModel, dateUtils, formUtils, deliveryServiceService, userService) {

    var getPurgeJobs = function() {
        $scope.loadingPurgeJobs = true;
        deliveryServiceService.getPurgeJobs($scope.deliveryService.id, false)
            .then(function(response) {
                $scope.loadingPurgeJobs = false;
                $scope.purgeJobs = response;
            });
    };

    var createPurgeJob = function(jobParams) {
        userService.createUserPurgeJob(jobParams)
            .then(function(response) {
                getPurgeJobs();
            });
    };

    $scope.deliveryService = deliveryServicesModel.getDeliveryService($stateParams.deliveryServiceId);

    $scope.dateFormat = dateUtils.dateFormat;

    $scope.resetPurgeJobData = function() {
        $scope.newPurgeJobData = {
            dsId: $scope.deliveryService.id,
            dsXmlId: $scope.deliveryService.xmlId,
            regex: '',
            ttl: 672,
            startTime: moment().format('YYYY-MM-DD HH:mm:ss')
        };
    };
    $scope.resetPurgeJobData();

    $scope.confirmPurgeJobCreate = function(newPurgeJob) {
        var params = {
            title: 'Confirmation required',
            message: 'Are you sure you want to invalidate content for ' + $scope.deliveryService.orgServerFqdn + newPurgeJob.regex
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
            controller: 'DialogConfirmController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result
            .then(
            function() {
                createPurgeJob($scope.newPurgeJobData);
                $scope.resetPurgeJobData();
                $scope.purgeForm.$setPristine();
            },
            function () {
            });
    };

    $scope.toDate = function(date) {
        return moment(date).toDate(); // hack for https://issues.apache.org/jira/browse/TC-78
    };

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    angular.element(document).ready(function () {
        getPurgeJobs();
    });

};

ToolsPurgeController.$inject = ['$scope', '$uibModal', '$stateParams', 'messageModel', 'deliveryServicesModel', 'dateUtils', 'formUtils', 'deliveryServiceService', 'userService'];
module.exports = ToolsPurgeController;
