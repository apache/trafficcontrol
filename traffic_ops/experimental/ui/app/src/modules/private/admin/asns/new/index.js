module.exports = angular.module('trafficOps.private.admin.asns.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.asns.new', {
                url: '/new',
                views: {
                    asnsContent: {
                        templateUrl: 'common/modules/form/asn/form.asn.tpl.html',
                        controller: 'FormNewASNController',
                        resolve: {
                            asn: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
