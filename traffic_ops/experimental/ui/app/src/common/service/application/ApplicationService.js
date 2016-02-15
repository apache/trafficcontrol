var ApplicationService = function($http) {

    this.startup = function() {
        // anything you need to do at startup
    };

    var init = function() {
        $http.defaults.withCredentials = true;
    };
    init();

};

ApplicationService.$inject = ['$http'];
module.exports = ApplicationService;
