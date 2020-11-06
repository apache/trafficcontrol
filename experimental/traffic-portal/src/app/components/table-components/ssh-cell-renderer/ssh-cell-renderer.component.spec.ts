import { waitForAsync, ComponentFixture, TestBed } from "@angular/core/testing";

import { SSHCellRendererComponent } from "./ssh-cell-renderer.component";

describe("SshCellRendererComponent", () => {
	let component: SSHCellRendererComponent;
	let fixture: ComponentFixture<SSHCellRendererComponent>;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [ SSHCellRendererComponent ]
		})
		.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(SSHCellRendererComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
