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
                    'traffic_portal_release.json',
                    'traffic_portal_properties.json'
                ]
            },
            {
                expand: true,
                dot: true,
                cwd: "<%= globalConfig.importdir %>",
                dest: "<%= globalConfig.resourcesdir %>/assets/js/",
                src: ["ag-grid-community/dist/ag-grid-community.min.js"]
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
                    'traffic_portal_release.json',
                    'traffic_portal_properties.json'
                ]
            },
            {
                expand: true,
                dot: true,
                cwd: "<%= globalConfig.importdir %>",
                dest: "<%= globalConfig.resourcesdir %>/assets/js/",
                src: ["ag-grid-community/dist/ag-grid-community.min.js"]
            }
        ]
    }
};
