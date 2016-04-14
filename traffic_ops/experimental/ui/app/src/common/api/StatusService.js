var StatusService = function(Restangular, messageModel) {

    this.getStatuses = function() {
        return Restangular.all('status').getList();
    };

    this.getStatus = function(id) {
        return Restangular.one("status", id).get();
    };

    this.createStatus = function(status) {
        return Restangular.service('status').post(status)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Status created' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Status create failed' } ], false);
                }
            );
    };

    this.updateStatus = function(status) {
        return status.put()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Status updated' } ], false);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'Status update failed' } ], false);
            }
        );
    };

    this.deleteStatus = function(id) {
        return Restangular.one("status", id).remove()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Status deleted' } ], true);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'Status delete failed' } ], false);
            }
        );
    };

};

StatusService.$inject = ['Restangular', 'messageModel'];
module.exports = StatusService;