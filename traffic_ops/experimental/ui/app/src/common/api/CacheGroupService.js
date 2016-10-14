var CacheGroupService = function(Restangular, locationUtils, messageModel) {

    this.getCacheGroups = function() {
        return Restangular.all('cachegroups').getList();
    };

    this.getCacheGroup = function(id) {
        return Restangular.one("cachegroups", id).get();
    };

    this.createCacheGroup = function(cacheGroup) {
        return Restangular.service('cachegroups').post(cacheGroup)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'CacheGroup created' } ], true);
                    locationUtils.navigateToPath('/configure/cache-groups');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateCacheGroup = function(cacheGroup) {
        return cacheGroup.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Cache group updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteCacheGroup = function(id) {
        return Restangular.one("cachegroups", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Cache group deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

};

CacheGroupService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = CacheGroupService;