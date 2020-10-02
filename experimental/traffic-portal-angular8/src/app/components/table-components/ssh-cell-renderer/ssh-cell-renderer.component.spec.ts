import { async, ComponentFixture, TestBed } from "@angular/core/testing";

import { SSHCellRendererComponent } from "./ssh-cell-renderer.component";

describe("SshCellRendererComponent", () => {
	let component: SSHCellRendererComponent;
	let fixture: ComponentFixture<SSHCellRendererComponent>;

	beforeEach(async(() => {
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
