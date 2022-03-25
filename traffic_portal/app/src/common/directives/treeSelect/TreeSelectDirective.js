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

var TreeSelectDirective = function($document) {
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
        link: function(scope, element, attrs, ngFormController) {
            /**
             * Non-recursed ordered list of rows to display (before filtering)
             * @type TreeSelectDirective.RowData[]
             */
            scope.treeRows = [];
            /** @type string */
            scope.searchText = "";
            /** @type boolean */
            scope.shown = false;
            /** @type TreeSelectDirective.RowData */
            scope.selected = null;

            // Bound variables
            /** @type TreeSelectDirective.TreeData[] */
            scope.treeData;
            /** @type string */
            scope.initialValue;
            /**
             * Used for form validation, will be assigned to an id attribute
             * @type string
             */
            scope.handle;
            /**
             * Used to properly update the parent on value change, useful for validation.
             * @callback onUpdate
             * @param {string} value
             */
            scope.onUpdate;


            element.bind("click",
                /**
                 * Handles click events within this element, prevents $document propagation as well as binds
                 * the close event to it.
                 * @param evt
                 */
                function(evt) {
                    if(scope.shown) {
                        evt.stopPropagation();
                        $document.on("click", closeSelect);
                    }
            });
            /**
             * Detects when a click is trigger outside this element. Used to close dropdown and ensure
             * there is no $document event pollution.
             * @param evt
             */
            const closeSelect = function(evt){
                scope.close();
                $document.off("click", closeSelect);
                scope.$apply();
            };

            /**
             * Initializes treeRows, must be called whenever treeData changes
             */
            const reInit = function() {
                scope.treeRows = [];
                scope.selected = null;
                for(let data of scope.treeData) {
                    if (data != undefined)
                        addNode(data, 0);
                }
            }

            /**
             * Converts a tree data node into row data recursively.
             * @param {TreeSelectDirective.TreeData} row
             * @param {number} depth
             * @returns TreeSelectDirective.RowData
             */
            const addNode = function(row, depth) {
                scope.treeRows.push({
                    label: row.name,
                    value: row.id,
                    depth: depth,
                    children: [],
                    collapsed: false,
                    hidden: false
                });
                const last = scope.treeRows.length - 1;
                if(row.id === scope.initialValue || row.name === scope.initialValue) {
                    scope.selected = scope.treeRows[last];
                }
                if(row.children != null) {
                    for(const child of row.children) {
                        if(child === undefined) continue;
                        scope.treeRows[last].children.push(addNode(child, depth + 1));
                    }
                }
                return scope.treeRows[last];
            }

            /**
             * Collapses a row data and recursively hides and collapses its children.
             * @param {TreeSelectDirective.RowData} row
             * @param {boolean?} state
             */
            const collapseRecurse = function(row, state) {
                if(row.children.length === 0) return;
                for(const treeRow of scope.treeRows) {
                    if (treeRow.value === row.value) {
                        if(state === undefined)
                            treeRow.collapsed = !treeRow.collapsed;
                        else
                            treeRow.collapsed = state;
                        for(let treeChild of treeRow.children) {
                            treeChild.hidden = treeRow.collapsed;
                            collapseRecurse(treeChild, treeRow.collapsed);
                        }
                    }
                }
            }

            /**
             * Returns true if the inputs letters are also present in text with the same order
             * @param {string} text
             * @param {string} input
             * @returns {boolean}
             */
            const fuzzyMatch = function(text, input) {
                if(input === "") return true;
                if(text === undefined) return false;
                text = text.toString().toLowerCase();
                input = input.toString().toLowerCase();
                let n = -1;
                for(let i in input) {
                    const letter = input[i];
                    if (!~(n = text.indexOf(letter, n + 1))) return false;
                }
                return true;
            }

            /**
             * Triggers onUpdate binding
             */
            scope.update = function() {
                scope.onUpdate({value: scope.selected.value ?? ""});
            }

            /**
             * Toggle the dropdown menu
             */
            scope.toggle = function() {
                scope.shown = !scope.shown;
            }
            /**
             * Close the dropdown menu.
             */
            scope.close = function() {
                scope.shown = false;
            }


            /**
             * Updates the selection when clicking a dropdown option
             * @param {TreeSelectDirective.RowData} row
             */
            scope.select = function(row) {
                scope.selected = row;
                scope.selection = row.value;
                scope.update();

                if(row.value !== this.initialValue) {
                    ngFormController[this.handle].$setDirty();
                } else {
                    ngFormController[this.handle].$setPristine();
                }
                scope.close();
            }
            /**
             * When collapse icon is clicked on row data
             * @param {TreeSelectDirective.RowData} row
             * @param evt
             */
            scope.collapse = function(row, evt) {
                evt.stopPropagation();
                return collapseRecurse(row);
            }

            /**
             * Returns true if the row data
             * @param {TreeSelectDirective.RowData} testRow
             * @returns {boolean}
             */
            scope.checkFilters = function(testRow) {
                if(testRow.hidden && scope.searchText.length === 0)
                    return false;
                return fuzzyMatch(testRow.label, scope.searchText);
            }

            /**
             * Gets the FontAwesome icon class based on if the row data has children and is collapsed
             * @param {TreeSelectDirective.RowData} row
             * @returns {string}
             */
            scope.getClass = function(row) {
                if(row.collapsed && row.children.length > 0) return "fa-minus";
                else if(row.children.length > 0) return "fa-plus";
                else return "fa-users";
            }

            scope.$watch("treeData", function(newVal, oldVal) {
                reInit();
            });
        }
    }
};

TreeSelectDirective.$inject = ['$document'];
module.exports = TreeSelectDirective;
