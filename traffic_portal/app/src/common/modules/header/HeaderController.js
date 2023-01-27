/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/** @typedef {import("jquery")} $ */

const defaultBannerColor = "#EDEDED";
const defaultSidebarColor = "#2A3F54";
const defaultTextColor = "#515356";

const prodTextColor = "white";
const prodBannerColor = "#B22222";

/** @typedef {import("moment")} moment */

/**
 * @param {import("angular").IRootScopeService} $rootScope
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").ILocationService} $location
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../service/utils/PermissionUtils")} permissionUtils
 * @param {import("../../api/AuthService")} authService
 * @param {import("../../api/TrafficPortalService")} trafficPortalService
 * @param {import("../../api/ChangeLogService")} changeLogService
 * @param {import("../../api/CDNService")} cdnService
 * @param {import("../../models/ChangeLogModel")} changeLogModel
 * @param {import("../../models/UserModel")} userModel
 * @param {import("../../models/PropertiesModel")} propertiesModel
 */
var HeaderController = function($rootScope, $scope, $state, $uibModal, $location, $anchorScroll, locationUtils, permissionUtils, authService, trafficPortalService, changeLogService, cdnService, changeLogModel, userModel, propertiesModel) {

    let getCDNs = function(notifications) {
        cdnService.getCDNs(true)
            .then(function(cdns) {
                cdns.forEach(function(cdn) {
                    cdn.hasNotifications = notifications.find(function(notification){ return cdn.name === notification.cdn });
                });
                $scope.cdns = cdns;
            });
    };

    $scope.isCollapsed = true;

    $scope.userLoaded = userModel.loaded;

    $scope.enviroName = (propertiesModel.properties.environment) ? propertiesModel.properties.environment.name : '';

    if (propertiesModel.properties.environment && propertiesModel.properties.environment.isProd) {
        document.documentElement.style.setProperty("--banner-color", prodBannerColor);
        document.documentElement.style.setProperty("--sidebar-color", prodBannerColor);
        document.documentElement.style.setProperty("--banner-text-color", prodTextColor);
    } else {
        document.documentElement.style.setProperty("--banne-color", defaultBannerColor);
        document.documentElement.style.setProperty("--sidebar-color", defaultSidebarColor);
        document.documentElement.style.setProperty("--banner-text-color", defaultTextColor);
    }

    /* we don't want real time changes to the user showing up. we want the ability to revert changes
    if necessary. thus, we will only update this on save. see userModel::userUpdated event below.
     */
    $scope.user = angular.copy(userModel.user);

    $scope.newLogCount = changeLogModel.newLogCount;

    $scope.changeLogs = [];

    $scope.hasCapability = cap => permissionUtils.hasCapability(cap);

    $scope.isState = function(state) {
        return $state.current.name.indexOf(state) !== -1;
    };

    $scope.getChangeLogs = function() {
        $scope.loadingChangeLogs = true;
        $scope.changeLogs = [];
        changeLogService.getChangeLogs({ limit: 6 })
            .then(function(response) {
                $scope.loadingChangeLogs = false;
                $scope.changeLogs = response;
            });
    };

    $scope.getNotifications = function(cdn) {
        $scope.loadingNotifications = true;
        $scope.notifications = [];
        cdnService.getNotifications({ cdn: cdn.name })
            .then(function(response) {
                $scope.loadingNotifications = false;
                $scope.notifications = response;
            });
    };

    $scope.getRelativeTime = function(date) {
        return moment(date).fromNow();
    };

    $scope.logout = function() {
        authService.logout();
    };

    $scope.dbDump = function() {
        trafficPortalService.dbDump();
    };

    $scope.confirmQueueServerUpdates = function() {
        var params = {
            title: 'Queue Server Updates',
            message: "Please select a CDN"
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
            controller: 'DialogSelectController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                },
                collection: function() {
                    return $scope.cdns.filter(function(cdn) {
                        return cdn.name != 'ALL';
                    });
                }
            }
        });
        modalInstance.result.then(function(cdn) {
            cdnService.queueServerUpdates(cdn.id);
        }, function () {
            // do nothing
        });
    };

    $scope.lockCDN = function() {
        const modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/select/lock/dialog.select.lock.tpl.html',
            controller: 'DialogSelectLockController',
            size: 'md',
            resolve: {
                cdns: function() {
                    return $scope.cdns;
                }
            }
        });
        modalInstance.result.then(function(lock) {
            cdnService.createLock(lock).
            then(
                function() {
                    $state.reload();
                }
            );
        }, function () {
            // do nothing
        });
    };

    $scope.snapshot = function() {
        var params = {
            title: 'Diff CDN Config Snapshot',
            message: "Please select a CDN"
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
            controller: 'DialogSelectController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                },
                collection: function() {
                    return $scope.cdns.filter(function(cdn) {
                        return cdn.name != 'ALL';
                    });
                }
            }
        });
        modalInstance.result.then(function(cdn) {
            $location.path('/cdns/' + cdn.id + '/config/changes');
        }, function () {
            // do nothing
        });
    };

    $scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

    var scrollToTop = function() {
        $anchorScroll(); // hacky?
    };

    var initToggleMenu = function() {
        $('#menu_toggle').click(function () {
            var isBig = $('body').hasClass('nav-md');
            if (isBig) {
                // shrink side menu
                $('body').removeClass('nav-md');
                $('body').addClass('nav-sm');
                $('.main-nav').removeClass('scroll-view');
                $('.main-nav').removeAttr('style');
                $('.sidebar-footer').hide();

                if ($('#sidebar-menu li').hasClass('active')) {
                    $('#sidebar-menu li.active').addClass('active-sm');
                    $('#sidebar-menu li.active').removeClass('active');
                }

                $('.side-menu-category ul').hide();

            } else {
                // expand side menu
                $('body').removeClass('nav-sm');
                $('body').addClass('nav-md');
                $('.sidebar-footer').show();

                if ($('#sidebar-menu li').hasClass('active-sm')) {
                    $('#sidebar-menu li.active-sm').addClass('active');
                    $('#sidebar-menu li.active-sm').removeClass('active-sm');
                }

                $rootScope.$broadcast('HeaderController::navExpanded', {});

            }
        });
    };

    $scope.$on('userModel::userUpdated', function() {
        $scope.user = angular.copy(userModel.user);
    });

    $rootScope.$on('notificationsController::refreshNotifications', function(event, options) {
        getCDNs(options.notifications);
    });

    var init = function () {
        scrollToTop();
        initToggleMenu();
    };
    init();
};

HeaderController.$inject = ['$rootScope', '$scope', '$state', '$uibModal', '$location', '$anchorScroll', 'locationUtils', 'permissionUtils', 'authService', 'trafficPortalService', 'changeLogService', 'cdnService', 'changeLogModel', 'userModel', 'propertiesModel'];
module.exports = HeaderController;
