var ASNService = function(Restangular, locationUtils, messageModel) {

    this.getASNs = function() {
        return Restangular.all('asns').getList();
    };

    this.getASN = function(id) {
        return Restangular.one("asns", id).get();
    };

    this.createASN = function(asn) {
        return Restangular.service('asns').post(asn)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'ASN created' } ], true);
                    locationUtils.navigateToPath('/admin/asns');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateASN = function(asn) {
        return asn.put()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'ASN updated' } ], false);
            },
            function(fault) {
                messageModel.setMessages(fault.data.alerts, false);
            }
        );
    };

    this.deleteASN = function(id) {
        return Restangular.one("asns", id).remove()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'ASN deleted' } ], true);
            },
            function(fault) {
                messageModel.setMessages(fault.data.alerts, true);
            }
        );
    };

};

ASNService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = ASNService;