var UserModel = function($rootScope, $window, messageModel) {

    var user = {};
    user.loaded = false;
    this.user = user;

    this.setUser = function(userData) {
        user.loaded = true;
        user = angular.extend(user, userData);
        if (user.newUser) {
            user.username = ''; // new users were given a temp username that needs to be ditched
        }
        if (!user.localUser) {
            messageModel.setMessages([ { level: 'success', text: 'Logged in as LDAP user.' } ], false);
        }
        $rootScope.$broadcast('userModel::userUpdated', user);
    };

    this.resetUser = function() {
        $window.sessionStorage.removeItem('token');

        user = {};
        user.loaded = false;
        this.user = user;
        $rootScope.$broadcast('userModel::userUpdated', user);
    };

};

UserModel.$inject = ['$rootScope', '$window', 'messageModel'];
module.exports = UserModel;