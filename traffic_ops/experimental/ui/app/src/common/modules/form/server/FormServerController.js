var FormServerController = function(server, $scope, $uibModal, $anchorScroll, formUtils, locationUtils, cacheGroupService, cdnService, locationService, profileService, serverService, statusService, typeService) {

    var deleteServer = function(server) {
        serverService.deleteServer(server.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/servers');
            });
    };

    var getLocations = function() {
        locationService.getLocations()
            .then(function(result) {
                $scope.locations = result;
            });
    };

    var getCacheGroups = function() {
        cacheGroupService.getCacheGroups()
            .then(function(result) {
                $scope.cacheGroups = result;
            });
    };

    var getTypes = function() {
        typeService.getTypes()
            .then(function(result) {
                $scope.types = result;
            });
    };

    var getCDNs = function() {
        cdnService.getCDNs()
            .then(function(result) {
                $scope.cdns = result;
            });
    };

    var getStatuses = function() {
        statusService.getStatuses()
            .then(function(result) {
                $scope.statuses = result;
            });
    };

    var getProfiles = function() {
        profileService.getProfiles()
            .then(function(result) {
                $scope.profiles = result;
            });
    };

    $scope.serverCopy = angular.copy(server);

    $scope.server = server;

    $scope.props = [
        { name: 'id', type: 'number', required: true, readonly: true },
        { name: 'hostName', type: 'text', required: true, maxLength: 45 },
        { name: 'domainName', type: 'text', required: true, maxLength: 45 },
        { name: 'tcpPort', type: 'number', required: false, maxLength: 10 },
        { name: 'xmppId', type: 'text', required: false, maxLength: 256 },
        { name: 'xmppPasswd', type: 'text', required: false, maxLength: 45 },
        { name: 'interfaceName', type: 'text', required: true, maxLength: 45 },
        { name: 'ipAddress', type: 'text', required: true, maxLength: 45 },
        { name: 'ipNetmask', type: 'text', required: true, maxLength: 45 },
        { name: 'ipGateway', type: 'text', required: true, maxLength: 45 },
        { name: 'ip6Address', type: 'text', required: false, maxLength: 50 },
        { name: 'ip6Gateway', type: 'text', required: false, maxLength: 50 },
        { name: 'interfaceMtu', type: 'number', required: true, maxLength: 11, pattern: new RegExp('(^1500$|^9000$)'), invalidMsg: '1500 or 9000' },
        { name: 'rack', type: 'text', required: false, maxLength: 64 },
        { name: 'mgmtIpAddress', type: 'text', required: false, maxLength: 50 },
        { name: 'mgmtIpNetmask', type: 'text', required: false, maxLength: 45 },
        { name: 'mgmtIpGateway', type: 'text', required: false, maxLength: 45 },
        { name: 'iloIpAddress', type: 'text', required: false, maxLength: 45 },
        { name: 'iloIpNetmask', type: 'text', required: false, maxLength: 45 },
        { name: 'iloIpGateway', type: 'text', required: false, maxLength: 45 },
        { name: 'iloUsername', type: 'text', required: false, maxLength: 45 },
        { name: 'iloPassword', type: 'text', required: false, maxLength: 45 },
        { name: 'routerHostName', type: 'text', required: false, maxLength: 256 },
        { name: 'routerPortName', type: 'text', required: false, maxLength: 256 }
    ];

    $scope.update = function(server) {
        serverService.updateServer(server).
            then(function() {
                $scope.serverCopy = angular.copy(server);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(server) {
        var params = {
            title: 'Confirm Delete',
            message: 'This action CANNOT be undone. This will permanently delete ' + server.hostName + '. Are you sure you want to delete ' + server.hostName + '?'
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
            controller: 'DialogConfirmController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            deleteServer(server);
        }, function () {
            // do nothing
        });
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getLocations();
        getCacheGroups();
        getTypes();
        getCDNs();
        getStatuses();
        getProfiles();
    };
    init();

};

FormServerController.$inject = ['server', '$scope', '$uibModal', '$anchorScroll', 'formUtils', 'locationUtils', 'cacheGroupService', 'cdnService', 'locationService', 'profileService', 'serverService', 'statusService', 'typeService'];
module.exports = FormServerController;