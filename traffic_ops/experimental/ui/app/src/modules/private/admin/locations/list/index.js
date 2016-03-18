module.exports = angular.module('trafficOps.private.admin.locations.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.locations.list', {
                url: '',
                views: {
                    locationsContent: {
                        templateUrl: 'common/modules/table/locations/table.locations.tpl.html',
                        controller: 'TableLocationsController',
                        resolve: {
                            locations: function(locationService) {
                                return locationService.getLocations();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
