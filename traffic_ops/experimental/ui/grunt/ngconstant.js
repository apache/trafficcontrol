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
                    "base_url": "/api/2.0/"
                }
            }
        }
    },
    prod: {
        constants: {
            ENV: {
                api: {
                    "base_url": "/api/2.0/"
                }
            }
        }
    }
};

