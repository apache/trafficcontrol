var CDNService = function(Restangular, messageModel) {

    this.getCDNs = function() {
        return Restangular.all('cdn').getList();
    };

    this.getCDN = function(id) {
        return Restangular.one("cdn", id).get();
    };

    this.updateCDN = function(cdn) {
        return cdn.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'CDN updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'CDN update failed' } ], false);
                }
            );
    };

    this.deleteCDN = function(id) {
        return Restangular.one("cdn", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'CDN deleted' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'CDN delete failed' } ], false);
                }
            );
    };

};

CDNService.$inject = ['Restangular', 'messageModel'];
module.exports = CDNService;