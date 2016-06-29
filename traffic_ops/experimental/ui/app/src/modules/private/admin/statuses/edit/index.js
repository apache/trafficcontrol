module.exports = angular.module('trafficOps.private.admin.statuses.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.statuses.edit', {
                url: '/{statusId}/edit',
                views: {
                    statusesContent: {
                        templateUrl: 'common/modules/form/status/form.status.tpl.html',
                        controller: 'FormEditStatusController',
                        resolve: {
                            status: function($stateParams, statusService) {
                                return statusService.getStatus($stateParams.statusId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
