var ApplicationService = function($rootScope, $http, $window) {

    /**
     *  Application Startup Process
     */
    this.startup = function() {
        registerOnlineListener();
    };

    var registerOnlineListener = function() {
        $rootScope.online = $window.navigator.onLine;
        $window.addEventListener("offline", function () {
            $rootScope.$apply(function() {
                $rootScope.online = false;
            });
        }, false);
        $window.addEventListener("online", function () {
            $rootScope.$apply(function() {
                $rootScope.online = true;
            });
        }, false);
    };

    var init = function() {
        $http.defaults.withCredentials = true;
    };
    init();

};

ApplicationService.$inject = ['$rootScope', '$http', '$window'];
module.exports = ApplicationService;
