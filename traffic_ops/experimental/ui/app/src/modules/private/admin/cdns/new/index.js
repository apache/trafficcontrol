module.exports = angular.module('trafficOps.private.admin.cdns.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.cdns.new', {
                url: '/new',
                views: {
                    cdnsContent: {
                        templateUrl: 'common/modules/form/cdn/form.cdn.tpl.html',
                        controller: 'FormNewCDNController',
                        resolve: {
                            cdn: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
