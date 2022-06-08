import { Component, type OnInit } from "@angular/core";
import type { MatSelectChange } from "@angular/material/select";
import { ActivatedRoute } from "@angular/router";

import { UserService } from "src/app/api";
import type { Role, Tenant, User } from "src/app/models";

/**
 * UserDetailsComponent is the controller for the page for viewing/editing a
 * user.
 */
@Component({
	selector: "tp-user-details",
	styleUrls: ["./user-details.component.scss"],
	templateUrl: "./user-details.component.html"
})
export class UserDetailsComponent implements OnInit {

	public user!: User;
	public roles = new Array<Role>();
	public tenants = new Array<Tenant>();

	constructor(private readonly userService: UserService, private readonly route: ActivatedRoute) {
		this.userService.getRoles().then(rs=>this.roles=rs);
		this.userService.getTenants().then(ts=>this.tenants=ts);
	}

	public async ngOnInit(): Promise<void> {
		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			console.error("missing required route parameter 'id'");
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			console.error("route parameter 'id' was non-number:", ID);
		}
		this.user = await this.userService.getUsers(numID);
		console.log(this.user);
	}

	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		this.user = await this.userService.updateUser(this.user);
	}

	public role(): Role {
		const role = this.roles.find(r=>r.id === this.user.role);
		if (!role) {
			throw new Error(`user's Role "${this.user.rolename}" (#${this.user.role}) does not exist`);
		}
		return role;
	}

	public tenant(): Tenant {
		const tenant = this.tenants.find(t=>t.id === this.user.tenantId);
		if (!tenant) {
			throw new Error(`user's Tenant "${this.user.tenant}" (#${this.user.tenantId}) does not exist`);
		}
		return tenant;
	}

	public updateRole(r: MatSelectChange & {value: Role}): void {
		this.user.role = r.value.id;
		this.user.rolename = r.value.name;
	}

	public updateTenant(r: MatSelectChange & {value: Tenant}): void {
		this.user.tenantId = r.value.id;
		this.user.tenant = r.value.name;
	}
}
