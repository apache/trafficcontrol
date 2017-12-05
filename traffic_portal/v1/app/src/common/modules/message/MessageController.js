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

var MessageController = function($rootScope, $scope, messageModel, chartModel) {

    $scope.messageData = messageModel;

    $scope.chartData = chartModel.chart;

    $scope.dismissMessage = function(message) {
        messageModel.removeMessage(message);
    };

    $scope.showConnectionLostMsg = function() {
        $('#lostConnectionMsg').show();
    };

    $scope.hideConnectionLostMsg = function() {
        $('#lostConnectionMsg').hide();
    };

    $rootScope.$watch('online', function(newStatus) {
        if (newStatus === false && $scope.chartData.autoRefresh === true) {
            $scope.showConnectionLostMsg();
        }
    });

};

MessageController.$inject = ['$rootScope', '$scope', 'messageModel', 'chartModel'];
module.exports = MessageController;
