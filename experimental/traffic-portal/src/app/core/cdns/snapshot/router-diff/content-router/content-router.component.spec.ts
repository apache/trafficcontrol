import { ComponentFixture, TestBed } from "@angular/core/testing";

import { ContentRouterComponent } from "./content-router.component";

describe("ContentRouterComponent", () => {
	let component: ContentRouterComponent;
	let fixture: ComponentFixture<ContentRouterComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ ContentRouterComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(ContentRouterComponent);
		component = fixture.componentInstance;
		component.router = {
			// eslint-disable-next-line @typescript-eslint/naming-convention
			"api.port": "",
			fqdn: "",
			httpsPort: -1,
			ip: "",
			ip6: "",
			location: "",
			port: -1,
			profile: "",
			// eslint-disable-next-line @typescript-eslint/naming-convention
			"secure.api.port": "",
			status: ""
		};
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
