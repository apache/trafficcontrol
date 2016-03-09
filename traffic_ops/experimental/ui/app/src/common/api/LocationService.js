var LocationService = function(Restangular, messageModel) {

    this.getLocations = function() {
        return Restangular.all('phys_location').getList();
    };

    this.getLocation = function(id) {
        return Restangular.one("phys_location", id).get();
    };

    this.updateLocation = function(location) {
        return location.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Location updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Location update failed' } ], false);
                }
            );
    };

};

LocationService.$inject = ['Restangular', 'messageModel'];
module.exports = LocationService;