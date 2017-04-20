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

var HeaderController = function($rootScope, $scope, $log, $state, $anchorScroll, $interval, locationUtils, authService, changeLogService, trafficOpsService, userModel) {

    $scope.isCollapsed = true;

    $scope.userLoaded = userModel.loaded;

    /* we don't want real time changes to the user showing up. we want the ability to revert changes
    if necessary. thus, we will only update this on save. see userModel::userUpdated event below.
     */
    $scope.user = angular.copy(userModel.user);

    $scope.newLogCount = 0;

    $scope.changeLogs = [];

    $scope.isState = function(state) {
        return $state.current.name.indexOf(state) !== -1;
    };

    $scope.getChangeLogs = function() {
        $scope.changeLogs = [];
        changeLogService.getChangeLogs({ limit: 6 })
            .then(function(response) {
                $scope.changeLogs = response;
            });
    };

    $scope.getRelativeTime = function(date) {
        return moment(date).fromNow();
    };

    $scope.logout = function() {
        authService.logout();
    };

    $scope.dumpDB = function() {
        trafficOpsService.dumpDB();
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    var changeLogInterval,
        increment = 1;

    var createChangeLogInterval = function() {
        killChangeLogInterval();
        changeLogInterval = $interval(function() { getNewLogCount() }, (increment*60*1000)); // every X minutes
    };

    var killChangeLogInterval = function() {
        if (angular.isDefined(changeLogInterval)) {
            $interval.cancel(changeLogInterval);
            changeLogInterval = undefined;
        }
    };

    var getNewLogCount = function() {
        changeLogService.getNewLogCount()
            .then(function(result) {
                $scope.newLogCount = result.data.response.newLogcount;
            });
    };

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

    $scope.$on('changeLogService::getChangeLogs', function() {
        $scope.newLogCount = 0;
    });

    var init = function () {
        scrollToTop();
        initToggleMenu();
        getNewLogCount();
        createChangeLogInterval();
    };
    init();
};

HeaderController.$inject = ['$rootScope', '$scope', '$log', '$state', '$anchorScroll', '$interval', 'locationUtils', 'authService', 'changeLogService', 'trafficOpsService', 'userModel'];
module.exports = HeaderController;
