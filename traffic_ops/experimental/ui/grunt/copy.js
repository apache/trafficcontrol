module.exports = {
    dev: {
        files: [
            {
                expand: true,
                dot: true,
                cwd: '<%= globalConfig.srcdir %>',
                dest: '<%= globalConfig.resourcesdir %>',
                src: [
                    'assets/css/**/*',
                    'assets/fonts/**/*',
                    'assets/images/**/*',
                    'assets/js/**/*'
                ]
            },
            {
                expand: true,
                dot: true,
                cwd: '<%= globalConfig.srcdir %>',
                dest: '<%= globalConfig.distdir %>/public',
                src: [
                    '*.html',
                    'trafficOps_release.json'
                ]
            }
        ]
    },
    dist: {
        files: [
            {
                expand: true,
                dot: true,
                cwd: '<%= globalConfig.srcdir %>',
                dest: '<%= globalConfig.resourcesdir %>',
                src: [
                    'assets/css/**/*',
                    'assets/fonts/**/*',
                    'assets/images/**/*',
                    'assets/js/**/*'
                ]
            },
            {
                expand: true,
                dot: true,
                cwd: '<%= globalConfig.srcserverdir %>',
                dest: '<%= globalConfig.distdir %>/server',
                src: [
                    'server.js'
                ]
            },
            {
                expand: true,
                dot: true,
                cwd: '<%= globalConfig.srcdir %>',
                dest: '<%= globalConfig.distdir %>',
                src: [
                    'package.json'
                ]
            },
            {
                expand: true,
                dot: true,
                cwd: '<%= globalConfig.srcdir %>',
                dest: '<%= globalConfig.distdir %>/public',
                src: [
                    '*.html',
                    'trafficOps_release.json'
                ]
            }
        ]
    }
};
