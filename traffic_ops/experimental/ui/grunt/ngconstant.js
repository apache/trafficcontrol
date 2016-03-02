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
                    "login": "/api/2.0/login",
                    'get_current_user': '/api/2.0/tm_user/current',
                    'update_current_user': '/api/2.0/tm_user/current',
                    'get_users': '/api/2.0/tm_user'
                }
            }
        }
    },
    prod: {
        constants: {
            ENV: {
                apiEndpoint: {
                    "login": "/api/2.0/login",
                    'get_current_user': '/api/2.0/tm_user/current',
                    'update_current_user': '/api/2.0/tm_user/current',
                    'get_users': '/api/2.0/tm_user'
                }
            }
        }
    }
};

