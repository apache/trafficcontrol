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
import { animate, state, style, transition, trigger } from "@angular/animations";
import { NestedTreeControl } from "@angular/cdk/tree";
import { AfterViewInit, Component, OnInit, ViewChild } from "@angular/core";
import { MatSidenav } from "@angular/material/sidenav";
import { MatTreeNestedDataSource } from "@angular/material/tree";
import { Router, RouterEvent, Event, NavigationEnd, IsActiveMatchOptions } from "@angular/router";
import { filter } from "rxjs/operators";

import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService, TreeNavNode } from "src/app/shared/navigation/navigation.service";

/**
 * TpSidebarComponent is the controller for the sidebar.
 */
@Component({
	animations: [
		trigger("slide", [
			state(
				"hide",
				style({
					height: 0,
					visibility: "hidden"
				})
			),
			state(
				"show",
				style({
					height: "*"
				})
			),
			transition("hide <=> show", [
				animate("125ms linear")
			])
		])
	],
	selector: "tp-sidebar",
	styleUrls: ["./tp-sidebar.component.scss"],
	templateUrl: "./tp-sidebar.component.html",
})
export class TpSidebarComponent implements OnInit, AfterViewInit {
	public dataSource = new MatTreeNestedDataSource<TreeNavNode>();
	public treeCtrl = new NestedTreeControl<TreeNavNode>(node => node.children);

	private lastRoute = "";
	public readonly routeOptions: IsActiveMatchOptions = {
		fragment: "exact",
		matrixParams: "ignored",
		paths: "exact",
		queryParams: "ignored"
	};

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

	/**
	 * Determines if the node is a root node (no parents)
	 *
	 * @param node The node to check
	 * @returns If the node is root
	 */
	public isRoot(node: TreeNavNode): boolean {
		return !this.childToParent.has(this.nodeHandle(node));
	}

	constructor(
		private readonly navService: NavigationService,
		private readonly route: Router,
		public readonly user: CurrentUserService,
		private readonly log: LoggingService,
	) { }

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
	 * Not done in OnInit because we are dependent on the sidenav
	 */
	public ngAfterViewInit(): void {
		this.navService.sidebarHidden.subscribe(hidden => {
			if(hidden && this.sidenav.opened) {
				this.sidenav.close().catch(err => {
					this.log.error(`Unable to close sidebar: ${err}`);
				});
			} else if (!hidden && !this.sidenav.opened) {
				this.sidenav.open().catch(err => {
					this.log.error(`Unable to open sidebar: ${err}`);
				});
			}
		});
	}

	/**
	 * Angular lifecycle hook.
	 */
	public ngOnInit(): void {
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
							// Prevent infinite loops
							let depth = 0;
							while(parent !== undefined) {
								this.treeCtrl.expand(parent);
								parent = this.childToParent.get(this.nodeHandle(parent));
								if(depth++ > 5) {
									this.log.error(`Maximum depth ${depth} reached, aborting expand on ${parent?.name ?? "unknown"}`);
									break;
								}
							}
							return;
						}
					}
				}
			}
		});
	}

	/**
	 * Collapse all tree entries
	 *
	 * @private
	 */
	public collapseAll(): void {
		this.treeCtrl.collapseAll();
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

	/**
	 * Determines whether href is absolute or not - for TPv2 to TPv1 redirect
	 *
	 * @param href route string
	 * @returns true if href absolute
	 */
	public isAbsoluteURL(href: string): boolean {
		const regexPattern = /^(?:[a-z]+:)?\/\//i;
		return regexPattern.test(href);
	}
}
