var TableCacheGroupsController = function(cacheGroups, $scope, locationUtils) {

    $scope.cacheGroups = cacheGroups;

    $scope.editCacheGroup = function(id) {
        locationUtils.navigateToPath('/configure/cache-groups/' + id + '/edit');
    };

    $scope.createCacheGroup = function() {
        locationUtils.navigateToPath('/configure/cache-groups/new');
    };

    angular.element(document).ready(function () {
        $('#cacheGroupsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableCacheGroupsController.$inject = ['cacheGroups', '$scope', 'locationUtils'];
module.exports = TableCacheGroupsController;