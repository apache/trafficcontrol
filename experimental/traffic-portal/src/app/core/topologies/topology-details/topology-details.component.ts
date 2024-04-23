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

import { NestedTreeControl } from "@angular/cdk/tree";
import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { MatTreeNestedDataSource } from "@angular/material/tree";
import { ActivatedRoute, Router } from "@angular/router";
import { ResponseTopology } from "trafficops-types";

import { TopologyService, TopTreeNode } from "src/app/api";
import {
	DecisionDialogComponent,
	DecisionDialogData,
} from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * TopologyDetailComponent is the controller for a Topology's "detail" page.
 */
@Component({
	selector: "tp-topology-details",
	styleUrls: ["./topology-details.component.scss"],
	templateUrl: "./topology-details.component.html",
})
export class TopologyDetailsComponent implements OnInit {
	public new = false;

	/** Loader status for the actions */
	public loading = true;

	public oldName: string | undefined = undefined;

	public topology: ResponseTopology = {
		description: "",
		lastUpdated: new Date(),
		name: "",
		nodes: [],
	};
	public showErrors = false;
	public topologies: Array<ResponseTopology> = [];
	public topologySource = new MatTreeNestedDataSource<TopTreeNode>();
	public topologyControl = new NestedTreeControl<TopTreeNode>(node => node.children);

	constructor(
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly api: TopologyService,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) { }

	/**
	 * Angular lifecycle hook where data is initialized.
	 */
	public async ngOnInit(): Promise<void> {
		const name = this.route.snapshot.paramMap.get("name");

		const topologiesPromise = this.api.getTopologies().then(topologies => this.topologies = topologies);
		if (name === null) {
			this.new = true;
			this.setTitle();
			await topologiesPromise;
			this.loading = false;
			return;
		}
		await topologiesPromise;
		const index = this.topologies.findIndex(c => c.name === name);
		if (index < 0) {
			this.log.error(`no such Topology: ${name}`);
			this.loading = false;
			return;
		}
		this.oldName = name;
		this.topology = this.topologies.splice(index, 1)[0];
		this.topologySource.data = TopologyService.topologyToTree(this.topology);
		this.loading = false;
	}

	/**
	 * Used by angular to determine if this node should be a nested tree node.
	 *
	 * @param _ Index of the current node.
	 * @param node Node to test.
	 * @returns If the node has children.
	 */
	public hasChild(_: number, node: TopTreeNode): boolean {
		return Array.isArray(node.children) && node.children.length > 0;
	}

	/**
	 * Sets the title of the page to either "new" or the name of the displayed
	 * Topology, depending on the value of TopologyDetailComponent.new.
	 */
	private setTitle(): void {
		const title = this.new ? "New Topology" : `Topology: ${this.topology.name}`;
		this.navSvc.headerTitle.next(title);
	}

	/**
	 * Deletes the Topology.
	 */
	public async delete(): Promise<void> {
		if (this.new) {
			this.log.error("Unable to delete new Topology");
			return;
		}
		const ref = this.dialog.open<DecisionDialogComponent, DecisionDialogData, boolean>(
			DecisionDialogComponent,
			{
				data: {
					message: `Are you sure you want to delete Topology ${this.topology.name}?`,
					title: "Confirm Delete"
				}
			}
		);
		ref.afterClosed().subscribe(result => {
			if (result) {
				this.api.deleteTopology(this.topology);
				this.router.navigate(["core/topologies"]);
			}
		});
	}

	/**
	 * Submits new/updated Topology.
	 *
	 * @param e HTML form submission event.
	 */
	public async submit(e: Event): Promise<void> {
		this.topology = TopologyService.treeToTopology(this.topology.name, this.topology.description, this.topologySource.data);

		e.preventDefault();
		e.stopPropagation();
		this.showErrors = true;
		if (this.new) {
			this.topology = await this.api.createTopology(this.topology);
			this.new = false;
			await this.router.navigate(["core/topologies", this.topology.name]);
		} else {
			this.topology = await this.api.updateTopology(this.topology, this.oldName);
		}
		this.setTitle();
	}
}
