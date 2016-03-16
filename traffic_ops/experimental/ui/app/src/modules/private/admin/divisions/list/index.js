module.exports = angular.module('trafficOps.private.admin.divisions.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.divisions.list', {
                url: '',
                views: {
                    divisionsContent: {
                        templateUrl: 'common/modules/table/divisions/table.divisions.tpl.html',
                        controller: 'TableDivisionsController',
                        resolve: {
                            divisions: function(divisionService) {
                                return divisionService.getDivisions();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
