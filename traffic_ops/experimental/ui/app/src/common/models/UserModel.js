var UserModel = function($rootScope, $window, jwtHelper) {

    this.loaded = false;

    this.userId = angular.isDefined($window.sessionStorage.token) ? jwtHelper.decodeToken($window.sessionStorage.token)['userid'] : 0;

    this.user = {};

    var removeToken = function() {
        $window.sessionStorage.removeItem('token');
    };

    this.setToken = function(token) {
        $window.sessionStorage.token = token;
        this.userId = jwtHelper.decodeToken(token)['userid'];
    };

    this.setUser = function(userData) {
        this.loaded = true;
        this.user = userData;
        $rootScope.$broadcast('userModel::userUpdated', this.user);
    };

    this.resetUser = function() {
        removeToken();
        this.loaded = false;
        this.userId = 0;
        this.user = {};
        $rootScope.$broadcast('userModel::userUpdated', this.user);
    };

};

UserModel.$inject = ['$rootScope', '$window', 'jwtHelper'];
module.exports = UserModel;