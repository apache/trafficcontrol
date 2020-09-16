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
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientModule } from '@angular/common/http';
import { RouterTestingModule } from '@angular/router/testing';
import { of } from 'rxjs';

import { DsCardComponent } from './ds-card.component';
import { LoadingComponent } from '../loading/loading.component';

import { LinechartDirective } from '../../directives/linechart.directive';

import { DeliveryService, GeoLimit, GeoProvider } from '../../models';

import { APIService } from '../../services/api.service';

describe('DsCardComponent', () => {
	let component: DsCardComponent;
	let fixture: ComponentFixture<DsCardComponent>;

	beforeEach(async(() => {
		// mock the API
		const mockAPIService = jasmine.createSpyObj(['getDSKBPS', 'getDSCapacity', 'getDSHealth']);
		mockAPIService.getDSKBPS.and.returnValue(of([]), of([]));
		mockAPIService.getDSCapacity.and.returnValue(of({
			availablePercent: 34,
			maintenance: 42,
			utilized: 24
		}));
		mockAPIService.getDSHealth.and.returnValue(of({
			totalOnline: 80,
			totalOffline: 20
		}));

		TestBed.configureTestingModule({
			declarations: [
				DsCardComponent,
				LoadingComponent,
				LinechartDirective
			],
			imports: [
				HttpClientModule,
				RouterTestingModule
			]
		});
		TestBed.overrideProvider(APIService, { useValue: mockAPIService });
		TestBed.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(DsCardComponent);
		component = fixture.componentInstance;
		component.deliveryService = {xmlId: "test-ds"} as DeliveryService;
		fixture.detectChanges();
	});

	it('should exist', () => {
		expect(component).toBeTruthy();
	});

	afterAll(() => {
		TestBed.resetTestingModule();
	});
});
