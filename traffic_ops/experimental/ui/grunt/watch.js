
module.exports = {
    options: {
        livereload: true
    },
    js: {
        files: ['<%= globalConfig.srcdir %>/**/*.js'],
        tasks: [ 'build-dev']
    },
    compass: {
        files: ['<%= globalConfig.srcdir %>/**/*.{scss,sass}'],
        tasks: ['compass:dev']
    },
    html: {
        files: ['app/**/*.tpl.html', 'app/**/index.html'],
        tasks: ['copy:dist', 'build-dev']
    }
};