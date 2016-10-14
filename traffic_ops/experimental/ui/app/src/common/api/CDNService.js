var CDNService = function(Restangular, locationUtils, messageModel) {

    this.getCDNs = function() {
        return Restangular.all('cdns').getList();
    };

    this.getCDN = function(id) {
        return Restangular.one("cdns", id).get();
    };

    this.createCDN = function(cdn) {
        return Restangular.service('cdns').post(cdn)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'CDN created' } ], true);
                    locationUtils.navigateToPath('/admin/cdns');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateCDN = function(cdn) {
        return cdn.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'CDN updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteCDN = function(id) {
        return Restangular.one("cdns", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'CDN deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

};

CDNService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = CDNService;