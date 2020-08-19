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

const dockerComposeArgs = [
	"-f",
	`${process.env.GITHUB_WORKSPACE}/infrastructure/docker/build/docker-compose.yml`,
	"up"
];

const spawnArgs = {stdio: "inherit", stderr: "inherit"};

const proc = child_process.spawnSync(
	"docker-compose",
	dockerComposeArgs,
	spawnArgs
);

if (proc.status !== 0) {
	console.error("Building the RPMs failed");
} else {
	console.log("Finished building RPMS");
}
process.exit(proc.status);
