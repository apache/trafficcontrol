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
import { type ComponentFixture, TestBed, fakeAsync } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { RouterTestingModule } from "@angular/router/testing";
import {ReplaySubject} from "rxjs";

import { APITestingModule } from "src/app/api/testing";
import { TpHeaderService } from "src/app/shared/tp-header/tp-header.service";

import { CacheGroupTableComponent } from "./cache-group-table.component";

describe("CacheGroupTableComponent", () => {
	let component: CacheGroupTableComponent;
	let fixture: ComponentFixture<CacheGroupTableComponent>;

	const headerSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});
	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ CacheGroupTableComponent ],
			imports: [
				APITestingModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule
			],
			providers: [
				{ provide: TpHeaderService, useValue: headerSvc}
			]
		}).compileComponents();
		fixture = TestBed.createComponent(CacheGroupTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("emits the search box value", fakeAsync(() => {
		component.fuzzControl.setValue("query");
		component.updateURL();
		expectAsync(component.fuzzySubject.toPromise()).toBeResolvedTo("query");
	}));

	it("doesn't throw errors when handling context menu events", () => {
		expect(()=>component.handleContextMenu({action: "something", data: []})).not.toThrow();
	});
});
