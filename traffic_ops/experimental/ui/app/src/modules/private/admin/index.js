module.exports = angular.module('trafficOps.private.admin', [])
    .controller('AdminController', require('./AdminController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin', {
                url: 'admin',
                abstract: true,
                views: {
                    privateContent: {
                        templateUrl: 'modules/private/admin/admin.tpl.html',
                        controller: 'AdminController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
