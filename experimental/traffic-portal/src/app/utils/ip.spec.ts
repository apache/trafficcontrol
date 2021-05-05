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
import { serviceInterface } from "./ip";

describe("serviceInterface", () => {
	it("finds a service interface in a group of interfaces", () => {
		// the first interface is the service interface
		let infs = [
			{
				ipAddresses: [
					{
						address: "192.0.2.0",
						gateway: null,
						serviceAddress: true
					}
				],
				maxBandwidth: null,
				monitor: false,
				mtu: null,
				name: "correct"
			},
			{
				ipAddresses: [
					{
						address: "192.0.2.1",
						gateway: null,
						serviceAddress: false
					}
				],
				maxBandwidth: null,
				monitor: false,
				mtu: null,
				name: "wrong"
			},
			{
				ipAddresses: [
					{
						address: "192.0.2.2",
						gateway: null,
						serviceAddress: false
					}
				],
				maxBandwidth: null,
				monitor: false,
				mtu: null,
				name: "incorrect"
			}
		];
		expect(serviceInterface(infs).name).toBe("correct");

		// the middle interface is the service interface
		infs = [
			{
				ipAddresses: [
					{
						address: "192.0.2.0",
						gateway: null,
						serviceAddress: false
					}
				],
				maxBandwidth: null,
				monitor: false,
				mtu: null,
				name: "wrong"
			},
			{
				ipAddresses: [
					{
						address: "192.0.2.1",
						gateway: null,
						serviceAddress: true
					}
				],
				maxBandwidth: null,
				monitor: false,
				mtu: null,
				name: "correct"
			},
			{
				ipAddresses: [
					{
						address: "192.0.2.2",
						gateway: null,
						serviceAddress: false
					}
				],
				maxBandwidth: null,
				monitor: false,
				mtu: null,
				name: "incorrect"
			}
		];
		expect(serviceInterface(infs).name).toBe("correct");

		// the last interface is the service interface
		infs = [
			{
				ipAddresses: [
					{
						address: "192.0.2.0",
						gateway: null,
						serviceAddress: false
					}
				],
				maxBandwidth: null,
				monitor: false,
				mtu: null,
				name: "wrong"
			},
			{
				ipAddresses: [
					{
						address: "192.0.2.1",
						gateway: null,
						serviceAddress: false
					}
				],
				maxBandwidth: null,
				monitor: false,
				mtu: null,
				name: "incorrect"
			},
			{
				ipAddresses: [
					{
						address: "192.0.2.2",
						gateway: null,
						serviceAddress: true
					}
				],
				maxBandwidth: null,
				monitor: false,
				mtu: null,
				name: "correct"
			}
		];
		expect(serviceInterface(infs).name).toBe("correct");
	});

	it("throws an error if there is no service address", () => {
		expect(()=>serviceInterface([])).toThrow();
	});
});
