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
import { type ComponentFixture, TestBed, fakeAsync, tick } from "@angular/core/testing";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import {ReplaySubject, Subject} from "rxjs";
import { ResponseCurrentUser } from "trafficops-types";

import { UserService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";
import { TpHeaderComponent } from "src/app/shared/navigation/tp-header/tp-header.component";

import { CurrentuserComponent } from "./currentuser.component";

describe("CurrentuserComponent", () => {
	let component: CurrentuserComponent;
	let fixture: ComponentFixture<CurrentuserComponent>;
	let dialogClose: Subject<void>;
	let updateSpy: jasmine.Spy;
	let updateSucceeded: boolean;
	let currentUser: null | ResponseCurrentUser = null;
	let api: UserService;

	beforeEach(fakeAsync(() => {
		updateSucceeded = false;
		updateSpy = jasmine.createSpy("CurrentUserService's `updateCurrentuser` method", async () => {
			currentUser = await api.getCurrentUser();
			return updateSucceeded;
		});
		updateSpy.and.callThrough();

		dialogClose = new Subject();

		const headerSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});

		TestBed.configureTestingModule({
			declarations: [
				CurrentuserComponent,
				TpHeaderComponent
			],
			imports: [
				APITestingModule,
				HttpClientModule,
				RouterTestingModule.withRoutes([{component: CurrentuserComponent, path: ""}]),
				MatDialogModule
			],
			providers: [
				{
					provide: CurrentUserService,
					useValue: {
						get currentUser(): ResponseCurrentUser | null {
							return currentUser;
						},
						hasPermission: (): boolean => true,
						updateCurrentUser: async (): Promise<boolean> => updateSpy()
					}
				},
				{
					provide: MatDialog,
					useValue: {
						open: (): {afterClosed: () => Subject<void>} => ({
							afterClosed: () => dialogClose
						})
					}
				},
				{ provide: NavigationService, useValue: headerSvc}
			]
		});
		TestBed.compileComponents();
		api = TestBed.inject(UserService);
		api.getCurrentUser().then(
			u => {
				currentUser = u;
			}
		);
		tick();
		fixture = TestBed.createComponent(CurrentuserComponent);
		component = fixture.componentInstance;
		TestBed.inject(Router).initialNavigation();
		fixture.detectChanges();
		tick();
	}));

	it("should create", fakeAsync(() => {
		updateSucceeded = true;
		expect(component).toBeTruthy();
		component.currentUser = null;
		component.ngOnInit();
		tick();
		expect(updateSpy).toHaveBeenCalled();
		expect(component.currentUser).not.toBeNull();
	}));

	it("toggles editing mode", () => {
		expect(component.editMode).toBeFalse();
		expect(component.editUser).toBeNull();
		component.edit();
		expect(component.editMode).toBeTrue();
		expect(component.editUser).toEqual(component.currentUser);
		component.cancelEdit();
		expect(component.editMode).toBeFalse();
		component.currentUser = null;
		expect(()=>component.edit()).toThrow();
		expect(()=>component.cancelEdit()).toThrow();
	});

	it("sets edit mode from query parameters", fakeAsync(()=>{
		const router = TestBed.inject(Router);
		router.navigate(["."], {queryParams: {edit: true}});
		tick();
		component.ngOnInit();
		expect(component.editMode).toBeTrue();
		router.navigate(["."], {queryParams: {edit: true, updatePassword: true}});
		tick();
		component.ngOnInit();
		expect(component.editMode).toBeTrue();
		tick();
		expect(router.url).toContain("edit=true");
		expect(router.url).toContain("updatePassword=true");
		dialogClose.next(void undefined);
		tick();
		expect(router.url).toBe("/?edit=true");
	}));

	it("submits user update requests", fakeAsync(()=>{
		if (component.currentUser === null) {
			return fail("component initialized with null current User");
		}

		expectAsync(component.submitEdit(new Event("submit"))).toBeRejected();

		component.editUser = {
			...component.currentUser,
		};
		component.submitEdit(new Event("submit"));
		tick();

		expect(updateSpy).toHaveBeenCalledTimes(1);

		updateSucceeded = true;
		component.submitEdit(new Event("submit"));
		tick();

		expect(updateSpy).toHaveBeenCalledTimes(2);

		component.editUser = {
			...component.editUser,
			id: -1
		};
		component.submitEdit(new Event("submit"));
		tick();

		expect(updateSpy).toHaveBeenCalledTimes(2);
	}));

	it("determines whether a user has a 'bottom-level' address", () => {
		component.currentUser = null;
		expect(component.hasBottomAddress()).toBeFalse();
		component.currentUser = {
			addressLine1: null,
			addressLine2: null,
			changeLogCount: 0,
			city: null,
			company: null,
			country: null,
			email: "em@i.l",
			fullName: "",
			gid: null,
			id: 1,
			lastAuthenticated: null,
			lastUpdated: new Date(),
			localUser: false,
			newUser: false,
			phoneNumber: null,
			postalCode: null,
			publicSshKey: null,
			registrationSent: null,
			role: "",
			stateOrProvince: null,
			tenant: "",
			tenantId: 1,
			ucdn: "",
			uid: null,
			username: "",
		};
		expect(component.hasBottomAddress()).toBeFalse();
		component.currentUser.city = "Townsville";
		expect(component.hasBottomAddress()).toBeTrue();
		component.currentUser.city = null;
		component.currentUser.country = "Nationstan";
		expect(component.hasBottomAddress()).toBeTrue();
		component.currentUser.country = null;
		component.currentUser.stateOrProvince = "Provincia";
		expect(component.hasBottomAddress()).toBeTrue();
		component.currentUser.stateOrProvince = null;
		component.currentUser.postalCode = "00000";
		expect(component.hasBottomAddress()).toBeTrue();
	});
});
