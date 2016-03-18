var DivisionService = function(Restangular, messageModel) {

    this.getDivisions = function() {
        return Restangular.all('division').getList();
    };

    this.getDivision = function(id) {
        return Restangular.one("division", id).get();
    };

    this.createDivision = function(division) {
        return Restangular.service('division').post(division)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Division created' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Division create failed' } ], false);
                }
            );
    };

    this.updateDivision = function(division) {
        return division.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Division updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Division update failed' } ], false);
                }
            );
    };

    this.deleteDivision = function(id) {
        return Restangular.one("division", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Division deleted' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Division delete failed' } ], false);
                }
            );
    };

};

DivisionService.$inject = ['Restangular', 'messageModel'];
module.exports = DivisionService;