var HeaderController = function($rootScope, $scope, $log, $state, $anchorScroll, authService, userModel) {

    $scope.isCollapsed = true;

    $scope.userLoaded = userModel.loaded;

    /* we don't want real time changes to the user showing up. we want the ability to revert changes
    if necessary. thus, we will only update this on save. see userModel::userUpdated event below.
     */
    $scope.user = angular.copy(userModel.user);

    $scope.isState = function(state) {
        return $state.current.name.indexOf(state) !== -1;
    };

    $scope.logout = function() {
        authService.logout();
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

    $scope.$on('userModel::userUpdated', function(event) {
        $scope.user = angular.copy(userModel.user);
    });

    var init = function () {
        scrollToTop();
        initToggleMenu();
    };
    init();
};

HeaderController.$inject = ['$rootScope', '$scope', '$log', '$state', '$anchorScroll', 'authService', 'userModel'];
module.exports = HeaderController;
