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

import { HttpClientModule } from "@angular/common/http";
import { TestBed, waitForAsync } from "@angular/core/testing";
import { RouterTestingModule } from "@angular/router/testing";
import { MatSnackBarModule } from "@angular/material/snack-bar";

import { AppComponent } from "./app.component";
import { AlertComponent } from "./components/alert/alert.component";

describe("AppComponent", () => {
	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [
				AlertComponent,
				AppComponent
			],
			imports: [
				HttpClientModule,
				RouterTestingModule,
				MatSnackBarModule
			],
		}).compileComponents();
	}));

	it("should create the app", () => {
		const fixture = TestBed.createComponent(AppComponent);
		const app = fixture.componentInstance;
		expect(app).toBeTruthy();
	});

	it("should have as title 'Traffic Portal'", () => {
		const fixture = TestBed.createComponent(AppComponent);
		const app = fixture.componentInstance;
		expect(app.title).toEqual("Traffic Portal");
	});
});
