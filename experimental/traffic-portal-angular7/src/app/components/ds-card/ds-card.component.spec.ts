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

import { DsCardComponent } from './ds-card.component';
import { LoadingComponent } from '../loading/loading.component';
import { DeliveryService } from '../../models/deliveryservice';
import { LinechartDirective } from '../../directives/linechart.directive';

describe('DsCardComponent', () => {
	let component: DsCardComponent;
	let fixture: ComponentFixture<DsCardComponent>;

	beforeEach(async(() => {
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
		})
		.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(DsCardComponent);
		component = fixture.componentInstance;
		component.deliveryService = {} as DeliveryService;
		fixture.detectChanges();
	});

	it('should exist', () => {
		expect(component).toBeTruthy();
	});
});
