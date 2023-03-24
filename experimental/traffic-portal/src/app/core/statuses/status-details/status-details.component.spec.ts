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
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { RouterTestingModule } from "@angular/router/testing";
import { Observable, ReplaySubject, of } from "rxjs";

import { ServerService } from "src/app/api";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

import { StatusDetailsComponent } from "./status-details.component";

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
			afterClosed: (): Observable<boolean> => of(true)
		};
	}
}

describe("StatusDetailsComponent", () => {
	let component: StatusDetailsComponent;
	let fixture: ComponentFixture<StatusDetailsComponent>;

	beforeEach(async () => {

		const navSvc = jasmine.createSpyObj([], { headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>() });

		await TestBed.configureTestingModule({
			declarations: [StatusDetailsComponent],
			imports: [
				HttpClientModule,
				RouterTestingModule,
				FormsModule,
				ReactiveFormsModule
			],
			providers: [
				{ provide: MatDialog, useClass: MockDialog },
				{ provide: NavigationService, useValue: navSvc },
				ServerService
			]
		})
			.compileComponents();
		fixture = TestBed.createComponent(StatusDetailsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

});
