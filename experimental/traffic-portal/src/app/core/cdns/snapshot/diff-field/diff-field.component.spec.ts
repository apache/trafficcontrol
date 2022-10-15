import { ComponentFixture, TestBed } from "@angular/core/testing";

import { DiffFieldComponent } from "./diff-field.component";

describe("DiffFieldComponent", () => {
	let component: DiffFieldComponent;
	let fixture: ComponentFixture<DiffFieldComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ DiffFieldComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(DiffFieldComponent);
		component = fixture.componentInstance;
		component.value = {
			newValue: undefined,
			oldValue: undefined,
		};
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
