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
            function() {
                messageModel.setMessages([ { level: 'error', text: 'Region create failed' } ], false);
            }
        );
    };

    this.updateRegion = function(region) {
        return region.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Region updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Region update failed' } ], false);
                }
            );
    };

    this.deleteRegion = function(id) {
        return Restangular.one("regions", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Region deleted' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Region delete failed' } ], false);
                }
            );
    };

};

RegionService.$inject = ['Restangular', 'messageModel'];
module.exports = RegionService;