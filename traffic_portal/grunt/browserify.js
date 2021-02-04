/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

module.exports = {
    'shared-libs-prod': {
        src: ['./<%= globalConfig.srcdir %>/scripts/shared-libs.js'],
        dest: './<%= globalConfig.resourcesdir %>/assets/js/shared-libs.js',
        options: {
            alias: {
                "angular": "./<%= globalConfig.app %>/bower_components/angular/angular.min.js",
                "angular-animate": './<%= globalConfig.app %>/bower_components/angular-animate/angular-animate.min.js',
                "contextMenu": './<%= globalConfig.app %>/bower_components/angular-bootstrap-contextmenu/contextMenu.js',
                "ui-bootstrap": './<%= globalConfig.app %>/bower_components/angular-bootstrap/ui-bootstrap.min.js',
                "ui-bootstrap-tpls": './<%= globalConfig.app %>/bower_components/angular-bootstrap/ui-bootstrap-tpls.min.js',
                "angular-jwt": './<%= globalConfig.app %>/bower_components/angular-jwt/dist/angular-jwt.min.js',
                "loading-bar": './<%= globalConfig.app %>/bower_components/angular-loading-bar/build/loading-bar.min.js',
                "angular-resource": './<%= globalConfig.app %>/bower_components/angular-resource/angular-resource.min.js',
                "angular-route": './<%= globalConfig.app %>/bower_components/angular-route/angular-route.min.js',
                "angular-sanitize": './<%= globalConfig.app %>/bower_components/angular-sanitize/angular-sanitize.min.js',
                "angular-ui-router": './<%= globalConfig.app %>/bower_components/angular-ui-router/release/angular-ui-router.min.js',
                "bootstrap": './<%= globalConfig.app %>/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap.min.js',
                "es5-shim": './<%= globalConfig.app %>/bower_components/es5-shim/es5-shim.min.js',
                "jquery": './<%= globalConfig.app %>/bower_components/jquery/jquery.min.js',
                "json3": './<%= globalConfig.app %>/bower_components/json3/lib/json3.min.js',
                'jquery-flot': './<%= globalConfig.importdir %>/flot/dist/es5/jquery.flot.js',
                'jquery-flot-pie': './<%= globalConfig.importdir %>/flot/source/jquery.flot.pie.js',
                'jquery-flot-stack': './<%= globalConfig.importdir %>/flot/source/jquery.flot.stack.js',
                'jquery-flot-time': './<%= globalConfig.importdir %>/flot/source/jquery.flot.time.js',
                'jquery-flot-tooltip': './<%= globalConfig.importdir %>/jquery.flot.tooltip/js/jquery.flot.tooltip.js',
                'jquery-flot-axislabels': './<%= globalConfig.importdir %>/flot/source/jquery.flot.axislabels.js',
            },
        },
    },
    'shared-libs-dev': {
        src: ['./<%= globalConfig.srcdir %>/scripts/shared-libs.js'],
        dest: './<%= globalConfig.resourcesdir %>/assets/js/shared-libs.js',
        options: {
            alias: {
                "angular": "./<%= globalConfig.importdir %>/angular/angular.min.js",
                "angular-animate": './<%= globalConfig.importdir %>/angular-animate/angular-animate.min.js',
                "contextMenu": './<%= globalConfig.importdir %>/angular-bootstrap-contextmenu/contextMenu.js',
                "ui-bootstrap": './<%= globalConfig.importdir %>/angular-ui-bootstrap/ui-bootstrap.min.js',
                "ui-bootstrap-tpls": './<%= globalConfig.importdir %>/angular-ui-bootstrap/ui-bootstrap-tpls.min.js',
                "angular-jwt": './<%= globalConfig.importdir %>/angular-jwt/dist/angular-jwt.min.js',
                "loading-bar": './<%= globalConfig.importdir %>/angular-loading-bar/build/loading-bar.min.js',
                "angular-resource": './<%= globalConfig.importdir %>/angular-resource/angular-resource.min.js',
                "angular-route": './<%= globalConfig.importdir %>/angular-route/angular-route.min.js',
                "angular-sanitize": './<%= globalConfig.importdir %>/angular-sanitize/angular-sanitize.min.js',
                "angular-ui-router": './<%= globalConfig.importdir %>/@uirouter/angularjs/release/angular-ui-router.min.js',
                "bootstrap": './<%= globalConfig.app %>/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap.min.js',
                "es5-shim": './<%= globalConfig.app %>/bower_components/es5-shim/es5-shim.min.js',
                "jquery": './<%= globalConfig.app %>/bower_components/jquery/jquery.min.js',
                "json3": './<%= globalConfig.app %>/bower_components/json3/lib/json3.min.js',
                'jquery-flot': './<%= globalConfig.importdir %>/flot/dist/es5/jquery.flot.js',
                'jquery-flot-pie': './<%= globalConfig.importdir %>/flot/source/jquery.flot.pie.js',
                'jquery-flot-stack': './<%= globalConfig.importdir %>/flot/source/jquery.flot.stack.js',
                'jquery-flot-time': './<%= globalConfig.importdir %>/flot/source/jquery.flot.time.js',
                'jquery-flot-tooltip': './<%= globalConfig.importdir %>/jquery.flot.tooltip/js/jquery.flot.tooltip.js',
                'jquery-flot-axislabels': './<%= globalConfig.importdir %>/flot/source/jquery.flot.axislabels.js',
            },
        },
    },
    'app-prod': {
        src: ['./<%= globalConfig.srcdir %>/app.js'],
        dest: './<%= globalConfig.resourcesdir %>/assets/js/app.js',
        browserifyOptions: {
            debug: false,
        },
        options: {
            alias: {
                'app-templates': './<%= globalConfig.tmpdir %>/app-templates.js'
            }
        }
    },
    'app-dev': {
        src: ['./<%= globalConfig.srcdir %>/app.js'],
        dest: './<%= globalConfig.resourcesdir %>/assets/js/app.js',
        browserifyOptions: {
            debug: true,
        },
        options: {
            alias: {
                'app-templates':'./<%= globalConfig.tmpdir %>/app-templates.js'
            }
        }
    },
    'app-config': {
        src: ['./<%= globalConfig.srcdir %>/scripts/config.js'],
        dest: './<%= globalConfig.resourcesdir %>/assets/js/config.js'
    }
};
