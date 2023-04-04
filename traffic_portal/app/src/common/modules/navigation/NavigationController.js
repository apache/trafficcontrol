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

/**
 * @param {*} $scope
 * @param {*} $state
 * @param {import("angular").ILocationService} $location
 * @param {import("angular").IWindowService} $window
 * @param {import("angular").ITimeoutService} $timeout
 * @param {*} $uibModal
 * @param {import("../../service/utils/PermissionUtils")} permissionUtils
 * @param {import("../../api/AuthService")} authService
 * @param {import("../../api/TrafficPortalService")} trafficPortalService
 * @param {import("../../models/PropertiesModel")} propertiesModel
 * @param {import("../../models/UserModel")} userModel
 */
var NavigationController = function($scope, $state, $location, $window, $timeout, $uibModal, permissionUtils, authService, trafficPortalService, propertiesModel, userModel) {

    $scope.appName = propertiesModel.properties.name;

    $scope.isProd = (propertiesModel.properties.environment) ? propertiesModel.properties.environment.isProd : false;

    $scope.enforceCapabilities = propertiesModel.properties.enforceCapabilities;

    $scope.customMenu = propertiesModel.properties.customMenu;

    $scope.showCacheChecks = propertiesModel.properties.cacheChecks.show;

    $scope.dsRequestsEnabled = propertiesModel.properties.dsRequests.enabled;

    $scope.userLoaded = userModel.loaded;

    $scope.user = userModel.user;

    $scope.monitor = {
        isOpen: false,
        isDisabled: false
    };

    $scope.settings = {
        isOpen: false,
        isDisabled: false
    };

    $scope.hasCapability = cap => permissionUtils.hasCapability(cap);

    $scope.navigateToPath = function(path) {
        $location.url(path);
    };

    $scope.isState = function(state) {
        return $state.current.name.indexOf(state) !== -1;
    };

    $scope.logout = function() {
        authService.logout();
    };

    $scope.popout = function() {
        $window.open(
            $location.absUrl(),
            '_blank'
        );
    };

    $scope.releaseVersion = function() {
        trafficPortalService.getReleaseVersionInfo()
            .then(function(result) {
                $uibModal.open({
                    templateUrl: 'common/modules/release/release.tpl.html',
                    controller: 'ReleaseController',
                    size: 'sm',
                    resolve: {
                        releaseParams: function () {
                            return result.data;
                        }
                    }
                });
            });
    };

    $scope.customURL = function(item) {
        var url;
        if (item.embed) {
            url = '/#!/custom?url=' + encodeURIComponent(item.url);
        } else {
            url = item.url;
        }
        return url;
    };

    $scope.customTarget = function(item) {
        return (item.embed) ? '_self' : '_blank';
    };

    var explodeMenu = function() {
        var isBig = $('body').hasClass('nav-md');

        $('.side-menu-category ul').slideUp();
        $('.side-menu-category').removeClass('active');
        $('.side-menu-category').removeClass('active-sm');

        if (isBig) {
            $('.current-page').parent('ul').slideDown().parent().addClass('active');
        } else {
            $('.current-page').closest('.side-menu-category').addClass('active-sm');
        }
    };

    var registerMenuHandlers = function() {
        $('.side-menu-category').click(function() {
            var isBig = $('body').hasClass('nav-md');
            if (isBig) {
                if ($(this).is('.active')) {
                    $(this).removeClass('active');
                    $('ul', this).slideUp();
                    $(this).removeClass('nv');
                    $(this).addClass('vn');
                } else {
                    $('#sidebar-menu li ul').slideUp();
                    $(this).removeClass('vn');
                    $(this).addClass('nv');
                    $('ul', this).slideDown();
                    $('#sidebar-menu li').removeClass('active');
                    $(this).addClass('active');
                }
            } else {
                $('#sidebar-menu li ul').slideUp();
                if ($(this).is('.active-sm')) {
                    $(this).removeClass('active-sm');
                    $(this).removeClass('nv');
                    $(this).addClass('vn');
                } else {
                    $(this).removeClass('vn');
                    $(this).addClass('nv');
                    $('ul', this).slideDown();
                    $('#sidebar-menu li').removeClass('active-sm');
                    $(this).addClass('active-sm');
                }
            }
        });

        $('.side-menu-category-item').click(function(event) {
            event.stopPropagation();
            var isBig = $('body').hasClass('nav-md');
            if (!isBig) {
                // close the menu when child is clicked only in small mode
                $(event.currentTarget).closest('.child_menu').slideUp();
            }
        });
    };

    $scope.$on('HeaderController::navExpanded', function() {
        explodeMenu();
    });

    var init = function() {
        $timeout(function() {
            explodeMenu();
            registerMenuHandlers();
        }, 100);
    };
    init();

};

NavigationController.$inject = ['$scope', '$state', '$location', '$window', '$timeout', '$uibModal', 'permissionUtils', 'authService', 'trafficPortalService', 'propertiesModel', 'userModel'];
module.exports = NavigationController;
