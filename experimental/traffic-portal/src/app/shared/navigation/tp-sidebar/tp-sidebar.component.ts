import { NestedTreeControl } from "@angular/cdk/tree";
import { Component, OnInit, ViewChild } from "@angular/core";
import { MatSidenav } from "@angular/material/sidenav";
import { MatTreeNestedDataSource } from "@angular/material/tree";

import { NavigationService, TreeNavNode } from "src/app/shared/navigation/navigation.service";

@Component({
	selector: "tp-sidebar",
	styleUrls: ["./tp-sidebar.component.scss"],
	templateUrl: "./tp-sidebar.component.html",
})
export class TpSidebarComponent implements OnInit {
	public dataSource = new MatTreeNestedDataSource<TreeNavNode>();
	public treeCtrl = new NestedTreeControl<TreeNavNode>(node => node.children);

	@ViewChild("sidenav") public sidenav!: MatSidenav;

	/**
	 * Used by angular to determine if this node should be a nested tree node.
	 *
	 * @param _ Index of the current node.
	 * @param node Node to test.
	 * @returns If the node has children.
	 */
	public hasChild(_: number, node: TreeNavNode): boolean {
		return node.children !== undefined && node.children.length > 0;
	}

	constructor(private readonly navService: NavigationService) {
	}

	public ngOnInit(): void {
		this.navService.sidebarHidden.subscribe(hidden => {
			if(hidden) {
				this.sidenav.close();
			} else {
				this.sidenav.open();
			}
		});
		this.navService.sidebarNavs.subscribe(navs => {
			this.dataSource.data = navs;
		});
	}
}
