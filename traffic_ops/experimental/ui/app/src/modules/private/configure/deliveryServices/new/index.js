module.exports = angular.module('trafficOps.private.configure.deliveryServices.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.deliveryServices.new', {
                url: '/new',
                views: {
                    deliveryServicesContent: {
                        templateUrl: 'common/modules/form/deliveryService/form.deliveryService.tpl.html',
                        controller: 'FormNewDeliveryServiceController',
                        resolve: {
                            deliveryService: function() {
                                return {
                                    active: false,
                                    signed: false,
                                    qstringIgnore: "0",
                                    dscp: "0",
                                    geoLimit: "0",
                                    geoProvider: "0",
                                    initialDispersion: "1",
                                    ipv6RoutingEnabled: false,
                                    rangeRequestHandling: "0",
                                    multiSiteOrigin: false,
                                    regionalGeoBlocking: false,
                                    logsEnabled: false
                                };
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
