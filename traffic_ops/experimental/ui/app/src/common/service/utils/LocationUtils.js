var LocationUtils = function($location) {

    this.navigateToPath = function(path) {
        $location.url(path);
    };

};

LocationUtils.$inject = ['$location'];
module.exports = LocationUtils;
