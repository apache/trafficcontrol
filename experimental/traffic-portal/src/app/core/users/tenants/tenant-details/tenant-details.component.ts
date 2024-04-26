/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Router } from "@angular/router";
import { RequestTenant, ResponseTenant, Tenant } from "trafficops-types";

import { UserService } from "src/app/api";
import { TreeData } from "src/app/models";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * TenantsDetailsComponent is the controller for the tenant add/edit form.
 */
@Component({
	selector: "tp-tenant-details",
	templateUrl: "./tenant-details.component.html"
})
export class TenantDetailsComponent implements OnInit {
	public new = false;
	public disabled = false;
	public tenant!: Tenant;
	public tenants = new Array<ResponseTenant>();
	public displayTenant: TreeData;

	constructor(
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly userService: UserService,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) {
		this.displayTenant = {
			children: [],
			id: -1,
			name: ""
		};
	}

	/**
	 * Catches when tree-select outputs an update event
	 *
	 * @param evt The TreeData selected
	 */
	public update(evt: TreeData): void {
		const tenant = this.tenants.find(t => t.id === evt.id);
		if (tenant === undefined) {
			this.log.error(`Unknown tenant selected ${evt.id}`);
			return;
		}
		this.tenant.parentId = tenant.id;
	}

	/**
	 * Recursively fills out a nodes children.
	 *
	 * @param tenantByParentId All tenants grouped by parent id.
	 * @param currentTenant The tenant to populate.
	 */
	public breakTenantNode(tenantByParentId: Map<number, Array<TreeData>>, currentTenant: TreeData): void {
		currentTenant.children = (tenantByParentId.get(currentTenant.id) ?? []).map(t => ({...t, children: []} as TreeData));

		currentTenant.children.forEach(t => {
			this.breakTenantNode(tenantByParentId, t);
		});
	}

	/**
	 * Converts the tenants list into the tree-data structure needed by the tree-select component.
	 */
	public constructTreeData(): void {
		const tenantByParentId = new Map<number, Array<TreeData>>();
		this.tenants.forEach(t => {
			if (t.parentId === null) {
				return;
			}
			let children = tenantByParentId.get(t.parentId);
			if(!children) {
				children = [];
			}
			children.push({...t, children: []});
			tenantByParentId.set(t.parentId, children);
		});
		const rootTenant = this.tenants.find(t => t.parentId === null);
		if (rootTenant === undefined) {
			return;
		}
		const rootNode = {...rootTenant, children: []} as TreeData;
		this.breakTenantNode(tenantByParentId, rootNode);

		this.displayTenant = rootNode;
	}

	/**
	 * Angular lifecycle hook.
	 */
	public async ngOnInit(): Promise<void> {
		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			this.log.error("missing required route parameter 'id'");
			return;
		}

		this.tenants = await this.userService.getTenants();
		this.constructTreeData();

		this.new = ID === "new";

		if (this.new) {
			this.setTitle();
			this.new = true;
			this.tenant = {
				active: true,
				name: "",
			} as RequestTenant;
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			this.log.error("route parameter 'id' was non-number:", ID);
			return;
		}
		const tenant = this.tenants.find(t => t.id === numID);
		if (!tenant) {
			this.log.error(`Unable to find tenant with id ${numID}`);
			return;
		}
		this.tenant = tenant;
		this.disabled = this.isRoot();
		this.setTitle();
	}

	/**
	 * Sets the headerTitle based on current Tenant state.
	 *
	 * @private
	 */
	private setTitle(): void {
		const title = this.new ? "New Tenant" : `Tenant: ${this.tenant.name}`;
		this.navSvc.headerTitle.next(title);
	}

	/**
	 * Submits new/changed tenant.
	 *
	 * @param e Html event generated from click
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		if (this.tenant.parentId === undefined) {
			return;
		}
		if (this.new) {
			this.tenant = await this.userService.createTenant(this.tenant as RequestTenant);
			this.new = false;
			await this.router.navigate(["core/tenants", (this.tenant as ResponseTenant).id]);
		} else {
			this.tenant = await this.userService.updateTenant(this.tenant as ResponseTenant);
		}
		this.setTitle();
	}

	/**
	 * Deletes the current tenant.
	 */
	public async deleteTenant(): Promise<void> {
		if (this.new) {
			this.log.error("Unable to delete new tenant");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete tenantn ${this.tenant.name}`,
				title: "Confirm Delete"}
		});
		ref.afterClosed().subscribe(result => {
			if (result) {
				this.userService.deleteTenant((this.tenant as ResponseTenant).id);
				this.router.navigate(["core/tenants"]);
			}
		});
	}

	/**
	 * Determines if the current tenant is the root tenant.
	 *
	 * @returns if a tenant is root
	 */
	public isRoot(): boolean {
		return this.tenant && this.tenant.name === "root";
	}

}
