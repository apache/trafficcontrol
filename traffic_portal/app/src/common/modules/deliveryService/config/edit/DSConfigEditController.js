var DSConfigEditController = function(deliveryService, $scope, $uibModalInstance) {

    $scope.deliveryService = deliveryService;

    $scope.close = function () {
        $uibModalInstance.dismiss('close');
    };

    $scope.isHTTP = function(ds) {
        return ds.type.indexOf('HTTP') !== -1;
    };

    $scope.edgeFQDNs = function(ds) {
        var urlString = '';
        if (_.isArray(ds.exampleURLs) && ds.exampleURLs.length > 0) {
            for (var i = 0; i < ds.exampleURLs.length; i++) {
                urlString += ds.exampleURLs[i] + '\n';
            }
        }
        return urlString;
    };

    $scope.rangeRequestHandling = function(ds) {
        var rrh = '';
        if (ds.rangeRequestHandling == '0') {
            rrh = 'Do not cache range requests';
        } else if (ds.rangeRequestHandling == '1') {
            rrh = 'Background fetch';
        } else if (ds.rangeRequestHandling == '2') {
            rrh = 'Cache range requests';
        }
        return rrh;
    };

};

DSConfigEditController.$inject = ['deliveryService', '$scope', '$uibModalInstance'];
module.exports = DSConfigEditController;
