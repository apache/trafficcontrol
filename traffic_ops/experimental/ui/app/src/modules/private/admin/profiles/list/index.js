module.exports = angular.module('trafficOps.private.admin.profiles.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.profiles.list', {
                url: '',
                views: {
                    profilesContent: {
                        templateUrl: 'common/modules/table/profiles/table.profiles.tpl.html',
                        controller: 'TableProfilesController',
                        resolve: {
                            profiles: function(profileService) {
                                return profileService.getProfiles();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
