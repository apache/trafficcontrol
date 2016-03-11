module.exports = angular.module('trafficOps.private.configure.divisions.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.divisions.edit', {
                url: '/{divisionId}',
                views: {
                    divisionsContent: {
                        templateUrl: 'common/modules/form/division/form.division.tpl.html',
                        controller: 'FormDivisionController',
                        resolve: {
                            division: function($stateParams, divisionService) {
                                return divisionService.getDivision($stateParams.divisionId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
module.exports = angular.module('trafficOps.private.configure.divisions.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.divisions.edit', {
                url: '/{divisionId}',
                views: {
                    divisionsContent: {
                        templateUrl: 'common/modules/form/division/form.division.tpl.html',
                        controller: 'FormDivisionController',
                        resolve: {
                            division: function($stateParams, divisionService) {
                                return divisionService.getDivision($stateParams.divisionId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
