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
