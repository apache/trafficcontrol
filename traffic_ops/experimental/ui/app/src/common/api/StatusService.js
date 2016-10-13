var StatusService = function(Restangular, locationUtils, messageModel) {

    this.getStatuses = function() {
        return Restangular.all('statuses').getList();
    };

    this.getStatus = function(id) {
        return Restangular.one("statuses", id).get();
    };

    this.createStatus = function(status) {
        return Restangular.service('statuses').post(status)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Status created' } ], true);
                    locationUtils.navigateToPath('/admin/statuses');
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
        return Restangular.one("statuses", id).remove()
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

StatusService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = StatusService;