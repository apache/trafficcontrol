var CacheGroupService = function(Restangular, messageModel) {

    this.getCacheGroups = function() {
        return Restangular.all('cachegroup').getList();
    };

    this.getCacheGroup = function(id) {
        return Restangular.one("cachegroup", id).get();
    };

    this.createCacheGroup = function(cacheGroup) {
        return Restangular.service('cachegroup').post(cacheGroup)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'CacheGroup created' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'CacheGroup create failed' } ], false);
                }
            );
    };

    this.updateCacheGroup = function(cacheGroup) {
        return cacheGroup.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Cache group updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Cache group update failed' } ], false);
                }
            );
    };

    this.deleteCacheGroup = function(id) {
        return Restangular.one("cachegroup", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Cache group deleted' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Cache group delete failed' } ], false);
                }
            );
    };

};

CacheGroupService.$inject = ['Restangular', 'messageModel'];
module.exports = CacheGroupService;