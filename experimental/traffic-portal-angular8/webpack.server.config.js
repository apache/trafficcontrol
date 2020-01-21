/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

// Work around for https://github.com/angular/angular-cli/issues/7200

const path = require('path');
const webpack = require('webpack');

module.exports = {
	mode: 'none',
	entry: {
		// This is our Express server for Dynamic universal
		server: './server.ts'
	},
	target: 'node',
	resolve: { extensions: ['.ts', '.js'] },
	optimization: {
		minimize: false
	},
	output: {
		// Puts the output at the root of the dist folder
		path: path.join(__dirname, 'dist'),
		filename: '[name].js'
	},
	module: {
		rules: [
			{ test: /\.ts$/, loader: 'ts-loader' },
			{
			// Mark files inside `@angular/core` as using SystemJS style dynamic imports.
			// Removing this will cause deprecation warnings to appear.
			test: /(\\|\/)@angular(\\|\/)core(\\|\/).+\.js$/,
			parser: { system: true },
			},
		]
	},
	plugins: [
		new webpack.ContextReplacementPlugin(
			// fixes WARNING Critical dependency: the request of a dependency is an expression
			/(.+)?angular(\\|\/)core(.+)?/,
			path.join(__dirname, 'src'), // location of your src
			{} // a map of your routes
		),
		new webpack.ContextReplacementPlugin(
			// fixes WARNING Critical dependency: the request of a dependency is an expression
			/(.+)?express(\\|\/)(.+)?/,
			path.join(__dirname, 'src'),
			{}
		)
	]
};
