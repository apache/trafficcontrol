
module.exports = {
    options: {
        base: './app/src'
    },
    'dist': {
        src: ['<%= globalConfig.srcfiles.tpl %>'],
        dest: '<%= globalConfig.tmpdir %>/app-templates.js',
        module: 'app.templates'
    }
};