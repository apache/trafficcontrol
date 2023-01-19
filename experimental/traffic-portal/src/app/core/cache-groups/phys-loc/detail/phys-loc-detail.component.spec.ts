import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { ReplaySubject } from "rxjs";

import { APITestingModule } from "src/app/api/testing";
import { PhysLocDetailComponent } from "src/app/core/cache-groups/phys-loc/detail/phys-loc-detail.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

describe("PhysLocDetailComponent", () => {
	let component: PhysLocDetailComponent;
	let fixture: ComponentFixture<PhysLocDetailComponent>;
	let route: ActivatedRoute;
	let paramMap: jasmine.Spy;

	const navSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});
	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ PhysLocDetailComponent ],
			imports: [ APITestingModule, RouterTestingModule, MatDialogModule ],
			providers: [ { provide: NavigationService, useValue: navSvc } ]
		})
			.compileComponents();

		route = TestBed.inject(ActivatedRoute);
		paramMap = spyOn(route.snapshot.paramMap, "get");
		paramMap.and.returnValue(null);
		fixture = TestBed.createComponent(PhysLocDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
		expect(paramMap).toHaveBeenCalled();
	});

	it("new physicalLocation", async () => {
		paramMap.and.returnValue("new");

		fixture = TestBed.createComponent(PhysLocDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.physLocation).not.toBeNull();
		expect(component.physLocation.name).toBe("");
		expect(component.new).toBeTrue();
	});

	it("existing physicalLocation", async () => {
		paramMap.and.returnValue("1");

		fixture = TestBed.createComponent(PhysLocDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.physLocation).not.toBeNull();
		expect(component.physLocation.name).toBe("phys");
		expect(component.new).toBeFalse();
	});
});
