module.exports = angular.module('trafficOps.private.admin.statuses', [])
    .controller('StatusesController', require('./StatusesController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.statuses', {
                url: '/statuses',
                abstract: true,
                views: {
                    adminContent: {
                        templateUrl: 'modules/private/admin/statuses/statuses.tpl.html',
                        controller: 'StatusesController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
