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

'use strict';

module.exports = function (grunt) {

    var os = require("os");
    var globalConfig = require('./grunt/globalConfig');

    // load time grunt - helps with optimizing build times
    require('time-grunt')(grunt);

    // load grunt task configurations
    require('load-grunt-config')(grunt);

    // load string replacement functionality. used to add query params to index.html to bust cache on upgrade
    require('grunt-string-replace')(grunt);

    // default task - when you type 'grunt' it really runs as 'grunt dev'
    grunt.registerTask('default', ['dev']);

    // dev task - when you type 'grunt dev' <-- builds unminified app and puts it in in app/dist folder and starts express server which reads server.js
    grunt.registerTask('dev', [
        'build-dev',
        'express:dev',
        'watch'
    ]);

    // dist task - when you type 'grunt dist' <-- builds minified app for distribution and generates node dependencies all wrapped up nicely in app/dist folder
    grunt.registerTask('dist', [
        'build'
    ]);

    // build tasks
    grunt.registerTask('build', [
        'clean',
        'copy:dist',
        'string-replace',
        'build-css',
        'build-js',
        'build-shared-libs'
    ]);

    grunt.registerTask('build-dev', [
        'clean',
        'copy:dev',
        'string-replace',
        'build-css-dev',
        'build-js-dev',
        'build-shared-libs-dev'
    ]);

    // css
    grunt.registerTask('build-css', [
        'dart-sass:prod'
    ]);

    grunt.registerTask('build-css-dev', [
        'dart-sass:dev'
    ]);

    // js (custom)
    grunt.registerTask('build-js', [
        'html2js',
        'browserify:app-prod',
        'browserify:app-config'
    ]);

    grunt.registerTask('build-js-dev', [
        'html2js',
        'browserify:app-dev',
        'browserify:app-config'
    ]);

    // js (libraries)
    grunt.registerTask('build-shared-libs', [
        'browserify:shared-libs-prod'
    ]);

    grunt.registerTask('build-shared-libs-dev', [
        'browserify:shared-libs-dev'
    ]);

};
