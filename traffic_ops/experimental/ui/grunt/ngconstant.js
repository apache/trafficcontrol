module.exports = {
    options: {
        space: '  ',
        wrap: '"use strict";\n\n {%= __ngModule %}',
        name: 'config',
        dest: '<%= globalConfig.srcdir %>/scripts/config.js'
    },
    dev: {
        constants: {
            ENV: {
                apiEndpoint: {
                    '1.1': '/api/1.1/',
                    '1.2': '/api/1.2/',
                    '2.0': '/api/2.0/'
                }
            }
        }
    },
    prod: {
        constants: {
            ENV: {
                apiEndpoint: {
                    '1.1': '/api/1.1/',
                    '1.2': '/api/1.2/',
                    '2.0': '/api/2.0/'
                }
            }
        }
    }
};

