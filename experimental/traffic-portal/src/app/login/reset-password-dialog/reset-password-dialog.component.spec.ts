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
import { type ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogRef } from "@angular/material/dialog";

import { UserService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";

import { ResetPasswordDialogComponent } from "./reset-password-dialog.component";

describe("ResetPasswordDialogComponent", () => {
	let component: ResetPasswordDialogComponent;
	let fixture: ComponentFixture<ResetPasswordDialogComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ ResetPasswordDialogComponent ],
			imports: [
				APITestingModule,
				HttpClientModule
			],
			providers: [
				// The controller doesn't pass any arguments or check the return
				// value of `close` - this literally needs to do nothing but be
				// callable.
				// eslint-disable-next-line @typescript-eslint/no-empty-function
				{ provide: MatDialogRef, useValue: {close: (): void => {}} },
			]
		})
			.compileComponents();
	});

	beforeEach(() => {
		fixture = TestBed.createComponent(ResetPasswordDialogComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("submits a password reset request", () => {
		const spy = spyOn(TestBed.inject(UserService), "resetPassword");
		expect(spy).not.toHaveBeenCalled();
		component.submit(new SubmitEvent("submit"));
		expect(spy).toHaveBeenCalled();
	});
});
