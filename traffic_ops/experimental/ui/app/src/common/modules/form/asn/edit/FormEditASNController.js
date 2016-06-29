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
            title: 'Delete ASN: ' + asn.asn,
            key: asn.asn.toString()
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
            deleteASN(asn);
        }, function () {
            // do nothing
        });
    };

};

FormEditASNController.$inject = ['asn', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'asnService'];
module.exports = FormEditASNController;