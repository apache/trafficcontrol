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

import { UserCardComponent } from './user-card.component';
import { User } from '../../models';

describe('UserCardComponent', () => {
	let component: UserCardComponent;
	let fixture: ComponentFixture<UserCardComponent>;

	beforeEach(async(() => {
		TestBed.configureTestingModule({
			declarations: [ UserCardComponent ],
			imports: [
				HttpClientModule
			]
		})
		.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(UserCardComponent);
		component = fixture.componentInstance;
		component.user = {lastUpdated: new Date(), id: 1, name: 'test', username: 'test', newUser: false} as User;
		fixture.detectChanges();
	});

	it('should exist', () => {
		expect(component).toBeTruthy();
	});

	afterAll(() => {
		TestBed.resetTestingModule();
	});
});
