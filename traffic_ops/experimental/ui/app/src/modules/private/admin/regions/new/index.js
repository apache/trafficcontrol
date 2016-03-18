module.exports = angular.module('trafficOps.private.admin.regions.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.regions.new', {
                url: '/new',
                views: {
                    regionsContent: {
                        templateUrl: 'common/modules/form/region/form.region.tpl.html',
                        controller: 'FormNewRegionController',
                        resolve: {
                            region: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
