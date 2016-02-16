module.exports = function() {
    var globalConfig = {
        app: 'app',
        resourcesdir: 'app/dist/public/resources',
        distdir: 'app/dist',
        srcserverdir: './server',
        srcdir: 'app/src',
        tmpdir: '.tmp',
        srcfiles: {
            js: ['./app/src/**/*.js'],
            tpl: ['./app/src/**/*.tpl.html']
        }
    };

    return globalConfig;
}
