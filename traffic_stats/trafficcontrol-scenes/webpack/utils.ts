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

import fs from "fs";
import os from "os";
import path from "path";
import process from "process";

import { glob } from "glob";

import { SOURCE_DIR } from "./constants";

/**
 * Check if the current environment is Windows Subsystem for Linux (WSL).
 *
 * @returns true if the current environment is WSL, false otherwise
 */
export function isWSL(): boolean {
	if (process.platform !== "linux") {
		return false;
	}

	if (os.release().toLowerCase().includes("microsoft")) {
		return true;
	}

	try {
		return fs.readFileSync("/proc/version", "utf8").toLowerCase().includes("microsoft");
	} catch {
		return false;
	}
}

/**
 * getPackageJson
 *
 * @returns package.json content
 */
// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types,@typescript-eslint/explicit-function-return-type
export function getPackageJson() {
	// eslint-disable-next-line @typescript-eslint/no-require-imports
	return require(path.resolve(process.cwd(), "package.json"));
}

/**
 * Check if a README.md file exists in the SOURCE_DIR directory.
 *
 * @returns true if README.md exists, false otherwise
 */
// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types,@typescript-eslint/explicit-function-return-type
export function getPluginJson() {
	// eslint-disable-next-line @typescript-eslint/no-require-imports
	return require(path.resolve(process.cwd(), `${SOURCE_DIR}/plugin.json`));
}

// Support bundling nested plugins by finding all plugin.json files in src directory
// then checking for a sibling module.[jt]sx? file.
/**
 * Asynchronously retrieves entries for plugins.
 *
 * @returns a Promise that resolves to a record of plugin entries
 */
export async function getEntries(): Promise<Record<string, string>> {
	const pluginsJson = await glob("**/src/**/plugin.json", { absolute: true });

	const plugins = await Promise.all(
		pluginsJson.map(async (pluginJson) => {
			const folder = path.dirname(pluginJson);
			return glob(`${folder}/module.{ts,tsx,js,jsx}`, { absolute: true });
		})
	);

	return plugins.reduce((result, modules) => modules.reduce((result, module) => {
		const pluginPath = path.dirname(module);
		const pluginName = path.relative(process.cwd(), pluginPath).replace(/src\/?/i, "");
		const entryName = pluginName === "" ? "module" : `${pluginName}/module`;

		// @ts-ignore
		result[entryName] = module;
		return result;
	}, result), {});
}
