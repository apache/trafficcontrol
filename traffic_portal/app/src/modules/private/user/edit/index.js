module.exports = angular.module('trafficPortal.user.edit', [])
    .controller('UserEditController', require('./UserEditController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.user.edit', {
                url: '',
                views: {
                    userContent: {
                        templateUrl: 'modules/private/user/edit/user.edit.tpl.html',
                        controller: 'UserEditController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
