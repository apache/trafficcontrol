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

import { APITestingModule } from "src/app/api/testing";
import { ASNDetailComponent } from "src/app/core/cache-groups/asns/detail/asn-detail.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

describe("AsnDetailComponent", () => {
	let component: ASNDetailComponent;
	let fixture: ComponentFixture<ASNDetailComponent>;
	let route: ActivatedRoute;
	let paramMap: jasmine.Spy;

	const headerSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});
	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ ASNDetailComponent ],
			imports: [ APITestingModule, RouterTestingModule, MatDialogModule ],
			providers: [ { provide: NavigationService, useValue: headerSvc } ]
		})
			.compileComponents();

		route = TestBed.inject(ActivatedRoute);
		paramMap = spyOn(route.snapshot.paramMap, "get");
		paramMap.and.returnValue(null);
		fixture = TestBed.createComponent(ASNDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
		expect(paramMap).toHaveBeenCalled();
	});

	it("new asn", async () => {
		paramMap.and.returnValue("new");

		fixture = TestBed.createComponent(ASNDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.asn).not.toBeNull();
		expect(component.asn.asn).toBe(0);
		expect(component.isNew).toBeTrue();
	});

	it("existing asn", async () => {
		paramMap.and.returnValue("1");

		fixture = TestBed.createComponent(ASNDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.asn).not.toBeNull();
		expect(component.asn.asn).toBe(0);
		expect(component.isNew).toBeFalse();
	});
});
