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
            function() {
                messageModel.setMessages([ { level: 'error', text: 'Role update failed' } ], false);
            }
        );
    };

    this.deleteRole = function(id) {
        return Restangular.one("roles", id).remove()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Role deleted' } ], true);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'Role delete failed' } ], false);
            }
        );
    };

};

RoleService.$inject = ['Restangular', 'messageModel'];
module.exports = RoleService;