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
                    'assets/images/**/*',
                    'assets/js/**/*',
                    'assets/other/**/*'
                ]
            },
            {
                expand: true,
                dot: true,
                cwd: '<%= globalConfig.srcdir %>',
                dest: '<%= globalConfig.distdir %>/public',
                src: [
                    '*.html',
                    'traffic_portal_release.json',
                    'traffic_portal_properties.json'
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
                    'assets/images/**/*',
                    'assets/js/**/*',
                    'assets/other/**/*'
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
                    'traffic_portal_properties.json'
                ]
            }
        ]
    }
};
