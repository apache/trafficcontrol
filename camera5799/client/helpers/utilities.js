/**
 * @author Ramiro Arenivar
 * For CSCI 5799
 */

Utilities = {

    getUsername: function () {
        if (Session.get('login_response')) {
            if (Session.get('login_response')['username']) {
                return Session.get('login_response')['username'];
            }
        }
        return null;
    },

    getUserToken: function() {
        if (Session.get('login_response')) {
            if (Session.get('login_response')['token']) {
                return Session.get('login_response')['token'];
            }
        }
        return null;
    },

    clearForm: function(formId) {
        $('#' + formId).find(':input').each(function() {
            switch (this.type) {
                case 'password':
                case 'select-multiple':
                case 'select-one':
                case 'text':
                case 'textarea':
                    $(this).val('');
                    break;
                case 'checkbox':
                case 'radio':
                    this.checked = false;
            }
        });
    }
}