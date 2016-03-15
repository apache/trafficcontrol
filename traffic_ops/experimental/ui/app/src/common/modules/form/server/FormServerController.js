var FormServerController = function(server, $scope, formUtils, locationUtils, serverService) {

    $scope.serverOriginal = server;

    $scope.server = angular.copy(server);

    $scope.props = [
        { name: 'id', required: true, readonly: true },
        { name: 'hostName', required: true, maxLength: 45 },
        { name: 'domainName', required: true, maxLength: 45 },
        { name: 'tcpPort', required: false, maxLength: 10 },
        { name: 'xmppId', required: false, maxLength: 256 },
        { name: 'xmppPasswd', required: false, maxLength: 45 },
        { name: 'interfaceName', required: true, maxLength: 45 },
        { name: 'ipAddress', required: true, maxLength: 45 },
        { name: 'ipNetmask', required: true, maxLength: 45 },
        { name: 'ipGateway', required: true, maxLength: 45 },
        { name: 'ip6Address', required: false, maxLength: 50 },
        { name: 'ip6Gateway', required: false, maxLength: 50 },
        { name: 'interfaceMtu', required: true, maxLength: 11, pattern: new RegExp('(^1500$|^9000$)'), invalidMsg: '1500 or 9000' },
        { name: 'rack', required: false, maxLength: 64 },
        { name: 'mgmtIpAddress', required: false, maxLength: 50 },
        { name: 'mgmtIpNetmask', required: false, maxLength: 45 },
        { name: 'mgmtIpGateway', required: false, maxLength: 45 },
        { name: 'iloIpAddress', required: false, maxLength: 45 },
        { name: 'iloIpNetmask', required: false, maxLength: 45 },
        { name: 'iloIpGateway', required: false, maxLength: 45 },
        { name: 'iloUsername', required: false, maxLength: 45 },
        { name: 'iloPassword', required: false, maxLength: 45 },
        { name: 'routerHostName', required: false, maxLength: 256 },
        { name: 'routerPortName', required: false, maxLength: 256 }
    ];

    $scope.embeds = [
        { name: 'phys_location', required: true, maxLength: 11 },
        { name: 'cachegroup', required: true, maxLength: 11 },
        { name: 'type', required: true, maxLength: 11 },
        { name: 'status', required: true, maxLength: 11 },
        { name: 'profile', required: true, maxLength: 11 },
        { name: 'cdn', required: true, maxLength: 11 }
    ];

    $scope.update = function(server) {
        alert('implement update');
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormServerController.$inject = ['server', '$scope', 'formUtils', 'locationUtils', 'serverService'];
module.exports = FormServerController;