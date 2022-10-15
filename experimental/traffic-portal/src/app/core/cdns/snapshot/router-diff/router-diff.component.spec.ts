import { ComponentFixture, TestBed } from "@angular/core/testing";
import { of } from "rxjs";

import { RouterDiffComponent } from "./router-diff.component";

describe("RouterDiffComponent", () => {
	let component: RouterDiffComponent;
	let fixture: ComponentFixture<RouterDiffComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ RouterDiffComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(RouterDiffComponent);
		component = fixture.componentInstance;
		component.snapshots = of({current: {}, pending: {}});
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
