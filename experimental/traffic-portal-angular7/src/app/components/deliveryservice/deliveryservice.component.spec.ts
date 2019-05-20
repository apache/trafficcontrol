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

import { DeliveryserviceComponent } from './deliveryservice.component';
import { TpHeaderComponent } from '../tp-header/tp-header.component';

import { LinechartDirective } from '../../directives/linechart.directive';


describe('DeliveryserviceComponent', () => {
	let component: DeliveryserviceComponent;
	let fixture: ComponentFixture<DeliveryserviceComponent>;

	beforeEach(async(() => {
		TestBed.configureTestingModule({
			declarations: [
				DeliveryserviceComponent,
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
		fixture = TestBed.createComponent(DeliveryserviceComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it('should exist', () => {
		expect(component).toBeTruthy();
	});

	it('sets the "to" and "from" values to "so far today"', () => {
		const now = new Date();
		now.setUTCMilliseconds(0);
		const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());

		expect(component.to).toEqual(now);
		expect(component.from).toEqual(today);

		let nowParts: any[] = now.toISOString().split('T');
		nowParts[1] = now.toTimeString().split(':');
		nowParts[1] = [nowParts[1][0], nowParts[1][1]].join(':');
		let todayParts: any[] = today.toISOString().split('T');
		todayParts[1] = today.toTimeString().split(':');
		todayParts[1] = [todayParts[1][0], todayParts[1][1]].join(':');
		expect(nowParts[0]).toEqual(component.toDate.value);
		expect(nowParts[1]).toEqual(component.toTime.value);
		expect(todayParts[0]).toEqual(component.fromDate.value);
		expect(todayParts[1]).toEqual(component.fromTime.value);
	});
});
