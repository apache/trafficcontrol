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
 * @typedef Tenant
 * @property {string} name
 * @property {number} id
 * @property {number} parentId
 */

/**
 * @typedef HierarchyTenant
 * @property {string} [name]
 * @property {number} [id]
 * @property {number} [parentId]
 * @property {HierarchyTenant[]} children
 */

/**
 * @typedef HierarchicalStructure
 * @property {number} [level]
 * @property {HierarchicalStructure[]} children
 */

/**
 * TenantUtils provides methods for dealing with Tenant hierarchies.
 */
class TenantUtils {
	/**
	 * Recurse through a hierarchical structure and add a 'level' property to
	 * each item.
	 *
	 * @private
	 * @param {HierarchicalStructure[]} arr
	 * @param {number} level
	 */
	applyLevels(arr, level) {
		for (const item of arr) {
			item.level = level;
			this.applyLevels(item.children, level + 1);
		}
	}

	/**
	 * @private
	 * @param {{name: string}} a
	 * @param {{name: string}} b
	 */
	hierarchySortFunc(a, b) {
		return a.name > b.name ? 1 : 0;
	}

	/**
	 * Sort parent groups into an array representative of the hierarchy order.
	 *
	 * @param {Record<number, {id: number; name: string}[]>} parentGroups
	 * @param {number} rootParentId
	 * @param {{id: number; name: string}[]} result
	 */
	hierarchySort(parentGroups, rootParentId, result) {

		if (parentGroups[rootParentId] == undefined)
			return;
		const arr = parentGroups[rootParentId].sort(this.hierarchySortFunc);
		for (const item of arr) {
			result.push(item);
			this.hierarchySort(parentGroups, item.id, result);
		}

		return result;
	}

	/**
	 * Take a flat tenant array and group the tenants by parentId.
	 *
	 * @template {{id: number; name: string; parentId: number}} T
	 * @param {T[]} tenants
	 * @returns {Record<number, T[]>}
	 */
	groupTenantsByParent(tenants) {

		/** @type {Record<number, T[]>} */
		const parentGroups = {};

		for (const tenant of tenants) {
			if (!(tenant.parentId in parentGroups))
				parentGroups[tenant.parentId] = [];
			parentGroups[tenant.parentId].push(tenant);
		}

		return parentGroups;
	};

	/**
	 * Takes a flat tenant array and turns it into a hierarchical
	 * representation of tenants based on their parent/child relationships.
	 *
	 * @param {{id: number; parentId: number}[]} tenantsArr
	 */
	convertToHierarchy(tenantsArr) {
		/** @type {Record<PropertyKey, HierarchyTenant} */
		const map = {};

		for (let i = 0; i < tenantsArr.length; ++i) {
			const obj = tenantsArr[i];
			obj.children = [];

			map[obj.id] = obj;

			const parent = (i === 0) ? "-" : obj.parentId; // the first item is always the top-level tenant
			if (!map[parent]) {
				map[parent] = {
					children: []
				};
			}
			map[parent].children.push(obj);
		}

		return map["-"].children;
	}

	/**
	 * Adds level information to a group of tenants. This manipulates the
	 * passed collection *in place* rather than returning a new structure.
	 *
	 * @param {Tenant[]} tenants
	 */
	addLevels(tenants) {
		// convert a flat tenant list into a hierarchy
		const tenantHierarchy = this.convertToHierarchy(tenants);
		// and then walk down the hierarchy to find out how deeply nested each tenant is and add a 'level' property to each tenant
		this.applyLevels(tenantHierarchy, 1);
	}

}

TenantUtils.$inject = [];
module.exports = TenantUtils;
