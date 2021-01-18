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

/** LocalizationMethod values are those allowed in the 'localizationMethods' of CacheGroups */
export const enum LocalizationMethod {
	/** Coverage Zone file lookup. */
	CZ = "CZ",
	/** Deep Coverage Zone file lookup. */
	DEEP_CZ = "DEEP_CZ",
	/** Geographic database search. */
	GEO = "GEO"
}

/**
 * Converts a LocalizationMethod to a human-readable string.
 *
 * @param l The LocalizationMethod to convert.
 * @returns A textual representation of 'l'.
 */
export function localizationMethodToString(l: LocalizationMethod): string {
	switch (l) {
		case LocalizationMethod.CZ:
			return "Coverage Zone File";
		case LocalizationMethod.DEEP_CZ:
			return "Deep Coverage Zone File";
		case LocalizationMethod.GEO:
			return "Geo-IP Database";
	}
}

/**
 * Represents a Cache Group.
 *
 * Refer to https://traffic-control-cdn.readthedocs.io/en/latest/overview/cache_groups.html
 */
export interface CacheGroup {
	fallbacks: Array<string>;
	fallbackToClosest: boolean;
	readonly id?: number;
	lastUpdated?: Date;
	latitude: number;
	localizationMethods: Array<LocalizationMethod>;
	longitude: number;
	name: string;
	parentCacheGroupID: number | null;
	parentCacheGroupName: string | null;
	secondaryParentCacheGroupID: number | null;
	secondaryParentCacheGroupName: string | null;
	shortName: string;
	typeId: number;
	typeName: string;
}
