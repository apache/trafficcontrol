var TableCacheGroupsController = function(cacheGroups, $scope, $location) {

    $scope.cacheGroups = cacheGroups;

    $scope.editCacheGroup = function(id) {
        $location.path($location.path() + id);
    };

    angular.element(document).ready(function () {
        $('#cacheGroupsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableCacheGroupsController.$inject = ['cacheGroups', '$scope', '$location'];
module.exports = TableCacheGroupsController;