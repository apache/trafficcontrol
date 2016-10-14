var DivisionService = function(Restangular, locationUtils, messageModel) {

    this.getDivisions = function() {
        return Restangular.all('divisions').getList();
    };

    this.getDivision = function(id) {
        return Restangular.one("divisions", id).get();
    };

    this.createDivision = function(division) {
        return Restangular.service('divisions').post(division)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Division created' } ], true);
                    locationUtils.navigateToPath('/admin/divisions');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateDivision = function(division) {
        return division.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Division updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteDivision = function(id) {
        return Restangular.one("divisions", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Division deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

};

DivisionService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = DivisionService;