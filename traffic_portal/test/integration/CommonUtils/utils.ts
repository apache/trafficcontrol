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
 * hasProperty checks whether some variable passed as `o` has the string
 * property `k`.
 *
 * @example
 * hasProperty({}, "id", "string"); // returns false
 * hasProperty({id: 8}, "id", "string"); // returns false
 * hasProperty({id: undefined}, "id", "string"); // returns false
 *
 * @param o The object to check.
 * @param k The key for which to check in the object.
 * @param type Specifies that `o.k` should be a string.
 * @returns Whether or not `o` has the string property `k`.
 */
export function hasProperty<T extends object, K extends PropertyKey>(o: T, k: K, type: "string"): o is T & Record<K, string>;
/**
 * hasProperty checks whether some variable passed as `o` has the number
 * property `k`.
 *
 * @example
 * hasProperty({}, "id", "number"); // returns false
 * hasProperty({id: 8}, "id", "number"); // returns true
 * hasProperty({id: undefined}, "id", "number"); // returns false
 *
 * @param o The object to check.
 * @param k The key for which to check in the object.
 * @param type Specifies that `o.k` should be a number.
 * @returns Whether or not `o` has the number property `k`.
 */
export function hasProperty<T extends object, K extends PropertyKey>(o: T, k: K, type: "number"): o is T & Record<K, number>;
/**
 * hasProperty checks whether some variable passed as `o` has the boolean
 * property `k`.
 *
 * @example
 * hasProperty({}, "id", "boolean"); // returns false
 * hasProperty({id: 8}, "id", "boolean"); // returns false
 * hasProperty({id: undefined}, "id", "boolean"); // returns false
 * hasProperty({id: true}, "id", "boolean"); // returns true
 *
 * @param o The object to check.
 * @param k The key for which to check in the object.
 * @param type Specifies that `o.k` should be a boolean.
 * @returns Whether or not `o` has the boolean property `k`.
 */
export function hasProperty<T extends object, K extends PropertyKey>(o: T, k: K, type: "boolean"): o is T & Record<K, boolean>;
/**
 * hasProperty checks whether some variable passed as `o` has the Array property
 * `k`.
 *
 * @example
 * hasProperty({}, "id", "Array"); // returns false
 * hasProperty({id: 8}, "id", "Array"); // returns false
 * hasProperty({id: undefined}, "id", "Array"); // returns false
 * hasProperty({id: []}, "id", "Array"); // returns true
 * hasProperty({id: [undefined, null, -7, true]}, "id", "Array"); // returns true
 *
 * @param o The object to check.
 * @param k The key for which to check in the object.
 * @param type Specifies that `o.k` should be a (potentially non-homogenous)
 * Array.
 * @returns Whether or not `o` has the Array property `k`.
 */
export function hasProperty<T extends object, K extends PropertyKey>(o: T, k: K, type: "Array"): o is T & Record<K, Array<unknown>>;
/**
 * hasProperty checks, generically, whether some variable passed as `o` has the
 * property `k`.
 *
 * @example
 * hasProperty({}, "id"); // returns false
 * hasProperty({id: 8}, "id"); // returns true
 * hasProperty({id: undefined}, "id"); // returns true
 *
 * @param o The object to check.
 * @param k The key for which to check in the object.
 * @returns Whether or not `o` has the property `k`.
 */
export function hasProperty<T extends object, K extends PropertyKey>(o: T, k: K): o is T & Record<K, unknown>;
/**
 * hasProperty checks, generically, whether some variable passed as `o` has the
 * property `k`.
 *
 * @example
 * hasProperty({}, "id"); // returns false
 * hasProperty({id: 8}, "id"); // returns true
 * hasProperty({id: undefined}, "id"); // returns true
 * hasProperty({id: 8}, "id", "number"); // returns true
 * hasProperty({id: undefined}, "id", "number"); // returns false
 * hasProperty({id: 8}, "id", "string"); // returns false
 *
 * @param o The object to check.
 * @param k The key for which to check in the object.
 * @param type Optionally specify a type to check for.
 * @returns Whether or not `o` has the property `k`.
 */
export function hasProperty<T extends object, K extends PropertyKey, S>(o: T, k: K, type?: "string" | "number" | "boolean" | "Array"): o is T & Record<K, S> {
	if (!Object.prototype.hasOwnProperty.call(o, k)) {
		return false;
	}
	if (!type) {
		return true;
	}
	const val = (o as Record<K, unknown>)[k];
	switch (type) {
		case "string":
			return typeof(val) === "string";
		case "number":
			return typeof(val) === "number";
		case "boolean":
			return typeof(val) === "boolean";
		case "Array":
			return val instanceof Array;
	}
}
