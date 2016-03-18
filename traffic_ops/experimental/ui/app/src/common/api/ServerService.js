var ServerService = function(Restangular, messageModel) {

    this.getServers = function() {
        return Restangular.all('server').getList();
    };

    this.getServer = function(id) {
        return Restangular.one("server", id).get();
    };

    this.createServer = function(server) {
        return Restangular.service('server').post(server)
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Server created' } ], true);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'Server create failed' } ], false);
            }
        );
    };

    this.updateServer = function(server) {
        return server.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Server updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Server update failed' } ], false);
                }
            );
    };

    this.deleteServer = function(id) {
        return Restangular.one("server", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Server deleted' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Server delete failed' } ], false);
                }
            );
    };

};

ServerService.$inject = ['Restangular', 'messageModel'];
module.exports = ServerService;