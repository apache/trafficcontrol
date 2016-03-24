var FormEditASNController = function(asn, $scope, $controller, $uibModal, $anchorScroll, locationUtils, asnService) {

    // extends the FormASNController to inherit common methods
    angular.extend(this, $controller('FormASNController', { asn: asn, $scope: $scope }));

    var deleteASN = function(asn) {
        asnService.deleteASN(asn.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/asns');
            });
    };

    $scope.asnName = angular.copy(asn.asn);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(asn) {
        asnService.updateASN(asn).
            then(function() {
                $scope.asnName = angular.copy(asn.asn);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(asn) {
        var params = {
            title: 'Confirm Delete',
            message: 'This action CANNOT be undone. This will permanently delete ' + asn.asn + '. Are you sure you want to delete ' + asn.asn + '?'
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
            deleteASN(asn);
        }, function () {
            // do nothing
        });
    };

};

FormEditASNController.$inject = ['asn', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'asnService'];
module.exports = FormEditASNController;