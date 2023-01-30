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
import { ComponentFixture, TestBed } from "@angular/core/testing";

import { TreeData } from "src/app/models";

import { TreeSelectComponent } from "./tree-select.component";

const treeData: Array<TreeData> = [{
	children: [{
		children: [],
		id: 11,
		name: "n11"
	}],
	id: 1,
	name: "n1",
}, {
	children: [{
		children: [{
			children: [],
			id: 211,
			name: "n211"
		}],
		id: 21,
		name: "n21"
	}],
	id: 2,
	name: "n2"
}];

/**
 * Returns all nodes in the component.
 *
 * @param component The component to get nodes from
 * @returns Every node in the tree
 */
function allNodes(component: TreeSelectComponent): Array<TreeData> {
	const ret = new Array<TreeData>();
	component.treeData.forEach(root => {
		ret.push(root);
		component.treeControl.getDescendants(root).forEach(node => {
			ret.push(node);
		});
	});
	return ret;
}

/**
 * Returns all nodes in the component that are visible (or unset).
 *
 * @param component Component to get nodes from
 * @returns All visible nodes
 */
function visibleData(component: TreeSelectComponent): Array<TreeData> {
	return allNodes(component).filter(node => node.visible === undefined || node.visible);
}

describe("TreeSelectComponent", () => {
	let component: TreeSelectComponent;
	let fixture: ComponentFixture<TreeSelectComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [TreeSelectComponent]
		})
			.compileComponents();

		fixture = TestBed.createComponent(TreeSelectComponent);
		component = fixture.componentInstance;
		component.treeData = treeData;
		fixture.detectChanges();
	});

	it("should create", async () => {
		expect(component).toBeTruthy();
		component.filter.next("");

		expect(component.selected.id).toBe(-1);
		expect(visibleData(component).length).toBe(5);
	});

	it("should filter", async () => {
		component.filter.next("n211");
		expect(visibleData(component).length).toBe(1);
		component.filter.next("");
		expect(visibleData(component).length).toBe(5);
		component.filter.next("n21");
		expect(visibleData(component).length).toBe(2);
	});

	it("should select initial value", () => {
		fixture = TestBed.createComponent(TreeSelectComponent);
		component = fixture.componentInstance;
		component.treeData = treeData;
		component.initialValue = treeData[1].id;
		fixture.detectChanges();

		expect(component.selected.id).toBe(treeData[1].id);

		fixture = TestBed.createComponent(TreeSelectComponent);
		component = fixture.componentInstance;
		component.treeData = treeData;
		component.initialValue = treeData[0].children[0].id;
		fixture.detectChanges();

		expect(component.selected.id).toBe(treeData[0].children[0].id);
	});
});
