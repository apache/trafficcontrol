var FooterController = function($scope, propertiesModel) {

    var chunk = function(array, n) {
        var retval = [];
        for (var i = 0, len = array.length; i < len; i += n) {
            retval.push(array.slice(i, i + n));
        }
        return retval;
    };

    $scope.footerChunks = chunk(propertiesModel.properties.footer.links, 4);

};

FooterController.$inject = ['$scope', 'propertiesModel'];
module.exports = FooterController;
