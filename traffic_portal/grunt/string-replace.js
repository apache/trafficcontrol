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
		files: {
			expand: true,
			cwd: '<%= globalConfig.distdir %>/public',
			dest: '<%= globalConfig.distdir %>/public',
			src: 'index.html',
		},
		options: {
			replacements: [
				{
					pattern: 'theme.css',
					replacement: 'theme.css?built=' + Date.now()
				},
				{
					pattern: 'loading.css',
					replacement: 'loading.css?built=' + Date.now()
				},
				{
					pattern: 'main.css',
					replacement: 'main.css?built=' + Date.now()
				},
				{
					pattern: 'custom.css',
					replacement: 'custom.css?built=' + Date.now()
				},
				{
					pattern: 'shared-libs.js',
					replacement: 'shared-libs.js?built=' + Date.now()
				},
				{
					pattern: 'app.js',
					replacement: 'app.js?built=' + Date.now()
				},
				{
					pattern: 'config.js',
					replacement: 'config.js?built=' + Date.now()
				}
			]
		}
};
