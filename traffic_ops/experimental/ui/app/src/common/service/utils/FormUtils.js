var FormUtils = function() {

    this.hasError = function(input) {
        return !input.$focused && input.$dirty && input.$invalid;
    };

    this.hasPropertyError = function(input, property) {
        return !input.$focused && input.$dirty && input.$error[property];
    };

};

FormUtils.$inject = [];
module.exports = FormUtils;
