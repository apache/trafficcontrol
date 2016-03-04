module.exports = angular.module('trafficOps.private.configure.divisions.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.divisions.list', {
                url: '',
                views: {
                    divisionsContent: {
                        templateUrl: 'common/modules/table/divisions/table.divisions.tpl.html',
                        controller: 'TableDivisionsController',
                        resolve: {
                            divisions: function(divisionService, ENV) {
                                return divisionService.getDivisions(ENV.api['root'] + 'division');
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
