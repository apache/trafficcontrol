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
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSelectModule } from "@angular/material/select";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { RouterTestingModule } from "@angular/router/testing";
import { ReplaySubject } from "rxjs";

import { APITestingModule } from "src/app/api/testing";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

import { CDNTableComponent } from "./cdn-table.component";

describe("CDNTableComponent", () => {
	let component: CDNTableComponent;
	let fixture: ComponentFixture<CDNTableComponent>;

	const navService = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [CDNTableComponent],
			imports: [
				APITestingModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule,
				MatDialogModule,
				NoopAnimationsModule,
				MatSelectModule,
			],
			providers: [
				{provide: NavigationService, useValue: navService},
			],
		})
			.compileComponents();

		fixture = TestBed.createComponent(CDNTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
