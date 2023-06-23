/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import { NestedTreeControl } from "@angular/cdk/tree";
import { Component, OnInit, ViewChild } from "@angular/core";
import { MatSidenav } from "@angular/material/sidenav";
import { MatTreeNestedDataSource } from "@angular/material/tree";
import { Router, RouterEvent, Event, NavigationEnd } from "@angular/router";
import { filter } from "rxjs/operators";

import { NavigationService, TreeNavNode } from "src/app/shared/navigation/navigation.service";

/**
 * TpSidebarComponent is the controller for the sidebar.
 */
@Component({
	selector: "tp-sidebar",
	styleUrls: ["./tp-sidebar.component.scss"],
	templateUrl: "./tp-sidebar.component.html",
})
export class TpSidebarComponent implements OnInit {
	public dataSource = new MatTreeNestedDataSource<TreeNavNode>();
	public treeCtrl = new NestedTreeControl<TreeNavNode>(node => node.children);

	public hidden = false;
	private lastRoute = "";

	/**
	 * Used in the sidebar to ensure the active page is visible.1
	 *
	 * @private
	 */
	private childToParent = new Map<string, TreeNavNode>();

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

	constructor(private readonly navService: NavigationService, private readonly route: Router) {
	}

	/**
	 * Adds to childToParent from the given node.
	 *
	 * @param node The node to map.
	 * @private
	 */
	private mapChild(node: TreeNavNode): void {
		if(node.children !== undefined) {
			for(const child of node.children) {
				this.childToParent.set(this.nodeHandle(child), node);
				this.mapChild(child);
			}
		}
	}

	/**
	 * Angular lifecycle hook.
	 */
	public ngOnInit(): void {
		this.navService.sidebarHidden.subscribe(hidden => {
			if(this.sidenav) {
				this.hidden = hidden;
				if(hidden && this.sidenav.opened) {
					this.sidenav.close().catch(err => {
						console.error(`Unable to close sidebar: ${err}`);
					});
				} else if (!this.sidenav.opened) {
					this.sidenav.open().catch(err => {
						console.error(`Unable to open sidebar: ${err}`);
					});
				}
			}
		});
		this.navService.sidebarNavs.subscribe(navs => {
			this.dataSource.data = navs;

			this.childToParent = new Map<string, TreeNavNode>();
			navs.forEach(nav => this.mapChild(nav));
		});

		this.route.events.pipe(
			filter((e: Event): e is NavigationEnd => e instanceof RouterEvent)
		).subscribe((e: NavigationEnd) => {
			const path = e.url.split("?")[0];
			if(path !== this.lastRoute) {
				this.lastRoute = path;
				for(const node of this.dataSource.data) {
					for(const child of this.treeCtrl.getDescendants(node)) {
						if(child.href === path) {
							this.treeCtrl.expand(node);
							let parent = this.childToParent.get(this.nodeHandle(child));
							let depth = 0;
							while(parent !== undefined && depth++ < 5) {
								this.treeCtrl.expand(parent);
								parent = this.childToParent.get(this.nodeHandle(parent));
							}
							return;
						}
					}
				}
			}
		});
	}

	/**
	 * Gets the key used in the parent map.
	 *
	 * @param node Node to get the key from.
	 * @returns node key
	 */
	public nodeHandle(node: TreeNavNode): string {
		return `${node.name}${node.href ?? ""}`;
	}
}
