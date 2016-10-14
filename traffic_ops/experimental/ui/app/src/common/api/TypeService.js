var TypeService = function(Restangular, locationUtils, messageModel) {

    this.getTypes = function(useInTable) {
        return Restangular.all('types').getList({ useInTable: useInTable });
    };

    this.getType = function(id) {
        return Restangular.one("types", id).get();
    };

    this.createType = function(type) {
        return Restangular.service('types').post(type)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Type created' } ], true);
                    locationUtils.navigateToPath('/admin/types');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateType = function(type) {
        return type.put()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Type updated' } ], false);
            },
            function(fault) {
                messageModel.setMessages(fault.data.alerts, false);
            }
        );
    };

    this.deleteType = function(id) {
        return Restangular.one("types", id).remove()
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

TypeService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = TypeService;