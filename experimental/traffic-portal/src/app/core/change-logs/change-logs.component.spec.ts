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
import { TestbedHarnessEnvironment } from "@angular/cdk/testing/testbed";
import { HttpClientModule } from "@angular/common/http";
import { ComponentFixture, fakeAsync, TestBed, tick } from "@angular/core/testing";
import { MatDialog, MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { MatDialogHarness } from "@angular/material/dialog/testing";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { of, ReplaySubject } from "rxjs";

import { APITestingModule } from "src/app/api/testing";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

import { ChangeLogsComponent } from "./change-logs.component";

const testCLEntry = {
	id: 1,
	lastUpdated: new Date(),
	level: "APICHANGE" as const,
	longTime: "",
	message: "testquest",
	relativeTime: "3 seconds ago",
	ticketNum: null,
	user: "admin"
};

describe("ChangeLogsComponent", () => {
	let component: ChangeLogsComponent;
	let fixture: ComponentFixture<ChangeLogsComponent>;

	beforeEach(async () => {
		const navSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});
		await TestBed.configureTestingModule({
			declarations: [ChangeLogsComponent],
			imports: [
				APITestingModule,
				HttpClientModule,
				RouterTestingModule,
				NoopAnimationsModule,
				MatDialogModule
			],
			providers: [
				{ provide: NavigationService, useValue: navSvc },
			]
		}).compileComponents();

		fixture = TestBed.createComponent(ChangeLogsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("handles unknown actions", async () => {
		await expectAsync(component.handleContextMenu({action: "something unknown", data: []})).toBeResolvedTo(undefined);
		await expectAsync(component.handleTitleButton("something unknown")).toBeResolvedTo(undefined);
	});

	it("updates the fuzzy search output", fakeAsync(() => {
		let called = false;
		const text = "testquest";
		const spy = jasmine.createSpy("subscriber", (txt: string): void =>{
			if (!called) {
				expect(txt).toBe("");
				called = true;
			} else {
				expect(txt).toBe(text);
			}
		});
		component.fuzzySubj.subscribe(spy);
		tick();
		expect(spy).toHaveBeenCalled();
		component.searchText = text;
		component.updateURL();
		tick();
		expect(spy).toHaveBeenCalledTimes(2);
	}));

	it("sets the fuzzy search subject based on the search query param", fakeAsync(() => {
		const router = TestBed.inject(ActivatedRoute);
		const searchString = "testquest";
		spyOnProperty(router, "queryParamMap").and.returnValue(of(new Map([["search", searchString]])));

		let searchValue = "not the right string";
		component.fuzzySubj.subscribe(
			s => searchValue = s
		);

		component.ngOnInit();
		tick();

		expect(searchValue).toBe(searchString);
	}));

	it("opens a dialog to view changelog item details (array data)", fakeAsync(async () => {
		const asyncExpectation = expectAsync(
			component.handleContextMenu({action: "viewChangeLog", data: [testCLEntry]})
		).toBeResolvedTo(undefined);

		tick();

		const loader = TestbedHarnessEnvironment.documentRootLoader(fixture);

		const dialogs = await loader.getAllHarnesses(MatDialogHarness);
		if (dialogs.length !== 1) {
			return fail(`expected exactly one dialog to be opened; got: ${dialogs.length}`);
		}

		await dialogs[0].close();

		await asyncExpectation;
	}));

	it("opens a dialog to view changelog item details (single object data)", fakeAsync(async () => {
		const asyncExpectation = expectAsync(
			component.handleContextMenu({action: "viewChangeLog", data: testCLEntry})
		).toBeResolvedTo(undefined);

		tick();

		const loader = TestbedHarnessEnvironment.documentRootLoader(fixture);

		const dialogs = await loader.getAllHarnesses(MatDialogHarness);
		if (dialogs.length !== 1) {
			return fail(`expected exactly one dialog to be opened; got: ${dialogs.length}`);
		}

		await dialogs[0].close();

		await asyncExpectation;
	}));

	it("opens a dialog to set the number of days to which to filter the changelogs", fakeAsync(async () => {
		const numDays = 5;
		component.lastDays = numDays + 1;

		const dialogService = TestBed.inject(MatDialog);
		const openSpy = spyOn(dialogService, "open").and.returnValue({
			afterClosed: () => of(numDays)
		} as MatDialogRef<unknown>);
		expect(openSpy).not.toHaveBeenCalled();

		const asyncExpectation = expectAsync(component.handleTitleButton("lastDays")).toBeResolvedTo(undefined);
		tick();
		await asyncExpectation;
		expect(openSpy).toHaveBeenCalled();
		expect(component.lastDays).toBe(numDays);
	}));
});
