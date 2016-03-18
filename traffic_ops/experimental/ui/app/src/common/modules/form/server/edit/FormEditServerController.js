var FormEditServerController = function(server, $scope, $controller, $uibModal, $anchorScroll, locationUtils, serverService) {

    // extends the FormServerController to inherit common methods
    angular.extend(this, $controller('FormServerController', { server: server, $scope: $scope }));

    var deleteServer = function(server) {
        serverService.deleteServer(server.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/servers');
            });
    };

    $scope.serverName = angular.copy(server.hostName);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(server) {
        serverService.updateServer(server).
            then(function() {
                $scope.serverName = angular.copy(server.hostName);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(server) {
        var params = {
            title: 'Confirm Delete',
            message: 'This action CANNOT be undone. This will permanently delete ' + server.hostName + '. Are you sure you want to delete ' + server.hostName + '?'
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
        modalInstance.result.then(function() {
            deleteServer(server);
        }, function () {
            // do nothing
        });
    };

};

FormEditServerController.$inject = ['server', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'serverService'];
module.exports = FormEditServerController;