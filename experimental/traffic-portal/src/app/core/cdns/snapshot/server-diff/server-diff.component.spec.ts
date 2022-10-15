import { ComponentFixture, TestBed } from "@angular/core/testing";
import { of } from "rxjs";

import { ServerDiffComponent } from "./server-diff.component";

describe("ServerDiffComponent", () => {
	let component: ServerDiffComponent;
	let fixture: ComponentFixture<ServerDiffComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ ServerDiffComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(ServerDiffComponent);
		component = fixture.componentInstance;
		component.snapshots = of({current: {}, pending: {}});
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
