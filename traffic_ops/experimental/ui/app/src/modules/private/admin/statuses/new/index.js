module.exports = angular.module('trafficOps.private.admin.statuses.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.statuses.new', {
                url: '/new',
                views: {
                    statusesContent: {
                        templateUrl: 'common/modules/form/status/form.status.tpl.html',
                        controller: 'FormNewStatusController',
                        resolve: {
                            status: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
