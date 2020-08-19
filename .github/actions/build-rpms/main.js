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

const child_process = require("child_process");
const spawnArgs = {
	stdio: "inherit",
	stderr: "inherit",
};

let atcComponent = process.env.ATC_COMPONENT;
const dockerComposeArgs = ["-f", `${process.env.GITHUB_WORKSPACE}/infrastructure/docker/build/docker-compose.yml`, "run", "--rm"];
if (typeof atcComponent !== "string" || atcComponent.length === 0) {
	console.error("Missing environment variable ATC_COMPONENT");
	process.exit(1);
}
const nonRpmComponents = ["source", "weasel", "docs"];
if (nonRpmComponents.indexOf(atcComponent) === -1) {
	atcComponent += "_build";
}
dockerComposeArgs.push(atcComponent);
const proc = child_process.spawnSync(
	"docker-compose",
	dockerComposeArgs,
	spawnArgs
);
process.exit(proc.status);
