import { waitForAsync, ComponentFixture, TestBed } from "@angular/core/testing";

import { ServersTableComponent } from "./servers-table.component";

import { TpHeaderComponent } from "../../tp-header/tp-header.component";

describe("ServersTableComponent", () => {
	let component: ServersTableComponent;
	let fixture: ComponentFixture<ServersTableComponent>;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [ ServersTableComponent, TpHeaderComponent ]
		})
		.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(ServersTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
