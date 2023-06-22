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

import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { ReplaySubject } from "rxjs";

import { ProfileService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

import { ProfileDetailComponent } from "./profile-detail.component";

describe("ProfileDetailComponent", () => {
	let component: ProfileDetailComponent;
	let fixture: ComponentFixture<ProfileDetailComponent>;
	let route: ActivatedRoute;
	let paramMap: jasmine.Spy;
	let service: ProfileService;

	const navSvc = jasmine.createSpyObj([], {
		headerHidden: new ReplaySubject<boolean>(),
		headerTitle: new ReplaySubject<string>()
	});
	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ProfileDetailComponent],
			imports: [APITestingModule, RouterTestingModule, MatDialogModule],
			providers: [{provide: NavigationService, useValue: navSvc}]
		})
			.compileComponents();

		route = TestBed.inject(ActivatedRoute);
		paramMap = spyOn(route.snapshot.paramMap, "get");
		service = TestBed.inject(ProfileService);
		paramMap.and.returnValue(null);
		fixture = TestBed.createComponent(ProfileDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("new profile", async () => {
		paramMap.and.returnValue("new");

		fixture = TestBed.createComponent(ProfileDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.profile).not.toBeNull();
		expect(component.new).toBeTrue();
	});

	it("existing profile", async () => {
		const id = 1;
		paramMap.and.returnValue(id);
		const profile = await service.getProfiles(id);
		fixture = TestBed.createComponent(ProfileDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.profile).not.toBeNull();
		expect(component.profile.name).toBe(profile.name);
		expect(component.new).toBeFalse();
	});
});
