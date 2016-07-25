var AboutController = function($scope, $sce, $timeout, propertiesModel) {

    var pinIframe = function() {
        var headerHeight = $('#header').css("height"),
            footerHeight = $('#footer').css("height");

        $('#aboutFrameWrapper').css("top", headerHeight);
        $('#aboutFrameWrapper').css("bottom", footerHeight);
    };

    $scope.properties = propertiesModel.properties;

    $scope.trustSrc = function(src) {
        return $sce.trustAsResourceUrl(src);
    };

    var init = function () {
        $timeout(function () {
            pinIframe();
        }, 200);
    };
    init();

};

AboutController.$inject = ['$scope', '$sce', '$timeout', 'propertiesModel'];
module.exports = AboutController;