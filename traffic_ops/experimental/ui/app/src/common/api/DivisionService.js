var DivisionService = function(Restangular, messageModel) {

    this.getDivisions = function() {
        return Restangular.all('division').getList();
    };

    this.getDivision = function(id) {
        return Restangular.one("division", id).get();
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

};

DivisionService.$inject = ['Restangular', 'messageModel'];
module.exports = DivisionService;