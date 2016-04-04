var ParameterService = function(Restangular, messageModel) {

    this.getParameters = function() {
        return Restangular.all('parameter').getList();
    };

    this.getParameter = function(id) {
        return Restangular.one("parameter", id).get();
    };

    this.createParameter = function(parameter) {
        return Restangular.service('parameter').post(parameter)
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Parameter created' } ], true);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'Parameter create failed' } ], false);
            }
        );
    };

    this.updateParameter = function(parameter) {
        return parameter.put()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Parameter updated' } ], false);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'Parameter update failed' } ], false);
            }
        );
    };

    this.deleteParameter = function(id) {
        return Restangular.one("parameter", id).remove()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Parameter deleted' } ], true);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'Parameter delete failed' } ], false);
            }
        );
    };

};

ParameterService.$inject = ['Restangular', 'messageModel'];
module.exports = ParameterService;