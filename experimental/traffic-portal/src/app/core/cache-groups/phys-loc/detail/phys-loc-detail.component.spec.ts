import { ComponentFixture, TestBed } from "@angular/core/testing";

import { PhysLocDetailComponent } from "src/app/core/cache-groups/phys-loc/detail/phys-loc-detail.component";

describe("PhysLocDetailComponent", () => {
	let component: PhysLocDetailComponent;
	let fixture: ComponentFixture<PhysLocDetailComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ PhysLocDetailComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(PhysLocDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
