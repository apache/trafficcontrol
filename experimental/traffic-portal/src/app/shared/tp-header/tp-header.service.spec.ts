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
import { TestBed } from "@angular/core/testing";

import {TpHeaderComponent} from "src/app/shared/tp-header/tp-header.component";

import {HeaderNavigation, TpHeaderService} from "./tp-header.service";

describe("TpHeaderService", () => {
	let service: TpHeaderService;
	let mockHeaderComp: jasmine.SpyObj<TpHeaderComponent>;

	beforeEach(() => {
		mockHeaderComp = jasmine.createSpyObj<TpHeaderComponent>([], {hidden: false, title: ""});
		TestBed.configureTestingModule({
			providers: [
				TpHeaderService
			]
		});
		service = TestBed.inject(TpHeaderService);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("set header component", () => {
		expect(mockHeaderComp).toBeTruthy();
		expect(mockHeaderComp?.hidden).toBeFalse();
		expect(mockHeaderComp?.title).toBe("");

		service.headerHidden.next(true);
		service.headerTitle.next("something else");
		service.horizontalNavsUpdated.next(new Array<HeaderNavigation>());
		service.verticalNavsUpdated.next(new Array<HeaderNavigation>());
	});
});
