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
import { orderBy } from "./order-by";

describe("orderBy", () => {
	it("sorts single properties properly", () => {
		const input = [{bar: 2, foo: 1}, {bar: 1, foo: 1}];
		const output = orderBy(input, "bar");
		expect(output[0].bar).toEqual(1);
		expect(output[1].bar).toEqual(2);
	});

	it("sorts multiple properties properly", () => {
		const input = [{bar: 2, foo: 2}, {bar: 3, foo: 1}, {bar: 1, foo: 3}, {bar: 3, foo: 2}];
		const output = orderBy(input, ["bar", "foo"]);
		expect(output[0].bar).toEqual(1);
		expect(output[1].bar).toEqual(2);
		expect(output[2].bar).toEqual(3);
		expect(output[2].foo).toEqual(1);
		expect(output[3].bar).toEqual(3);
		expect(output[3].foo).toEqual(2);
	});

	it("handles null properties", () => {
		const input = [{foo: 2}, {foo: null}, {foo: 1}];
		const output = orderBy(input, "foo");
		expect(output[0].foo).toEqual(1);
		expect(output[1].foo).toEqual(2);
		expect(output[2].foo).toBeNull();
	});

	it("ignores values of differing types", () => {
		const input = [{foo: 2}, {foo: "1"}, {foo: 3}, {foo: 1}];
		const output = orderBy(input, "foo");
		expect(output[0].foo).toEqual(2);
		expect(output[1].foo).toEqual("1");
		expect(output[2].foo).toEqual(1);
		expect(output[3].foo).toEqual(3);
	});

	it("won't order by a property that objects don't have", () => {
		const input = [{foo: 2}, {foo: 1}];
		const output = orderBy(input, "bar");
		expect(output[0].foo).toEqual(2);
		expect(output[1].foo).toEqual(1);
	});

	it("doesn't do any ordering when there are no ordering properties provided", () => {
		const input = [{foo: 3}, {foo: 2}, {foo: 1}];
		const output = orderBy(input, []);
		expect(output[0].foo).toEqual(3);
		expect(output[1].foo).toEqual(2);
		expect(output[2].foo).toEqual(1);
	});

	it("sorts null values to be last", ()=>{
		let input = [{foo: 1}, {foo: -1}, {foo: null}];
		let output = orderBy(input, "foo");
		expect(output[0].foo).toBe(-1);
		expect(output[1].foo).toBe(1);
		expect(output[2].foo).toBeNull();

		input = input.reverse();
		output = orderBy(input, "foo");
		expect(output[0].foo).toBe(-1);
		expect(output[1].foo).toBe(1);
		expect(output[2].foo).toBeNull();

		input = [{foo: Infinity}, {foo: -Infinity}, {foo: null}, {foo: null}];
		output = orderBy(input, "foo");
		expect(output[0].foo).toBeNegativeInfinity();
		expect(output[1].foo).toBePositiveInfinity();
		expect(output[2].foo).toBeNull();
		expect(output[3].foo).toBeNull();
	});
});
