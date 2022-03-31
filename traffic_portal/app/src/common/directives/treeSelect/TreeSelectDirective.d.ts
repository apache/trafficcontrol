/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

export interface TreeData {
	name: string;
	id: string;
	children: TreeData[];
}
export interface RowData {
	label: string;
	value: string;
	depth: number;
	collapsed: boolean;
	hidden: boolean;
	children: RowData[];
}


/**
 * Properties added to an Angular Scope either by directive binding or by being
 * declared in the `link` function.
 */
export interface TreeSelectScopeProperties {
	/** Returns true if the row data should be displayed after filtering. */
 	checkFilters: (row: RowData)=>boolean;
 	/** When collapse icon is clicked on row data. */
	collapse: (row: RowData, evt: Event)=>void;
	/**
	 * Gets the FontAwesome icon class based on if the row data has children and
	 * is collapsed.
	 */
	getClass: (row: RowData)=>string;
	/** Used for form validation, will be assigned to an id attribute. */
	handle: string;
 	initialValue: string;
	/**
	 * Used to properly update the parent on value change, useful for
	 * validation.
	 */
	onUpdate: (output: {value: string})=>void;
	searchText: string;
	/** Updates the selection when clicking a dropdown option. */
	select: (row: RowData)=>void;
	selected: RowData | null;
	shown: boolean;
	/** Toggle the dropdown menu. */
	toggle: ()=>void;
	treeData: Array<TreeData>;
	/** Non-recursed ordered list of rows to display (before filtering). */
	treeRows: Array<RowData>;
}
