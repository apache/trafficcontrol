module.exports = angular.module('trafficOps.private.user', [])
    .controller('UserController', require('./UserController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.user', {
                url: 'user',
                abstract: true,
                views: {
                    privateContent: {
                        templateUrl: 'modules/private/user/user.tpl.html',
                        controller: 'UserController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
