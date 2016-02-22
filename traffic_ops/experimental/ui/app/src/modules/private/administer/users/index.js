module.exports = angular.module('trafficOps.private.administer.users', [])
    .controller('UsersController', require('./UsersController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.administer.users', {
                url: '/users',
                views: {
                    administerContent: {
                        templateUrl: 'modules/private/administer/users/users.tpl.html',
                        controller: 'UsersController'
                    }
                },
                resolve: {
                    users: function() {
                        return [];
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
