import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";

import { LastDaysComponent } from "./last-days.component";

describe("LastDaysComponent", () => {
	let component: LastDaysComponent;
	let fixture: ComponentFixture<LastDaysComponent>;
	let mockMatDialog: jasmine.SpyObj<MatDialogRef<number>>;

	beforeEach(async () => {
		mockMatDialog = jasmine.createSpyObj("MatDialogRef", ["close", "afterClosed"]);
		await TestBed.configureTestingModule({
			declarations: [LastDaysComponent],
			providers: [{provide: MatDialogRef, useValue: mockMatDialog},
				{provide: MAT_DIALOG_DATA, useValue: (): number => 3}]
		})
			.compileComponents();
	});

	beforeEach(() => {
		fixture = TestBed.createComponent(LastDaysComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
