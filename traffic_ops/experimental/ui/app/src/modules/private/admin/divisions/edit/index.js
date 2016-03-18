module.exports = angular.module('trafficOps.private.admin.divisions.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.divisions.edit', {
                url: '/{divisionId}/edit',
                views: {
                    divisionsContent: {
                        templateUrl: 'common/modules/form/division/form.division.tpl.html',
                        controller: 'FormEditDivisionController',
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
