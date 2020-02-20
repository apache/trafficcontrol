import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';

import { TpHeaderComponent } from '../tp-header/tp-header.component';

import { ServersPageComponent } from './servers-page.component';

describe('ServersPageComponent', () => {
	let component: ServersPageComponent;
	let fixture: ComponentFixture<ServersPageComponent>;

	beforeEach(async(() => {
		TestBed.configureTestingModule({
			declarations: [
				ServersPageComponent,
				TpHeaderComponent
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
		fixture = TestBed.createComponent(ServersPageComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it('should create', () => {
		expect(component).toBeTruthy();
	});
});
