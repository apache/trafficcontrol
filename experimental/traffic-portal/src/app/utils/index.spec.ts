/**
 * @license Apache-2.0
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
import { TestBed } from "@angular/core/testing";

import { hasProperty, isArray, isArrayBufferView, isBoolean, isNumber, isNumberArray, isRecord, isString, isStringArray } from ".";

describe("Typing utility functions", () => {
	beforeEach(() => TestBed.configureTestingModule({}));

	it("should check for property existence correctly", () => {
		let test = {};
		expect(hasProperty(test, "anything")).toBeFalse();
		expect(hasProperty(test, 0)).toBeFalse();

		test = {anything: "something"};
		expect(hasProperty(test, "anything")).toBeTrue();
		expect(hasProperty(test, 0)).toBeFalse();

		// eslint-disable-next-line @typescript-eslint/naming-convention
		test = {0: "something"};
		expect(hasProperty(test, "anything")).toBeFalse();
		expect(hasProperty(test, 0)).toBeTrue();
	});

	it("should check property types correctly", () => {
		const wrong = {wrong: "type"};
		const isNum = (x: unknown): x is number => typeof x === "number";
		expect(hasProperty(wrong, "wrong", isNum)).toBeFalse();
		const right = {right: 0};
		expect(hasProperty(right, "right", isNum)).toBeTrue();
	});

	it("should check for string types correctly", () => {
		let test: string | number = "";
		expect(isString(test)).toBeTrue();
		test = 5;
		expect(isString(test)).toBeFalse();
	});

	it("should check for numeric types correctly", () => {
		let test: string | number = "0";
		expect(isNumber(test)).toBeFalse();
		test = 5;
		expect(isNumber(test)).toBeTrue();
	});

	it("should check for boolean types correctly", () => {
		let test: string | boolean = "true";
		expect(isBoolean(test)).toBeFalse();
		test = true;
		expect(isBoolean(test)).toBeTrue();
	});

	it("should check ambiguous record types correctly", () => {
		expect(isRecord(null)).toBeFalse();
		expect(isRecord({})).toBeTrue();
	});

	it("should check homogenous record types correctly", () => {
		const passes = {
			all: "properties",
			are: "strings"
		};
		expect(isRecord(passes, "string")).toBeTrue();
		const numbers = {
			numeric: 1,
			only: 5.23e7,
			properties: 0x2e,
		};
		expect(isRecord(numbers, "number")).toBeTrue();
		const fails = {
			not: "all",
			properties: "are",
			strings: 0
		};
		expect(isRecord(fails, "string")).toBeFalse();
		const isZero = (x: unknown): x is 0 => x === 0;
		expect(isRecord(fails, isZero)).toBeFalse();
		const customPasses = {
			everything: 0,
			is: 0,
			zero: 0
		};
		expect(isRecord(customPasses, isZero)).toBeTrue();
	});

	it("should check homogenous array types correctly", () => {
		const a = {};
		expect(isArray(a)).toBeFalse();
		const b = new Array();
		expect(isArray(b)).toBeTrue();
		b.push(5);
		expect(isArray(b)).toBeTrue();
		expect(isArray(b, "number")).toBeTrue();
		expect(isArray(b, "string")).toBeFalse();
		b.push("test");
		expect(isArray(b)).toBeTrue();
		expect(isArray(b, "number")).toBeFalse();
		expect(isArray(b, "string")).toBeFalse();
		expect(isArray(b, (x): x is Array<number | string> => typeof x === "number" || typeof x === "string")).toBeTrue();
	});

	it("should verify that only objects can be records", ()=>{
		expect(isRecord(0, (_): _ is unknown => true)).toBeFalse();
	});

	it("should be able to verify existence and type of array properties", () => {
		const a = {
			test: new Array<number|string>()
		};
		expect(hasProperty(a, "test", isArray)).toBeTrue();
		a.test.push("test");
		expect(hasProperty(a, "test", isStringArray)).toBeTrue();
		expect(hasProperty(a, "test", isNumberArray)).toBeFalse();
		a.test.push(5);
		expect(hasProperty(a, "test", isArray)).toBeTrue();
		expect(hasProperty(a, "test", isStringArray)).toBeFalse();
		expect(hasProperty(a, "test", isNumberArray)).toBeFalse();
	});

	it("knows if an object is an ArrayBufferView", () => {
		expect(isArrayBufferView(undefined)).toBeFalse();
		expect(isArrayBufferView(undefined)).toBeFalse();
		expect(isArrayBufferView(undefined)).toBeFalse();
		expect(isArrayBufferView(new Int8Array())).toBeTrue();
		expect(isArrayBufferView(new Uint8Array())).toBeTrue();
		expect(isArrayBufferView(new Uint8ClampedArray())).toBeTrue();
		expect(isArrayBufferView(new Int16Array())).toBeTrue();
		expect(isArrayBufferView(new Uint16Array())).toBeTrue();
		expect(isArrayBufferView(new Int32Array())).toBeTrue();
		expect(isArrayBufferView(new Uint32Array())).toBeTrue();
		expect(isArrayBufferView(new DataView(new ArrayBuffer(0)))).toBeTrue();
		expect(isArrayBufferView([1, 2, 3])).toBeFalse();
		expect(isArrayBufferView([0x01n, 2n, 3n])).toBeFalse();
	});
});
