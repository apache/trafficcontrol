var ServerService = function(Restangular, locationUtils, messageModel) {

    this.getServers = function() {
        return Restangular.all('servers').getList();
    };

    this.getServer = function(id) {
        return Restangular.one("servers", id).get();
    };

    this.createServer = function(server) {
        return Restangular.service('servers').post(server)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Server created' } ], true);
                    locationUtils.navigateToPath('/configure/servers');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateServer = function(server) {
        return server.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Server updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteServer = function(id) {
        return Restangular.one("servers", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Server deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

};

ServerService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = ServerService;