var FormEditCDNController = function(cdn, $scope, $controller, $uibModal, $anchorScroll, locationUtils, cdnService) {

    // extends the FormCDNController to inherit common methods
    angular.extend(this, $controller('FormCDNController', { cdn: cdn, $scope: $scope }));

    var deleteCDN = function(cdn) {
        cdnService.deleteCDN(cdn.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/cdns');
            });
    };

    $scope.cdnName = angular.copy(cdn.name);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(cdn) {
        cdnService.updateCDN(cdn).
            then(function() {
                $scope.cdnName = angular.copy(cdn.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(cdn) {
        var params = {
            title: 'Delete CDN: ' + cdn.name,
            key: cdn.name
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
            controller: 'DialogDeleteController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            deleteCDN(cdn);
        }, function () {
            // do nothing
        });
    };

};

FormEditCDNController.$inject = ['cdn', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'cdnService'];
module.exports = FormEditCDNController;