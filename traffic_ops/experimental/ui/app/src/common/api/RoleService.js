var RoleService = function(Restangular, messageModel) {

    this.getRoles = function() {
        return Restangular.all('roles').getList();
    };

    this.getRole = function(id) {
        return Restangular.one("roles", id).get();
    };

    this.updateRole = function(role) {
        return role.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Role updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
        );
    };

    this.deleteRole = function(id) {
        return Restangular.one("roles", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Role deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
        );
    };

};

RoleService.$inject = ['Restangular', 'messageModel'];
module.exports = RoleService;