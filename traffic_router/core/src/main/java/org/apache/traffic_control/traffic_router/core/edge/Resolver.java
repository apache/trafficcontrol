/*
 *
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

package org.apache.traffic_control.traffic_router.core.edge;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;

public class Resolver {
	private static final Logger LOGGER = LogManager.getLogger(Resolver.class);

	public List<InetRecord> resolve(final String fqdn) {
		List<InetRecord> ipAddresses = null;
   		try {
			final InetAddress[] addresses = Inet4Address.getAllByName(fqdn);
			ipAddresses = new ArrayList<InetRecord>();
			for (final InetAddress address : addresses) {
				if (!address.isAnyLocalAddress() && !address.isLoopbackAddress() && !address.isLinkLocalAddress()
						&& !address.isMulticastAddress()) {
					ipAddresses.add(new InetRecord(address, 0));
				}
			}
			if (ipAddresses.isEmpty()) {
				LOGGER.info(String.format("No public addresses found for: (%s)", fqdn));
//				ipAddresses = null; // jlaue - give it a chance to recover next time?  
			}
		} catch (final UnknownHostException e) {
			LOGGER.warn(String.format("Unable to determine IP Address for: (%s)", fqdn));
		}
   		return ipAddresses;
	}

}
