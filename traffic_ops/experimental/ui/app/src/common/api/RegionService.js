var RegionService = function(Restangular, messageModel) {

    this.getRegions = function() {
        return Restangular.all('regions').getList();
    };

    this.getRegion = function(id) {
        return Restangular.one("regions", id).get();
    };

    this.createRegion = function(region) {
        return Restangular.service('regions').post(region)
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Region created' } ], true);
            },
            function(fault) {
                messageModel.setMessages(fault.data.alerts, false);
            }
        );
    };

    this.updateRegion = function(region) {
        return region.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Region updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteRegion = function(id) {
        return Restangular.one("regions", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Region deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

};

RegionService.$inject = ['Restangular', 'messageModel'];
module.exports = RegionService;