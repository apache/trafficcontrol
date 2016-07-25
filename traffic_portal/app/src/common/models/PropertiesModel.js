var PropertiesModel = function() {

    this.properties = {};
    this.loaded = false;

    this.setProperties = function(properties) {
        this.properties = properties;
        this.loaded = true;
    };

};

PropertiesModel.$inject = [];
module.exports = PropertiesModel;