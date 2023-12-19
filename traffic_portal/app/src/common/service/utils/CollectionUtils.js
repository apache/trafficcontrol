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

/**
 * CollectionUtils provides methods for manipulating iterable collections of
 * arbitrary items.
 */
class CollectionUtils {
	/**
	 * minimizeArrayDiff reduces the size of an index-sensitive Array diff by
	 * making common elements in the old and new arrays have the same index.
	 *
	 * @template T
	 * @param {T[]} oldItems
	 * @param {T[]} newItems
	 * @returns {(T | undefined)[]}
	 */
	minimizeArrayDiff(oldItems, newItems) {
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

	/**
	 * uniqArray constructs an array from the combination of two arrays such
	 * that the elements are only objects having unique values for their 'key'
	 * property.
	 *
	 * The unique values chosen will be the first found by order of iterating
	 * first `array1` then `array2`.
	 *
	 * @template {PropertyKey} U
	 * @template {Record<U, unknown>} T
	 * @param {T[]} array1
	 * @param {T[]} array2
	 * @param {U} key
	 * @returns {T[]}
	 */
	uniqArray(array1, array2, key) {
		array1 = array1 || [];
		array2 = array2 || [];
		/** @type {Set<unknown>} */
		const keys = new Set();
		/** @type {T[]} */
		const uniq = new Array();
		array1.concat(array2).forEach(function (item) {
			if (!keys.has(item[key])) {
				uniq.push(item);
				keys.add(item[key]);
			}
		});
		return uniq;
	};
}

CollectionUtils.$inject = [];
module.exports = CollectionUtils;
