var RegionService = function(Restangular, messageModel) {

    this.getRegions = function() {
        return Restangular.all('region').getList();
    };

    this.getRegion = function(id) {
        return Restangular.one("region", id).get();
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

};

RegionService.$inject = ['Restangular', 'messageModel'];
module.exports = RegionService;