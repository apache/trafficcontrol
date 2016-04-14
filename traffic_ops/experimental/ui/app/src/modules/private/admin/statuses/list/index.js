module.exports = angular.module('trafficOps.private.admin.statuses.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.statuses.list', {
                url: '',
                views: {
                    statusesContent: {
                        templateUrl: 'common/modules/table/statuses/table.statuses.tpl.html',
                        controller: 'TableStatusesController',
                        resolve: {
                            statuses: function(statusService) {
                                return statusService.getStatuses();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
