module.exports = angular.module('trafficPortal.collateral', [])
    .controller('CollateralController', require('./CollateralController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.collateral', {
                url: 'collateral',
                views: {
                    privateContent: {
                        templateUrl: 'modules/private/collateral/collateral.tpl.html',
                        controller: 'CollateralController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
