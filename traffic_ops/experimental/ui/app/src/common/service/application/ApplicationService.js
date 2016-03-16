var ApplicationService = function($rootScope, $anchorScroll) {

    this.startup = function() {
        // anything you need to do at startup
    };

    $rootScope.$on("$viewContentLoaded", function() {
        $anchorScroll(); // scrolls window to top
    });

};

ApplicationService.$inject = ['$rootScope', '$anchorScroll'];
module.exports = ApplicationService;
