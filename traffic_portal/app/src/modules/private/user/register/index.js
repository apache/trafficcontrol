module.exports = angular.module('trafficPortal.user.register', [])
    .controller('UserRegisterController', require('./UserRegisterController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.user.register', {
                url: '/register',
                views: {
                    userContent: {
                        templateUrl: 'modules/private/user/edit/user.edit.tpl.html',
                        controller: 'UserRegisterController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
