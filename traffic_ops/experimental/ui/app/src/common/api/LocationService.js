var LocationService = function(Restangular, locationUtils, messageModel) {

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
                    locationUtils.navigateToPath('/admin/locations');

                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateLocation = function(location) {
        return location.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Location updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteLocation = function(id) {
        return Restangular.one("phys_locations", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Location deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

};

LocationService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = LocationService;