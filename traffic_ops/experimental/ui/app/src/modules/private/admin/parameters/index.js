module.exports = angular.module('trafficOps.private.admin.parameters', [])
    .controller('ParametersController', require('./ParametersController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.parameters', {
                url: '/parameters',
                abstract: true,
                views: {
                    adminContent: {
                        templateUrl: 'modules/private/admin/parameters/parameters.tpl.html',
                        controller: 'ParametersController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
