var TableStatusesController = function(statuses, $scope, locationUtils) {

    $scope.statuses = statuses;

    $scope.editStatus = function(id) {
        locationUtils.navigateToPath('/admin/statuses/' + id + '/edit');
    };

    $scope.createStatus = function() {
        locationUtils.navigateToPath('/admin/statuses/new');
    };

    angular.element(document).ready(function () {
        $('#statusesTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableStatusesController.$inject = ['statuses', '$scope', 'locationUtils'];
module.exports = TableStatusesController;