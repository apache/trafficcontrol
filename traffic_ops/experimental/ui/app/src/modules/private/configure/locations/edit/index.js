module.exports = angular.module('trafficOps.private.configure.locations.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.locations.edit', {
                url: '/{locationId}',
                views: {
                    locationContent: {
                        templateUrl: 'common/modules/form/location/form.location.tpl.html',
                        controller: 'FormLocationController',
                        resolve: {
                            location: function($stateParams, locationService) {
                                return locationService.getLocation($stateParams.locationId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
