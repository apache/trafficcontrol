module.exports = angular.module('trafficOps.private.admin.cdns', [])
    .controller('CdnsController', require('./CdnsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.cdns', {
                url: '/cdns',
                abstract: true,
                views: {
                    adminContent: {
                        templateUrl: 'modules/private/admin/cdns/cdns.tpl.html',
                        controller: 'CdnsController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
