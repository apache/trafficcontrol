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

import path from "path";

import CopyWebpackPlugin from "copy-webpack-plugin";
import ESLintPlugin from "eslint-webpack-plugin";
import ForkTsCheckerWebpackPlugin from "fork-ts-checker-webpack-plugin";
import ReplaceInFileWebpackPlugin from "replace-in-file-webpack-plugin";
import { Configuration } from "webpack";
import LiveReloadPlugin from "webpack-livereload-plugin";

import { DIST_DIR, SOURCE_DIR } from "./constants";
import { getEntries, getPackageJson, getPluginJson, isWSL } from "./utils";

const pluginJson = getPluginJson();

const config = async (env): Promise<Configuration> => {
	const baseConfig: Configuration = {
		cache: {
			buildDependencies: {
				config: [__filename],
			},
			type: "filesystem",
		},

		context: path.join(process.cwd(), SOURCE_DIR),

		devtool: env.production ? "source-map" : "eval-source-map",

		entry: await getEntries(),

		externals: [
			"lodash",
			"jquery",
			"moment",
			"slate",
			"emotion",
			"@emotion/react",
			"@emotion/css",
			"prismjs",
			"slate-plain-serializer",
			"@grafana/slate-react",
			"react",
			"react-dom",
			"react-redux",
			"redux",
			"rxjs",
			"react-router",
			"react-router-dom",
			"d3",
			"angular",
			"@grafana/ui",
			"@grafana/runtime",
			"@grafana/data",

			// Mark legacy SDK imports as external if their name starts with the "grafana/" prefix
			({request}, callback) => {
				const prefix = "grafana/";
				const hasPrefix = (request) => request.indexOf(prefix) === 0;
				const stripPrefix = (request) => request.substr(prefix.length);

				if (hasPrefix(request)) {
					return callback(undefined, stripPrefix(request));
				}

				callback();
			},
		],

		mode: env.production ? "production" : "development",

		module: {
			rules: [
				{
					exclude: /(node_modules)/,
					test: /\.[tj]sx?$/,
					use: {
						loader: "swc-loader",
						options: {
							jsc: {
								baseUrl: path.resolve(__dirname, "src"),
								loose: false,
								parser: {
									decorators: false,
									dynamicImport: true,
									syntax: "typescript",
									tsx: true,
								},
								target: "es2015",
							},
						},
					},
				},
				{
					test: /\.css$/,
					use: ["style-loader", "css-loader"],
				},
				{
					test: /\.s[ac]ss$/,
					use: ["style-loader", "css-loader", "sass-loader"],
				},
				{
					generator: {
						// Keep publicPath relative for host.com/grafana/ deployments
						filename: Boolean(env.production)
							? "[hash][ext]"
							: "[file]",
						outputPath: "img/",
						publicPath: `public/plugins/${pluginJson.id}/img/`,
					},
					test: /\.(png|jpe?g|gif|svg)$/,
					type: "asset/resource",
				},
				{
					generator: {
						// Keep publicPath relative for host.com/grafana/ deployments
						filename: Boolean(env.production)
							? "[hash][ext]"
							: "[name][ext]",
						outputPath: "fonts/",
						publicPath: `public/plugins/${pluginJson.id}/fonts/`,
					},
					test: /\.(woff|woff2|eot|ttf|otf)(\?v=\d+\.\d+\.\d+)?$/,
					type: "asset/resource",
				},
			],
		},

		output: {
			clean: {
				keep: new RegExp(
					"(.*?_(amd64|arm(64)?)(.exe)?|go_plugin_build_manifest)"
				),
			},
			filename: "[name].js",
			library: {
				type: "amd",
			},
			path: path.resolve(process.cwd(), DIST_DIR),
			publicPath: `public/plugins/${pluginJson.id}/`,
			uniqueName: pluginJson.id,
		},

		plugins: [
			new CopyWebpackPlugin({
				patterns: [
					// If src/README.md exists use it; otherwise the root README
					// To `compiler.options.output`
					{from: "plugin.json", to: "."},
					{from: "**/*.json", to: "."}, // TODO<Add an error for checking the basic structure of the repo>
					{from: "**/*.svg", noErrorOnMissing: true, to: "."}, // Optional
					{from: "**/*.png", noErrorOnMissing: true, to: "."}, // Optional
					{from: "**/*.html", noErrorOnMissing: true, to: "."}, // Optional
					{from: "img/**/*", noErrorOnMissing: true, to: "."}, // Optional
					{from: "libs/**/*", noErrorOnMissing: true, to: "."}, // Optional
					{from: "static/**/*", noErrorOnMissing: true, to: "."}, // Optional
					{
						from: "**/query_help.md",
						noErrorOnMissing: true,
						to: ".",
					}, // Optional
				],
			}),
			// Replace certain template-variables in the README and plugin.json
			new ReplaceInFileWebpackPlugin([
				{
					dir: DIST_DIR,
					files: ["plugin.json"],
					rules: [
						{
							replace: getPackageJson().version,
							search: /\%VERSION\%/g,
						},
						{
							replace: new Date().toISOString().substring(0, 10),
							search: /\%TODAY\%/g,
						},
						{
							replace: pluginJson.id,
							search: /\%PLUGIN_ID\%/g,
						},
					],
				},
			]),
			new ForkTsCheckerWebpackPlugin({
				async: Boolean(env.development),
				issue: {
					include: [{file: "**/*.{ts,tsx}"}],
				},
				typescript: {
					configFile: path.join(process.cwd(), "tsconfig.json"),
				},
			}),
			new ESLintPlugin({
				extensions: [".ts", ".tsx"],
				lintDirtyModulesOnly: Boolean(env.development), // don't lint on start, only lint changed files
			}),
			...(env.development ? [new LiveReloadPlugin()] : []),
		],

		resolve: {
			extensions: [".js", ".jsx", ".ts", ".tsx"],
			// handle resolving "rootDir" paths
			modules: [path.resolve(process.cwd(), "src"), "node_modules"],
			unsafeCache: true,
		},
	};

	if (isWSL()) {
		baseConfig.watchOptions = {
			ignored: /node_modules/,
			poll: 3000,
		};
	}

	return baseConfig;
};

export default config;
