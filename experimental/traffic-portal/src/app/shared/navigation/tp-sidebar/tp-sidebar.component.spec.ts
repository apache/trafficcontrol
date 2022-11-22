import { ComponentFixture, TestBed } from "@angular/core/testing";

import { TpSidebarComponent } from "./tp-sidebar.component";

describe("TpSidebarComponent", () => {
	let component: TpSidebarComponent;
	let fixture: ComponentFixture<TpSidebarComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ TpSidebarComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(TpSidebarComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
