import { ComponentFixture, TestBed } from "@angular/core/testing";

import { ContentServerComponent } from "./content-server.component";

describe("ContentServerComponent", () => {
	let component: ContentServerComponent;
	let fixture: ComponentFixture<ContentServerComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ ContentServerComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(ContentServerComponent);
		component = fixture.componentInstance;
		component.server = {
			cacheGroup: "",
			capabilities: [],
			fqdn: "",
			hashCount: -1,
			hashId: "",
			httpsPort: -1,
			interfaceName: "",
			ip: "",
			ip6: "",
			locationId: "",
			port: -1,
			profile: "",
			routingDisabled: 0,
			status: "",
			type: "MID"
		};
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
