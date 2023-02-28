/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import { HttpClientTestingModule } from "@angular/common/http/testing";
import { ComponentFixture, fakeAsync, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatCardModule } from "@angular/material/card";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatGridListModule } from "@angular/material/grid-list";
import { MatInputModule } from "@angular/material/input";
import { BrowserDynamicTestingModule } from "@angular/platform-browser-dynamic/testing";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { StatusesService } from "src/app/api/statuses.service";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { SharedModule } from "src/app/shared/shared.module";
import { StatusDetailsComponent } from "./status-details.component";

const status = { id: 1, name: 'test', description: 'test', lastUpdated: new Date('02/02/2023') };

describe("StatusDetailsComponent", () => {
	let component: StatusDetailsComponent;
	let fixture: ComponentFixture<StatusDetailsComponent>;
	let router: Router;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			imports: [
				HttpClientTestingModule,
				RouterTestingModule.withRoutes([
					{ path: 'core/statuses/:id', component: StatusDetailsComponent }
				]),
				ReactiveFormsModule,
				MatFormFieldModule,
				MatInputModule,
				MatGridListModule,
				MatCardModule,
				MatButtonModule,
				SharedModule
			],
			declarations: [StatusDetailsComponent, DecisionDialogComponent],
			providers: [StatusesService]
		})
			.compileComponents();
		TestBed.overrideModule(BrowserDynamicTestingModule, {
			set: {
				entryComponents: [DecisionDialogComponent]
			}
		});
		router = TestBed.inject(Router);
		fixture = TestBed.createComponent(StatusDetailsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("submits a update status request", fakeAsync(() => {
		const service = TestBed.inject(StatusesService);
		component.statusDetailsForm.setValue(status);
		spyOn(service, 'updateStatus').and.returnValue(Promise.resolve(status));
		component.updateStatus();

		service.updateStatus(component.statusDetailsForm.value, 1).then((result) => {
			expect(result).toEqual(status);
		})
	}));

	it("submits a status creation request", fakeAsync(() => {
		const service = TestBed.inject(StatusesService);
		component.statusDetailsForm.setValue(status);
		spyOn(service, 'createStatus').and.returnValue(Promise.resolve(status));
		component.createStatus();

		service.createStatus(component.statusDetailsForm.value).then((result) => {
			expect(result).toEqual(status);
			router.navigate(['/core/statuses/1']);
		})
	}));

});
