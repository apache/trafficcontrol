var TableCacheGroupsController = function(cacheGroups, $scope) {

    $scope.cacheGroups = cacheGroups.response;

    angular.element(document).ready(function () {
        $('#cacheGroupsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableCacheGroupsController.$inject = ['cacheGroups', '$scope'];
module.exports = TableCacheGroupsController;