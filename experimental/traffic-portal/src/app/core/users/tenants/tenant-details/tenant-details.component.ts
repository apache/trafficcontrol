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
import { Location } from "@angular/common";
import { Component, OnInit } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { RequestTenant, ResponseTenant, Tenant } from "trafficops-types";

import { UserService } from "src/app/api";
import { TreeData } from "src/app/models";

/**
 * TenantsDetailsComponent is the controller for the tenant add/edit form.
 */
@Component({
	selector: "tp-tenant-details",
	styleUrls: ["./tenant-details.component.scss"],
	templateUrl: "./tenant-details.component.html"
})
export class TenantDetailsComponent implements OnInit {
	public new = false;
	public disabled = false;
	public tenant!: Tenant;
	public tenants = new Array<ResponseTenant>();
	public displayTenant: TreeData;

	constructor(private readonly route: ActivatedRoute, private readonly userService: UserService,
		private readonly location: Location) {
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
			console.error(`Unknown tenant selected ${evt.id}`);
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
			console.error("missing required route parameter 'id'");
			return;
		}

		this.tenants = await this.userService.getTenants();
		this.constructTreeData();

		if (ID === "new") {
			this.new = true;
			this.tenant = {
				active: true,
				name: "",
			} as RequestTenant;
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			console.error("route parameter 'id' was non-number:", ID);
			return;
		}
		const tenant = this.tenants.find(t => t.id === numID);
		if (!tenant) {
			console.error(`Unable to find tenant with id ${numID}`);
			return;
		}
		this.tenant = tenant;
		this.disabled = this.isRoot();
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
		} else {
			this.tenant = await this.userService.updateTenant(this.tenant as ResponseTenant);
		}
	}

	/**
	 * Deletes the current tenant.
	 */
	public async deleteTenant(): Promise<void> {
		if (this.new) {
			console.error("Unable to delete new tenant");
			return;
		}
		await this.userService.deleteTenant((this.tenant as ResponseTenant).id);
		this.location.back();
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
