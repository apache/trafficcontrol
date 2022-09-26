import { ComponentFixture, TestBed } from "@angular/core/testing";

import { TextDialogComponent } from "./text-dialog.component";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";

describe("TextDialogComponent", () => {
	let component: TextDialogComponent;
	let fixture: ComponentFixture<TextDialogComponent>;

	beforeEach(async () => {
		const mockMatDialog = jasmine.createSpyObj("MatDialogRef", ["close", "afterClosed"]);
		await TestBed.configureTestingModule({
			declarations: [ TextDialogComponent ],
			providers: [{provide: MatDialogRef, useValue: mockMatDialog},
				{provide: MAT_DIALOG_DATA, useValue: () => {}}]
		})
			.compileComponents();

		fixture = TestBed.createComponent(TextDialogComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
