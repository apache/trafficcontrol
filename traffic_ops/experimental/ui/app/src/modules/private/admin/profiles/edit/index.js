module.exports = angular.module('trafficOps.private.admin.profiles.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.profiles.edit', {
                url: '/{profileId}/edit',
                views: {
                    profilesContent: {
                        templateUrl: 'common/modules/form/profile/form.profile.tpl.html',
                        controller: 'FormEditProfileController',
                        resolve: {
                            profile: function($stateParams, profileService) {
                                return profileService.getProfile($stateParams.profileId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
