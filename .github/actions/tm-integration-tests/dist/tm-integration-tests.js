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
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const child_process_1 = __importDefault(require("child_process"));
const path_1 = __importDefault(require("path"));
const spawnOptions = { stdio: "inherit" };
function runProcess(...commandArguments) {
    console.info(...commandArguments);
    const { status } = child_process_1.default.spawnSync(commandArguments[0], commandArguments.slice(1), spawnOptions);
    if (status === 0) {
        return;
    }
    console.error("Child process \"", ...commandArguments, "\" exited with status code", status, "!");
    process.exit(status !== null && status !== void 0 ? status : 1);
}
runProcess(path_1.default.join(__dirname, "../entrypoint.sh"));
