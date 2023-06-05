import { HarnessLoader } from "@angular/cdk/testing";
import { TestbedHarnessEnvironment } from "@angular/cdk/testing/testbed";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatFormFieldHarness } from "@angular/material/form-field/testing";

import { CertAuthorComponent } from "./cert-author.component";
import { AppUIModule } from "src/app/app.ui.module";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";

describe("CertAuthorComponent", () => {
	let component: CertAuthorComponent;
	let fixture: ComponentFixture<CertAuthorComponent>;
	let loader: HarnessLoader;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ CertAuthorComponent ],
			imports: [ AppUIModule, NoopAnimationsModule ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(CertAuthorComponent);
		component = fixture.componentInstance;
		component.author = {
			commonName: "name"
		};
		fixture.detectChanges();
		loader = TestbedHarnessEnvironment.documentRootLoader(fixture);
	});

	it("should create", async () => {
		expect(component).toBeTruthy();
		await fixture.whenRenderingDone();
		const formFields = await loader.getAllHarnesses(MatFormFieldHarness);
		expect(formFields.length).toBe(6);
	});
});
