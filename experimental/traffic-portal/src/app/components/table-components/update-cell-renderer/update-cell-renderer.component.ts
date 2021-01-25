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
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { faCheck, faClock } from "@fortawesome/free-solid-svg-icons";
import { ICellRendererAngularComp } from "ag-grid-angular";
import { ICellRendererParams } from "ag-grid-community";

@Component({
	selector: "tp-update-cell-renderer",
	styleUrls: ["./update-cell-renderer.component.scss"],
	templateUrl: "./update-cell-renderer.component.html",
})
export class UpdateCellRendererComponent implements ICellRendererAngularComp {

	public value = false;

	public get icon(): IconDefinition {
		return this.value ? faClock : faCheck;
	}

	public refresh(params: ICellRendererParams): true {
		this.value = params.value;
		return true;
	}

	public agInit(params: ICellRendererParams): void {
		this.value = params.value;
	}
}
