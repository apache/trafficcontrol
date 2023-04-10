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
import { ComponentFixture, fakeAsync, TestBed, tick } from "@angular/core/testing";
import { MatDialogRef } from "@angular/material/dialog";

import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";

import { UserRegistrationDialogComponent } from "./user-registration-dialog.component";

describe("UserRegistrationDialogComponent", () => {
	let component: UserRegistrationDialogComponent;
	let fixture: ComponentFixture<UserRegistrationDialogComponent>;
	let dialogOpen = true;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ UserRegistrationDialogComponent ],
			imports: [APITestingModule],
			providers: [
				{
					provide: CurrentUserService,
					useValue: {
						currentUser: {
							role: "admin",
							tenantId: 1
						},
						hasCapability: (): true => true,
					},
				},
				{
					provide: MatDialogRef,
					useValue: {
						close: (): void => {
							dialogOpen = false;
						}
					}
				}
			]
		})
			.compileComponents();

		fixture = TestBed.createComponent(UserRegistrationDialogComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("should close on success", fakeAsync(() => {
		component.submit(new Event("submit"));
		tick();
		expect(dialogOpen).toBeFalse();
	}));
});
