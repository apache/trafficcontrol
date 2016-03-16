module.exports = angular.module('trafficOps.private.admin.users', [])
    .controller('UsersController', require('./UsersController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.users', {
                url: '/users',
                abstract: true,
                views: {
                    adminContent: {
                        templateUrl: 'modules/private/admin/users/users.tpl.html',
                        controller: 'UsersController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
