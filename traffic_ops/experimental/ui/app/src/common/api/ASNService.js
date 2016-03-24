var ASNService = function(Restangular, messageModel) {

    this.getASNs = function() {
        return Restangular.all('asn').getList();
    };

    this.getASN = function(id) {
        return Restangular.one("asn", id).get();
    };

    this.createASN = function(asn) {
        return Restangular.service('asn').post(asn)
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'ASN created' } ], true);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'ASN create failed' } ], false);
            }
        );
    };

    this.updateASN = function(asn) {
        return asn.put()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'ASN updated' } ], false);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'ASN update failed' } ], false);
            }
        );
    };

    this.deleteASN = function(id) {
        return Restangular.one("asn", id).remove()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'ASN deleted' } ], true);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'ASN delete failed' } ], false);
            }
        );
    };

};

ASNService.$inject = ['Restangular', 'messageModel'];
module.exports = ASNService;