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
import { ElementRef } from "@angular/core";
import { BehaviorSubject } from "rxjs";

import type { DataSet } from "src/app/models";

import { LoggingService } from "../logging.service";

import { LinechartDirective } from "./linechart.directive";

describe("LinechartDirective", () => {
	let directive: LinechartDirective;
	let dataSets: BehaviorSubject<Array<DataSet | null> | null>;

	beforeEach(()=>{
		dataSets =  new BehaviorSubject<Array<DataSet | null>|null>(null);
		directive = new LinechartDirective(new ElementRef(document.createElement("canvas")), new LoggingService());
		directive.chartDataSets = dataSets;
		directive.ngAfterViewInit();
	});

	it("should create an instance", () => {
		expect(directive).toBeTruthy();
	});

	it("sets a default type when not given", () => {
		expect(directive.chartType).toBe("linear");
	});

	it("loads new data", () => {
		// TODO: check more than just that these don't throw errors
		dataSets.next([{data: [{x: 1, y: 2}, {x: 2, y: 4}], label: "label"}]);
		dataSets.next([null, {data: [{x: 1, y: 2}], label: "label"}]);
		dataSets.next([{data: [{x: 1, y: 2}, {x: 2, y: 4}], label: "label"}, null]);
		dataSets.next(null);
	});

	it("handles errors in the data stream", () => {
		dataSets.error(new Error("some kind of problem with the server"));
	});
});
