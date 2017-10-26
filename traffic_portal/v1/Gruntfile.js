/*


 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

 */

'use strict';

module.exports = function (grunt) {
    var os = require("os");
    var globalConfig = require('./grunt/globalConfig');

    // load time grunt - helps with optimizing build times
    require('time-grunt')(grunt);

    // load grunt task configurations
    require('load-grunt-config')(grunt);

    // default task - runs in dev mode
    grunt.registerTask('default', ['dev']);

    // dev task - when you type 'grunt dev' <-- builds unminified app and starts dev node.js server and watches for source file changes
    grunt.registerTask('dev', [
        'build-dev',
        'express:dev',
        'watch'
    ]);

    // dist task - when you type 'grunt dist' <-- builds minified app for distribution and generates node dependencies all wrapped up nicely in a /dist folder
    grunt.registerTask('dist', [
        'build',
        'install-dependencies'
    ]);

    // build tasks
    grunt.registerTask('build', [
        'ngconstant:prod',
        'clean',
        'copy:dist',
        'build-css',
        'build-js',
        'build-shared-libs'
    ]);

    grunt.registerTask('build-dev', [
        'ngconstant:dev',
        'clean',
        'copy:dev',
        'build-css-dev',
        'build-js-dev',
        'build-shared-libs-dev'
    ]);

    // css
    grunt.registerTask('build-css', [
        'compass:prod'
    ]);

    grunt.registerTask('build-css-dev', [
        'compass:dev'
    ]);

    // js (custom)
    grunt.registerTask('build-js', [
        'html2js',
        'browserify2:app-prod',
        'browserify2:app-config'
    ]);

    grunt.registerTask('build-js-dev', [
        'html2js',
        'browserify2:app-dev',
        'browserify2:app-config'
    ]);

    // js (libraries)
    grunt.registerTask('build-shared-libs', [
        'browserify2:shared-libs-prod'
    ]);

    grunt.registerTask('build-shared-libs-dev', [
        'browserify2:shared-libs-dev'
    ]);

};