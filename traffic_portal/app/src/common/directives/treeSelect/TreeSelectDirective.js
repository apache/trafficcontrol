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

/** @typedef {import("angular").IController} Controller */
/** @typedef {import("angular").IDirectiveLinkFn} IDirectiveLinkFn */
/** @typedef {import("angular").IDocumentService} NGDocument */
/** @typedef {import("angular").IRootElementService} NGElement */
/** @typedef {import("angular").IScope} Scope */

/** @typedef {import("./TreeSelectDirective").RowData} RowData */
/** @typedef {import("./TreeSelectDirective").TreeData} TreeData */
/** @typedef {import("./TreeSelectDirective").TreeSelectScopeProperties} TreeSelectScopeProperties */

/**
 *
 * @param {NGDocument} $document
 * @returns
 */
function TreeSelectDirective($document) {
	return {
		restrict: "E",
		require: "^form",
		templateUrl: "common/directives/treeSelect/tree.select.tpl.html",
		replace: true,
		scope: {
			treeData: '=',
			initialValue: '<',
			handle: '@',
			onUpdate: "&"
		},
		/**
		 * Controller logic for the Directive.
		 *
		 * @param {Scope & TreeSelectScopeProperties} scope
		 * @param {NGElement} element Host element
		 * @param {unknown} _ Host attribute bindings (unused)
		 * @param {Controller} ngFormController
		 */
		link: function(scope, element, _, ngFormController) {

			/**
			 * Close the dropdown menu.
			 */
			function close() {
				scope.shown = false;
			}

			/**
			 * Detects when a click is trigger outside this element. Used to close dropdown and ensure
			 * there is no $document event pollution.
			 */
			function closeSelect() {
				close();
				$document.off("click", closeSelect);
				scope.$apply();
			};

			/**
			 * Converts a tree data node into row data recursively.
			 *
			 * @param {TreeData} row
			 * @param {number} depth
			 * @returns {RowData}
			 */
			function addNode(row, depth) {
				/** @type {RowData} */
				const node = {
					label: row.name,
					value: row.id,
					depth: depth,
					children: [],
					collapsed: false,
					hidden: false
				};
				scope.treeRows.push(node);
				if(row.id === scope.initialValue || row.name === scope.initialValue) {
					scope.selected = node
				}
				if(row.children !== null && row.children !== undefined) {
					for(const child of row.children) {
						if(child === undefined) continue;
						node.children.push(addNode(child, depth + 1));
					}
				}
				return node;
			}

			/**
			 * Collapses a row data and recursively hides and collapses its children.
			 *
			 * @param {RowData} row
			 * @param {boolean} [state]
			 */
			 function collapseRecurse(row, state) {
				if(row.children.length === 0) return;
				for(const treeRow of scope.treeRows) {
					if (treeRow.value === row.value) {
						if(state === undefined)
							treeRow.collapsed = !treeRow.collapsed;
						else
							treeRow.collapsed = state;
						for(const treeChild of treeRow.children) {
							treeChild.hidden = treeRow.collapsed;
							collapseRecurse(treeChild, treeRow.collapsed);
						}
					}
				}
			}

			/**
			 * Initializes treeRows, must be called whenever treeData changes
			 */
			 function reInit() {
				scope.treeRows = [];
				scope.selected = null;
				for(const data of scope.treeData) {
					if (data !== undefined)
						addNode(data, 0);
				}
			}

			/**
			 * Returns true if the inputs letters are also present in text with the same order.
			 *
			 * @param {string} text
			 * @param {string} input
			 * @returns {boolean}
			 */
			function fuzzyMatch(text, input) {
				if(input === "") return true;
				if(text === undefined) return false;
				text = text.toString().toLowerCase();
				input = input.toString().toLowerCase();
				let n = -1;
				for(const letter of input) {
					if (!~(n = text.indexOf(letter, n + 1))) return false;
				}
				return true;
			}

			/**
			 * Triggers onUpdate binding
			 */
			function update() {
				scope.onUpdate({value: scope.selected?.value ?? ""});
			}

			scope.treeRows = [];
			scope.searchText = "";
			scope.shown = false;
			scope.selected = null;

			element.bind("click",
				evt => {
					if(scope.shown) {
						evt.stopPropagation();
						$document.on("click", closeSelect);
					}
				}
			);

			scope.toggle = () => scope.shown = !scope.shown;


			scope.select = row => {
				scope.selected = row;
				update();

				if(row.value !== this.initialValue) {
					ngFormController[scope.handle].$setDirty();
				}
				close();
			}

			scope.collapse = (row, evt) => {
				evt.stopPropagation();
				collapseRecurse(row);
			}

			scope.checkFilters = testRow => {
				if(testRow.hidden && scope.searchText.length === 0)
					return false;
				return fuzzyMatch(testRow.label, scope.searchText);
			}

			scope.getClass = function(row) {
				if(row.collapsed && row.children.length > 0) return "fa-minus";
				else if(row.children.length > 0) return "fa-plus";
				else return "fa-users";
			}

			scope.$watch("treeData", reInit);

			ngFormController.$removeControl(ngFormController[scope.handle + "searchText"]);
		}
	}
};

TreeSelectDirective.$inject = ['$document'];
module.exports = TreeSelectDirective;
