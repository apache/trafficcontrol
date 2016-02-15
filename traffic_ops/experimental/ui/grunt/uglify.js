module.exports = {
    dist: {
        files: {
            '<%= globalConfig.resourcesdir %>/assets/js/app.js': [
                '<%= globalConfig.tmpdir %>/app.js'
            ]
        }
    }
};