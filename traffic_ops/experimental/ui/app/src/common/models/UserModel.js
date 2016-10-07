var UserModel = function($rootScope) {

    this.loaded = false;

    this.user = {};

    this.setUser = function(userData) {
        this.loaded = true;
        this.user = userData;
        $rootScope.$broadcast('userModel::userUpdated', this.user);
    };

    this.resetUser = function() {
        this.loaded = false;
        this.userId = 0;
        this.user = {};
        $rootScope.$broadcast('userModel::userUpdated', this.user);
    };

};

UserModel.$inject = ['$rootScope'];
module.exports = UserModel;