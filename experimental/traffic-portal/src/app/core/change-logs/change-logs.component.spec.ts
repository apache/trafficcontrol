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
import { MatDialog } from "@angular/material/dialog";
import { RouterTestingModule } from "@angular/router/testing";
import { Observable, of, ReplaySubject } from "rxjs";

import { APITestingModule } from "src/app/api/testing";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

import { ChangeLogsComponent } from "./change-logs.component";

/**
 * Define the MockDialog
 */
class MockDialog {

	/**
	 * Fake opens the dialog
	 *
	 * @returns unknown
	 */
	public open(): unknown {
		return {
			afterClosed: (): Observable<number> => of(3)
		};
	}
}
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
				RouterTestingModule
			],
			providers: [
				{ provide: MatDialog, useClass: MockDialog },
				{ provide: NavigationService, useValue: navSvc },
			]
		})
			.compileComponents();
	});

	beforeEach(() => {
		fixture = TestBed.createComponent(ChangeLogsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
