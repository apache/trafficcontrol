var FormUtils = function() {

    this.hasError = function(input) {
        return !input.$focused && input.$invalid;
    };

    this.hasPropertyError = function(input, property) {
        return !input.$focused && input.$error[property];
    };

};

FormUtils.$inject = [];
module.exports = FormUtils;
