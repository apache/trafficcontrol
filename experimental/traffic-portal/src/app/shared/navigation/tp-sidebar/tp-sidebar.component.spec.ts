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
import { HarnessLoader } from "@angular/cdk/testing";
import { TestbedHarnessEnvironment } from "@angular/cdk/testing/testbed";
import { HttpClientModule } from "@angular/common/http";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatMenuModule } from "@angular/material/menu";
import { MatSidenavModule } from "@angular/material/sidenav";
import { MatSidenavHarness } from "@angular/material/sidenav/testing";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { RouterTestingModule } from "@angular/router/testing";
import { of, ReplaySubject } from "rxjs";

import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { NavigationService, TreeNavNode } from "src/app/shared/navigation/navigation.service";

import { TpSidebarComponent } from "./tp-sidebar.component";

describe("TpSidebarComponent", () => {
	let component: TpSidebarComponent;
	let fixture: ComponentFixture<TpSidebarComponent>;
	let navSvc: jasmine.SpyObj<NavigationService>;
	let loader: HarnessLoader;

	beforeEach(async () => {
		const mockCurrentUserService = jasmine.createSpyObj(
			["updateCurrentUser", "hasPermission", "login", "logout"], {userChanged: of(null)});
		navSvc = jasmine.createSpyObj(["addHorizontalNav", "addVerticalNav"],
			{
				headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>(),
				sidebarHidden: new ReplaySubject<boolean>(), sidebarNavs: new ReplaySubject<Array<TreeNavNode>>()
			});
		await TestBed.configureTestingModule({
			declarations: [TpSidebarComponent],
			imports: [APITestingModule,
				HttpClientModule,
				RouterTestingModule,
				MatMenuModule,
				MatSidenavModule,
				NoopAnimationsModule
			],
			providers: [
				{provide: CurrentUserService, useValue: mockCurrentUserService},
				{provide: NavigationService, useValue: navSvc}
			]
		})
			.compileComponents();

		fixture = TestBed.createComponent(TpSidebarComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		loader = TestbedHarnessEnvironment.documentRootLoader(fixture);
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("node handle set correctly", () => {
		expect(component.nodeHandle({name: "test"})).toBe("test");
		expect(component.nodeHandle({href: "href", name: "test"})).toBe("testhref");
	});

	it("handles visibility events", async () => {
		const sidenavContainer = await loader.getAllHarnesses(MatSidenavHarness);
		expect(sidenavContainer.length).toBe(1);
		const sidenav = sidenavContainer[0];
		const openedSpy = spyOn(component.sidenav, "open");
		const closedSpy = spyOn(component.sidenav, "close");
		openedSpy.and.callThrough();
		closedSpy.and.callThrough();
		// default state
		expect(await sidenav.isOpen()).toBeFalse();

		navSvc.sidebarHidden.next(false);
		await fixture.whenStable();
		expect(openedSpy).toHaveBeenCalled();
		expect(await sidenav.isOpen()).toBeTrue();

		navSvc.sidebarHidden.next(false);
		await fixture.whenStable();
		expect(openedSpy).toHaveBeenCalledTimes(1);
		expect(await sidenav.isOpen()).toBeTrue();

		navSvc.sidebarHidden.next(true);
		await fixture.whenStable();
		expect(closedSpy).toHaveBeenCalled();
		expect(await sidenav.isOpen()).toBeFalse();

		navSvc.sidebarHidden.next(true);
		await fixture.whenStable();
		expect(closedSpy).toHaveBeenCalledTimes(1);
		expect(await sidenav.isOpen()).toBeFalse();
	});

	it("child to parent association", async () => {
		const navs = [{
			name: "first"
		}, {
			children: [{
				href: "other",
				name: "first"
			}],
			name: "second"
		}];
		navSvc.sidebarNavs.next(navs);
		expect(component.dataSource.data).toEqual(navs);
		expect(component.isRoot(navs[0])).toBeTrue();
		expect(component.isRoot(navs[1])).toBeTrue();
		expect(component.isRoot((navs[1].children ?? [])[0])).toBeFalse();

		expect(component.isRoot({name: "madeup"})).toBeTrue();
	});

	it("collapse all nodes", () => {
		const collapseSpy = spyOn(component.treeCtrl, "collapseAll");
		component.collapseAll();
		expect(collapseSpy).toHaveBeenCalled();
	});

	it("node has children", () => {
		expect(component.hasChild(0, {name: "test"})).toBeFalse();
		expect(component.hasChild(0, {children: [], name: "test"})).toBeFalse();
		expect(component.hasChild(0, {children: [{name: "other"}], name: "test"})).toBeTrue();
	});
});
