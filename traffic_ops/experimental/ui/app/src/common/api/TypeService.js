var TypeService = function(Restangular, messageModel) {

    this.getTypes = function() {
        return Restangular.all('type').getList();
    };

    this.getType = function(id) {
        return Restangular.one("type", id).get();
    };

    this.createType = function(type) {
        return Restangular.service('type').post(type)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Type created' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Type create failed' } ], false);
                }
            );
    };

    this.updateType = function(type) {
        return type.put()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Type updated' } ], false);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'Type update failed' } ], false);
            }
        );
    };

    this.deleteType = function(id) {
        return Restangular.one("type", id).remove()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Type deleted' } ], true);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'Type delete failed' } ], false);
            }
        );
    };

};

TypeService.$inject = ['Restangular', 'messageModel'];
module.exports = TypeService;