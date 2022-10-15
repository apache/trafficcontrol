import { ComponentFixture, TestBed } from "@angular/core/testing";
import { RouterTestingModule } from "@angular/router/testing";

import { APITestingModule } from "src/app/api/testing";

import { SnapshotComponent } from "./snapshot.component";

describe("SnapshotComponent", () => {
	let component: SnapshotComponent;
	let fixture: ComponentFixture<SnapshotComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ SnapshotComponent ],
			imports: [RouterTestingModule, APITestingModule]
		})
			.compileComponents();

		fixture = TestBed.createComponent(SnapshotComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
