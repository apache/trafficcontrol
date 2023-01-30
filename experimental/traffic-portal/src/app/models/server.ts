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

/**
 * This file is for modeling and functionality related to Server objects
 */

/** IPAddress is a single IP address of a single network interface of a server. */
export interface IPAddress {
	/** The actual IP address. */
	address: string;
	/** The IP address of a gateway for this IP, if one exists/is known. */
	gateway: string | null;
	/** Whether or not this IP address is responsible for serving ATC traffic. */
	serviceAddress: boolean;
}

/** Interface is a server's network interface. */
export interface Interface {
	/** The IP addresses assigned to this network interface. */
	ipAddresses: Array<IPAddress>;
	/** The maximum bandwidth for considering the server healthy, if any. */
	maxBandwidth: number | null;
	/**
	 * Whether or not the Traffic Monitor should consider this network interface
	 * in health calculations.
	 */
	monitor: boolean;
	/** The maximum transmission unit of the network interface, if known. */
	mtu: number | null;
	/** The name of the network interface in the `/dev` directory. */
	name: string;
}
