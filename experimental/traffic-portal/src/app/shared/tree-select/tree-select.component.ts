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
import { Component, EventEmitter, HostListener, Input, OnInit, Output } from "@angular/core";
import { MatTreeNestedDataSource } from "@angular/material/tree";
import { Subject } from "rxjs";

import { TreeData } from "src/app/models";
import { fuzzyScore } from "src/app/utils";

/**
 * TreeSelectComponent is the controller for a tree select input
 */
@Component({
	selector: "tp-tree-select",
	styleUrls: ["./tree-select.component.scss"],
	templateUrl: "./tree-select.component.html"
})
export class TreeSelectComponent implements OnInit {
	@Input() public treeData = new Array<TreeData>();
	@Input() public label = "";
	@Input() public handle = "tree-select";
	@Input() public disabled = false;
	@Input() public initialValue = -1;
	@Output() public nodeSelected = new EventEmitter<TreeData>();
	public shown = false;
	public dataSource = new MatTreeNestedDataSource<TreeData>();
	public treeControl = new NestedTreeControl<TreeData>(node => node.children);
	public selected: TreeData = {children: [], id: -1, name: ""};
	public filter = new Subject<string>();

	/**
	 * Used by angular to determine if this node should be a nested tree node.
	 *
	 * @param _ Index of the current node.
	 * @param node Node to test.
	 * @returns If the node has children.
	 */
	public hasChild(_: number, node: TreeData): boolean {
		return node.children !== undefined && node.children.length > 0;
	}

	/**
	 * Used by angular to determine a node's visible property
	 *
	 * @param node The node to test.
	 * @returns Visible value, unset means visible.
	 */
	public isVisible(node: TreeData): boolean {
		return node?.visible ?? true;
	}

	/**
	 * Used by angular when the search input is changed
	 *
	 * @param $event The html input event.
	 */
	public filterChanged($event: Event): void {
		this.filter.next(($event.target as HTMLInputElement).value);
	}

	/**
	 * Listens for clicks outside this component to close the drop down.
	 */
	@HostListener("document:click", ["$event"])
	public documentClick(): void {
		if (this.shown) {
			this.shown = false;
		}
	}

	/**
	 * Called when a tree node is selected.
	 *
	 * @param node The selected node.
	 */
	public select(node: TreeData): void {
		this.shown = false;
		this.selected = node;
		this.nodeSelected.emit(node);
	}

	/**
	 * Called to toggle if the tree select drop down is visible.
	 *
	 * @param evt DOM event
	 */
	public toggle(evt: Event): void {
		evt.stopPropagation();
		evt.preventDefault();
		this.shown = !this.shown;
	}

	/**
	 * Angular lifecycle hook.
	 */
	public ngOnInit(): void {
		this.dataSource.data = this.treeData;

		for (const data of this.treeData) {
			if (data.id === this.initialValue) {
				this.selected = data;
				break;
			}
			const res = this.treeControl.getDescendants(data).find(desc => desc.id === this.initialValue);
			if (res !== undefined) {
				this.selected = res;
				break;
			}
		}

		this.filter.subscribe(value => {
			this.treeData.forEach(node => {
				if(value === "") {
					node.visible = true;
					node.containerNeeded = false;
					this.treeControl.getDescendants(node).forEach(desc => {
						desc.visible = true;
						desc.containerNeeded = false;
					});
				} else {
					this.filterNode(node, value);
				}
			});
		});
	}

	/**
	 * Recursively fuzzy filter a node on its name.
	 *
	 * @param node The node to filter.
	 * @param value The filter value.
	 * @returns If the node passes the filter.
	 */
	public filterNode(node: TreeData, value: string): boolean {
		let score: number;
		if(value === "") {
			score = 0;
		} else {
			score = fuzzyScore(node.name.toLocaleLowerCase(), value.toLocaleLowerCase());
		}
		node.visible = (score !== Infinity);
		if(node.containerNeeded) {
			node.containerNeeded = false;
		}
		this.treeControl.getDescendants(node).forEach(desc => node.containerNeeded = node.containerNeeded || this.filterNode(desc, value));
		return node.visible || (node.containerNeeded ?? false);
	}

}
