import { ComponentFixture, TestBed } from "@angular/core/testing";

import { PhysLocTableComponent } from "src/app/core/cache-groups/phys-loc/table/phys-loc-table.component";

describe("PhysLocTableComponent", () => {
	let component: PhysLocTableComponent;
	let fixture: ComponentFixture<PhysLocTableComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ PhysLocTableComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(PhysLocTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
