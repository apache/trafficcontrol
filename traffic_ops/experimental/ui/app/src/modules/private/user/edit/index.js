module.exports = angular.module('trafficOps.private.user.edit', [])
    .controller('UserEditController', require('./UserEditController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.user.edit', {
                url: '',
                views: {
                    userContent: {
                        templateUrl: 'common/modules/form/user/form.user.tpl.html',
                        controller: 'UserEditController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });

