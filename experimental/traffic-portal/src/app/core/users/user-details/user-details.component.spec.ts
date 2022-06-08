import { ComponentFixture, TestBed } from "@angular/core/testing";
import { RouterTestingModule } from "@angular/router/testing";

import { UserService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";

import { UserDetailsComponent } from "./user-details.component";

describe("UserDetailsComponent", () => {
	let component: UserDetailsComponent;
	let fixture: ComponentFixture<UserDetailsComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ UserDetailsComponent ],
			imports: [ APITestingModule, RouterTestingModule ]
		}).compileComponents();
		fixture = TestBed.createComponent(UserDetailsComponent);
		component = fixture.componentInstance;
		const service = TestBed.inject(UserService);
		component.roles = await service.getRoles();
		component.tenants = await service.getTenants();
		component.user = (await service.getUsers())[0];
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
