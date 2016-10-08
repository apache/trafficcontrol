var LocationService = function(Restangular, messageModel) {

    this.getLocations = function() {
        return Restangular.all('phys_locations').getList();
    };

    this.getLocation = function(id) {
        return Restangular.one("phys_locations", id).get();
    };

    this.createLocation = function(location) {
        return Restangular.service('phys_locations').post(location)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Location created' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Location create failed' } ], false);
                }
            );
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

    this.deleteLocation = function(id) {
        return Restangular.one("phys_locations", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Location deleted' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Location delete failed' } ], false);
                }
            );
    };

};

LocationService.$inject = ['Restangular', 'messageModel'];
module.exports = LocationService;