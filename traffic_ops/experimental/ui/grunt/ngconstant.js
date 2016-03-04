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
                api: {
                    "root": "/api/2.0/" // api base_url is defined in server.js
                }
            }
        }
    },
    prod: {
        constants: {
            ENV: {
                api: {
                    "root": "/api/2.0/" // api base_url is defined in server.js
                }
            }
        }
    }
};

