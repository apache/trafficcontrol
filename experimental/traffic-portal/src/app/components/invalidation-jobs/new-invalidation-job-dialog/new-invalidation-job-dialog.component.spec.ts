import { HttpClientModule } from "@angular/common/http";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule, MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";

import { NewInvalidationJobDialogComponent } from "./new-invalidation-job-dialog.component";

describe("NewInvalidationJobDialogComponent", () => {
	let component: NewInvalidationJobDialogComponent;
	let fixture: ComponentFixture<NewInvalidationJobDialogComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ NewInvalidationJobDialogComponent ],
			imports: [
				MatDialogModule,
				HttpClientModule
			],
			providers: [
				{provide: MatDialogRef, useValue: {close: (): void => {
					console.log("dialog closed");
				}}},
				{provide: MAT_DIALOG_DATA, useValue: -1}
			]
		}).compileComponents();
	});

	beforeEach(() => {
		fixture = TestBed.createComponent(NewInvalidationJobDialogComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
