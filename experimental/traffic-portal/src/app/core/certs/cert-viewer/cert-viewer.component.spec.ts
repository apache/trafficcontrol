import { ComponentFixture, TestBed } from "@angular/core/testing";

import { CertViewerComponent } from "./cert-viewer.component";

describe("CertViewerComponent", () => {
	let component: CertViewerComponent;
	let fixture: ComponentFixture<CertViewerComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ CertViewerComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(CertViewerComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
