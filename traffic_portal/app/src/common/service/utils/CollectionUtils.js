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

var CollectionUtils = function() {

	// minimizeArrayDiff reduces the size of an index-sensitive Array diff by
	// making common elements in the old and new arrays have the same index.
	this.minimizeArrayDiff = function (oldItems, newItems) {
		const minimalDiffItems = [];
		const addedItems = [];

		let newItemsIndex, newItem;
		const oldItemsIterator = oldItems.entries();
		const newItemsIterator = newItems.entries();
		for (let oldItemsNext = oldItemsIterator.next(), newItemsNext = newItemsIterator.next(); !(oldItemsNext.done || newItemsNext.done);) {
			const [, oldItem] = oldItemsNext.value;
			[newItemsIndex, newItem] = newItemsNext.value;
			if (oldItem < newItem) {
				newItemsIndex--;
				minimalDiffItems.push(undefined);
				oldItemsNext = oldItemsIterator.next();
				continue;
			} else if (oldItem > newItem) {
				addedItems.push(newItem);
				newItemsNext = newItemsIterator.next();
				continue;
			}
			minimalDiffItems.push(newItem);
			oldItemsNext = oldItemsIterator.next();
			newItemsNext = newItemsIterator.next();
		}
		minimalDiffItems.push(...addedItems);
		if (newItemsIndex !== undefined && newItemsIndex < newItems.length - 1) {
			minimalDiffItems.push(...newItems.slice(newItemsIndex + 1));
		}
		return minimalDiffItems;
	};

	this.uniqArray = function(array1, array2, key) {
		array1 = array1 || [];
		array2 = array2 || [];

		const keys = new Set();
		const uniq = new Array();
		array1.concat(array2).forEach(function(item) {
			if (!keys.has(item[key])) {
				uniq.push(item);
				keys.add(item[key]);
			}
		});
		return uniq;
	};

};

CollectionUtils.$inject = [];
module.exports = CollectionUtils;
