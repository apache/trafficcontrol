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
import { MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { RouterTestingModule } from "@angular/router/testing";
import { ResponseCurrentUser } from "trafficops-types";

import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";

import { UpdatePasswordDialogComponent } from "./update-password-dialog.component";

describe("UpdatePasswordDialogComponent", () => {
	let component: UpdatePasswordDialogComponent;
	let fixture: ComponentFixture<UpdatePasswordDialogComponent>;
	let dialogOpen = true;
	let updated = false;

	const mockAPIService = jasmine.createSpyObj(["updateCurrentUser", "getCurrentUser", "saveCurrentUser"], );
	mockAPIService.updateCurrentUser.and.returnValue(new Promise(resolve => resolve(true)));
	mockAPIService.getCurrentUser.and.returnValue(
		new Promise<ResponseCurrentUser>(
			resolve => resolve({id: -1, newUser: false, username: ""} as ResponseCurrentUser)
		)
	);
	mockAPIService.currentUser = {id: 1, newUser: false, username: "hello"};

	beforeEach(async () => {
		dialogOpen = true;
		updated = false;
		await TestBed.configureTestingModule({
			declarations: [ UpdatePasswordDialogComponent ],
			imports: [ APITestingModule, HttpClientModule, MatDialogModule, RouterTestingModule ],
			providers: [
				{provide: MatDialogRef, useValue: {close: (upd?: true): void => {
					dialogOpen = false;
					updated = upd ?? false;
				}}},
				{provide: CurrentUserService, useValue: mockAPIService}
			]
		}).compileComponents();
	});

	beforeEach(() => {
		fixture = TestBed.createComponent(UpdatePasswordDialogComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("doesn't allow submitting mismatched passwords", async () => {
		component.password = "password";
		component.confirm = "mismatch";
		await component.submit(new Event("submit"));
		expect(dialogOpen).toBeTrue();
		expect(updated).toBeFalse();
		component.confirmValid.subscribe(
			v => {
				expect(v).toBeTruthy();
			}
		);
		component.confirm = component.password;
		mockAPIService.saveCurrentUser.and.returnValue(new Promise(r => r(true)));
		await component.submit(new Event("submit"));
		expect(dialogOpen).toBeFalse();
		expect(updated).toBeTrue();
	});

	it("closes the dialog on cancel", () => {
		component.cancel();
		expect(dialogOpen).toBeFalse();
		expect(updated).toBeFalse();
	});
});
