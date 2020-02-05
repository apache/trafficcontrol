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
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { of } from 'rxjs';

import { UsersComponent } from './users.component';

import { LoadingComponent } from '../loading/loading.component';
import { TpHeaderComponent } from '../tp-header/tp-header.component';
import { UserCardComponent } from '../user-card/user-card.component';

import { APIService } from '../../services/api.service';

import { User } from '../../models';

describe('UsersComponent', () => {
	let component: UsersComponent;
	let fixture: ComponentFixture<UsersComponent>;

	beforeEach(async(() => {
		// mock the API
		const mockAPIService = jasmine.createSpyObj(["getUsers", "getRoles", "getCurrentUser"]);
		mockAPIService.getUsers.and.returnValue(of([]));
		mockAPIService.getRoles.and.returnValue(of([]));
		mockAPIService.getCurrentUser.and.returnValue(of({
			id: 0,
			newUser: false,
			username: "test"
		} as User));

		TestBed.configureTestingModule({
			declarations: [
				UsersComponent,
				LoadingComponent,
				TpHeaderComponent,
				UserCardComponent
			],
			imports: [
				FormsModule,
				HttpClientModule,
				ReactiveFormsModule
			]
		});
		TestBed.overrideProvider(APIService, { useValue: mockAPIService });
		TestBed.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(UsersComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it('should exist', () => {
		expect(component).toBeTruthy();
	});

	afterAll(() => {
		TestBed.resetTestingModule();
	});
});
