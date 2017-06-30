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

var DatesController = function(customLabel, dateRange, $rootScope, $scope) {

    $scope.dateRange = dateRange;

    $scope.customLabel = customLabel;

    $scope.openStart = function($event) {
        $event.preventDefault();
        $event.stopPropagation();

        $scope.startOpened = true;
    };

    $scope.openEnd = function($event) {
        $event.preventDefault();
        $event.stopPropagation();

        $scope.endOpened = true;
    };

    $scope.dateOptions = {
        formatYear: 'yy',
        startingDay: 0,
        showWeeks: false
    };

    $scope.changeDates = function(start, end) {
        $rootScope.$broadcast('datesController::dateChange', { start: start, end: end });
    };

    angular.element(document).ready(function () {
        $scope.changeDates(dateRange.start, dateRange.end);
    });

};

DatesController.$inject = ['customLabel', 'dateRange', '$rootScope', '$scope'];
module.exports = DatesController;
