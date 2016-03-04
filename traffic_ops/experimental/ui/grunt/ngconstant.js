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
                    "base_url": "/api/" // this is only used to identify api calls. See conf/config.js to define api location.
                }
            }
        }
    },
    prod: {
        constants: {
            ENV: {
                api: {
                    "base_url": "/api/" // this is only used to identify api calls. See conf/config.js to define api location.
                }
            }
        }
    }
};

