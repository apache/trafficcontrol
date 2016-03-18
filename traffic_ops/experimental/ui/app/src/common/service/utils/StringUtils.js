var FormUtils = function() {

    this.labelize = function(string) {
        return string.replace(/([A-Z])/g, ' $1').replace(/^./, function(str){ return str.toUpperCase(); });
    };

};

FormUtils.$inject = [];
module.exports = FormUtils;
