module.exports = angular.module('trafficOps.private.administer', [])
    .controller('AdministerController', require('./AdministerController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.administer', {
                url: 'administer',
                abstract: true,
                views: {
                    privateContent: {
                        templateUrl: 'modules/private/administer/administer.tpl.html',
                        controller: 'AdministerController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
