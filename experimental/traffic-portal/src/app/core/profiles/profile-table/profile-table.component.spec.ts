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

import { ComponentFixture, TestBed, fakeAsync, tick } from "@angular/core/testing";
import { MatDialog, MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { of } from "rxjs";
import { ProfileType } from "trafficops-types";

import { ProfileService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { FileUtilsService } from "src/app/shared/file-utils.service";
import { isAction } from "src/app/shared/generic-table/generic-table.component";

import { ProfileTableComponent } from "./profile-table.component";

describe("ProfileTableComponent", () => {
	let component: ProfileTableComponent;
	let fixture: ComponentFixture<ProfileTableComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ProfileTableComponent],
			imports: [
				APITestingModule,
				RouterTestingModule,
				MatDialogModule
			],
			providers:[
				FileUtilsService
			]
		})
			.compileComponents();

		fixture = TestBed.createComponent(ProfileTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("sets the fuzzy search subject based on the search query param", fakeAsync(() => {
		const router = TestBed.inject(ActivatedRoute);
		const searchString = "testquest";
		spyOnProperty(router, "queryParamMap").and.returnValue(of(new Map([["search", searchString]])));

		let searchValue = "not the right string";
		component.fuzzySubject.subscribe(
			s => searchValue = s
		);

		component.ngOnInit();
		tick();

		expect(searchValue).toBe(searchString);
	}));

	it("updates the fuzzy search output", fakeAsync(() => {
		let called = false;
		const text = "testquest";
		const spy = jasmine.createSpy("subscriber", (txt: string): void => {
			if (!called) {
				expect(txt).toBe("");
				called = true;
			} else {
				expect(txt).toBe(text);
			}
		});
		component.fuzzySubject.subscribe(spy);
		tick();
		expect(spy).toHaveBeenCalled();
		component.fuzzControl.setValue(text);
		component.updateURL();
		tick();
		expect(spy).toHaveBeenCalledTimes(2);
	}));

	it("handles the 'delete' context menu item", fakeAsync(async () => {
		const item = component.contextMenuItems.find(c => c.name === "Delete");
		if (!item) {
			return fail("missing 'Delete' context menu item");
		}
		if (!isAction(item)) {
			return fail("expected an action, not a link");
		}
		expect(item.multiRow).toBeFalsy();
		expect(item.disabled).toBeUndefined();

		const api = TestBed.inject(ProfileService);
		const spy = spyOn(api, "deleteProfile").and.callThrough();
		expect(spy).not.toHaveBeenCalled();

		const dialogService = TestBed.inject(MatDialog);
		const openSpy = spyOn(dialogService, "open").and.returnValue({
			afterClosed: () => of(true)
		} as MatDialogRef<unknown>);

		const profile = await api.createProfile({
			cdn: 1,
			description: "blah",
			name: "test",
			routingDisabled: false,
			type: ProfileType.ATS_PROFILE
		});
		expect(openSpy).not.toHaveBeenCalled();
		const asyncExpectation = expectAsync(component.handleContextMenu({ action: "delete", data: profile })).toBeResolvedTo(undefined);
		tick();

		expect(openSpy).toHaveBeenCalled();
		tick();

		expect(spy).toHaveBeenCalled();

		await asyncExpectation;
	}));

	it("constructs profile export context menu links", async () => {
		const item = component.contextMenuItems.find(c => c.name === "Export Profile");
		if (!item) {
			return fail("missing 'Export Profile' context menu item");
		}
		if (isAction(item)) {
			return fail("expected a link, not an action");
		}
		expect(item.newTab).toBeTrue();
		expect(item.disabled).toBeUndefined();
		expect(item.fragment).toBeUndefined();
		expect(item.queryParams).toBeUndefined();

		if (typeof(item.href) !== "function") {
			return fail(`expected a functional href property, got: ${typeof(item.href)}`);
		}

		const api = TestBed.inject(ProfileService);
		const profile = await api.createProfile({
			cdn: 1,
			description: "blah",
			name: "test",
			routingDisabled: false,
			type: ProfileType.ATS_PROFILE
		});

		expect(item.href(profile)).toBe(`/api/${api.apiVersion}/profiles/${profile.id}/export`);
		await api.deleteProfile(profile);
	});
});
