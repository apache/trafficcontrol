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
import { HttpClientModule } from "@angular/common/http";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatMenuModule } from "@angular/material/menu";
import { RouterTestingModule } from "@angular/router/testing";
import { of, ReplaySubject } from "rxjs";

import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { NavigationService, TreeNavNode } from "src/app/shared/navigation/navigation.service";

import { TpSidebarComponent } from "./tp-sidebar.component";

describe("TpSidebarComponent", () => {
	let component: TpSidebarComponent;
	let fixture: ComponentFixture<TpSidebarComponent>;

	beforeEach(async () => {
		const mockCurrentUserService = jasmine.createSpyObj(
			["updateCurrentUser", "hasPermission", "login", "logout"], {userChanged: of(null)});
		const navSvc = jasmine.createSpyObj(["addHorizontalNav", "addVerticalNav"],
			{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>(),
				sidebarHidden: new ReplaySubject<boolean>(), sidebarNavs: new ReplaySubject<Array<TreeNavNode>>()});
		await TestBed.configureTestingModule({
			declarations: [ TpSidebarComponent ],
			imports: [ APITestingModule, HttpClientModule, RouterTestingModule, MatMenuModule ],
			providers: [
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				{ provide: NavigationService, useValue: navSvc}
			]
		})
			.compileComponents();

		fixture = TestBed.createComponent(TpSidebarComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("node has children", () => {
		expect(component.hasChild(0, {name: "test"})).toBeFalse();
		expect(component.hasChild(0, {children: [], name: "test"})).toBeFalse();
		expect(component.hasChild(0, {children: [{name: "other"}], name: "test"})).toBeTrue();
	});
});
