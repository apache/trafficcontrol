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
import { Component } from "@angular/core";

/** LoadingComponent is the controller for a spinning "loading" icon. */
@Component({
	selector: "tp-loading",
	styleUrls: ["./loading.component.scss"],
	templateUrl: "./loading.component.html"
})
// need a class to bind to a template - even if there's no data or logic.
// eslint-disable-next-line @typescript-eslint/no-extraneous-class
export class LoadingComponent {
}
