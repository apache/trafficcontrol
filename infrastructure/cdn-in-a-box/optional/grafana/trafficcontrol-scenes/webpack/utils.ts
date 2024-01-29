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

/* eslint-disable */

import fs from "fs";
import os from "os";
import path from "path";
import process from "process";

import { glob } from "glob";

import { SOURCE_DIR } from "./constants";

export function isWSL() {
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

export function getPackageJson() {
	return require(path.resolve(process.cwd(), "package.json"));
}

export function getPluginJson() {
	return require(path.resolve(process.cwd(), `${SOURCE_DIR}/plugin.json`));
}

export function hasReadme() {
	return fs.existsSync(path.resolve(process.cwd(), SOURCE_DIR, "README.md"));
}

// Support bundling nested plugins by finding all plugin.json files in src directory
// then checking for a sibling module.[jt]sx? file.
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

		result[entryName] = module;
		return result;
	}, result), {});
}
