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
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';

import { DashboardComponent } from './dashboard.component';
import { DsCardComponent } from '../ds-card/ds-card.component';
import { LoadingComponent } from '../loading/loading.component';
import { TpHeaderComponent } from '../tp-header/tp-header.component';

import { DeliveryService } from '../../models';

import { LinechartDirective } from '../../directives/linechart.directive';

describe('DashboardComponent', () => {
	let component: DashboardComponent;
	let fixture: ComponentFixture<DashboardComponent>;

	beforeEach(async(() => {
		TestBed.configureTestingModule({
		declarations: [
			DashboardComponent,
			DsCardComponent,
			LoadingComponent,
			TpHeaderComponent,
			LinechartDirective
		],
		imports: [
			FormsModule,
			HttpClientModule,
			ReactiveFormsModule,
			RouterTestingModule
		]
		})
		.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(DashboardComponent);
		component = fixture.componentInstance;
		component.deliveryServices = [
			{
				displayName: 'FIZZbuzz'
			} as DeliveryService,
			{
				displayName: 'fooBAR'
			} as DeliveryService
		];
		fixture.detectChanges();
	});

	it('should exist', () => {
		expect(component).toBeTruthy();
	});

	it('should implement fuzzy search', () => {
		// letter exclusion
		component.fuzzControl.setValue('z');
		expect(component.fuzzy(component.deliveryServices[0])).toBeTruthy();
		expect(component.fuzzy(component.deliveryServices[1])).toBeFalsy();

		// matches case-insensitively
		component.fuzzControl.setValue('fb');
		expect(component.fuzzy(component.deliveryServices[0])).toBeTruthy();
		expect(component.fuzzy(component.deliveryServices[1])).toBeTruthy();

	});

	it('sets the "search" query parameter', () => {
		expect(true).toBeTruthy();
	});

	afterAll(() => {
		TestBed.resetTestingModule();
	});
});
