var ApplicationService = function($rootScope, $anchorScroll, $http) {

    this.startup = function() {
        // anything you need to do at startup
    };

    $rootScope.$on("$viewContentLoaded", function() {
        $anchorScroll(); // scrolls window to top
    });

    var init = function() {
        $http.defaults.withCredentials = true;
    };
    init();

};

ApplicationService.$inject = ['$rootScope', '$anchorScroll', '$http'];
module.exports = ApplicationService;
