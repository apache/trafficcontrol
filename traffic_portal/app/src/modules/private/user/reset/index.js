module.exports = angular.module('trafficPortal.user.reset', [])
    .controller('UserResetController', require('./UserResetController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.user.reset', {
                url: '/reset',
                views: {
                    userContent: {
                        templateUrl: 'modules/private/user/edit/user.edit.tpl.html',
                        controller: 'UserResetController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
