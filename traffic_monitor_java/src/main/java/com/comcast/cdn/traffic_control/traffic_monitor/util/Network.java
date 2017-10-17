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

package com.comcast.cdn.traffic_control.traffic_monitor.util;

import java.net.InetAddress;
import java.net.NetworkInterface;
import java.net.SocketException;
import java.net.UnknownHostException;
import java.util.Enumeration;

import org.apache.log4j.Logger;

public class Network {
	private static final Logger LOGGER = Logger.getLogger(Network.class);

	public static final boolean isIpAddressLocal(final String ip) {
		try {
			final InetAddress address = InetAddress.getByName(ip);
			final Enumeration<NetworkInterface> ifaceList = NetworkInterface.getNetworkInterfaces();

			while (ifaceList.hasMoreElements()) {
				final NetworkInterface iface = ifaceList.nextElement();
				final Enumeration<InetAddress> addressList = iface.getInetAddresses();

				while (addressList.hasMoreElements()) {
					final InetAddress thisAddress = addressList.nextElement();

					if (address.equals(thisAddress)) {
						LOGGER.debug(address + " found on " + iface.getName() + "; returning true");
						return true;
					}
				}
			}
		} catch (UnknownHostException ex) {
			LOGGER.fatal(ex, ex);
		} catch (SocketException ex) {
			LOGGER.fatal(ex, ex);
		}

		return false;
	}

	public static boolean isLocalName(final String name) {
		try {
			if (name.equals(InetAddress.getLocalHost().getHostName())) {
				return true;
			}
		} catch (UnknownHostException ex) {
			LOGGER.fatal(ex, ex);
		}

		return false;
	}
}
