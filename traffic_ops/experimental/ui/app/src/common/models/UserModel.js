var UserModel = function($rootScope, $window, jwtHelper) {

    this.userId = angular.isDefined($window.sessionStorage.token) ? jwtHelper.decodeToken($window.sessionStorage.token)['userid'] : 0;

    this.user = {
        loaded: false
    };

    var removeToken = function() {
        $window.sessionStorage.removeItem('token');
    };

    this.setToken = function(token) {
        $window.sessionStorage.token = token;
        this.userId = jwtHelper.decodeToken(token)['userid'];
    };

    this.setUser = function(userData) {
        this.user.loaded = true;
        this.user = angular.extend(this.user, userData);
        $rootScope.$broadcast('userModel::userUpdated', this.user);
    };

    this.resetUser = function() {
        removeToken();
        this.userId = 0;
        this.user = {
            loaded: false
        };
        $rootScope.$broadcast('userModel::userUpdated', this.user);
    };

};

UserModel.$inject = ['$rootScope', '$window', 'jwtHelper'];
module.exports = UserModel;