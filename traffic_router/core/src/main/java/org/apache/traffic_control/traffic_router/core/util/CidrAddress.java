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

package org.apache.traffic_control.traffic_router.core.util;

import org.apache.traffic_control.traffic_router.core.loc.NetworkNodeException;

import java.net.Inet4Address;
import java.net.Inet6Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Arrays;

public class CidrAddress implements Comparable<CidrAddress> {
    private final byte[] hostBytes;
    private final byte[] maskBytes;
    private final int netmaskLength;
    private final InetAddress address;

    public static CidrAddress fromString(final String cidrString) throws NetworkNodeException {
        final String[] hostNetworkArray = cidrString.split("/");
        final String host = hostNetworkArray[0].trim();

        InetAddress address;
        try {
            address = InetAddress.getByName(host);
        } catch (UnknownHostException ex) {
            throw new NetworkNodeException(ex);
        }

        if (hostNetworkArray.length == 1) {
            return new CidrAddress(address);
        }

        int netmaskLength;
        try {
            netmaskLength = Integer.parseInt(hostNetworkArray[1].trim());
        }
        catch (NumberFormatException e) {
            throw new NetworkNodeException(e);
        }

        return new CidrAddress(address, netmaskLength);
    }

    public CidrAddress(final InetAddress address) throws NetworkNodeException {
        this(address, address.getAddress().length * 8);
    }

    @SuppressWarnings("PMD.CyclomaticComplexity")
    public CidrAddress(final InetAddress address, final int netmaskLength) throws NetworkNodeException {
        this.netmaskLength = netmaskLength;
        this.address = address;
        final byte[] addressBytes = address.getAddress();

        if (address instanceof Inet4Address && (netmaskLength > 32 || netmaskLength < 0)) {
            throw new NetworkNodeException("Rejecting IPv4 subnet with invalid netmask: " + getCidrString());
        } else if (address instanceof Inet6Address && (netmaskLength > 128 || netmaskLength < 0)) {
            throw new NetworkNodeException("Rejecting IPv6 subnet with invalid netmask: " + getCidrString());
        }

        hostBytes = addressBytes;
        maskBytes = new byte[addressBytes.length];

        for (int i = 0; i < netmaskLength; i++) {
            maskBytes[i/8] |= 1<<(7-(i%8));
        }
    }

    public byte[] getHostBytes() {
        return hostBytes;
    }

    public byte[] getMaskBytes() {
        return maskBytes;
    }

    public int getNetmaskLength() {
        return netmaskLength;
    }

    public boolean includesAddress(final CidrAddress other) {
        if (netmaskLength >= other.netmaskLength) {
            return false;
        }

        return compareTo(other) == 0;
    }

    public boolean isIpV6() {
        return getHostBytes().length > 4;
    }

    @Override
	public int compareTo(final CidrAddress other) {
		byte[] mask = this.maskBytes;
		int len = netmaskLength;

		if (netmaskLength > other.netmaskLength) {
			mask = other.maskBytes;
			len = other.netmaskLength;
		}

		final int numNetmaskBytes = (int) Math.ceil((double) len / 8);

        for (int i = 0; i < numNetmaskBytes; i++) {
            final int diff = (hostBytes[i] & mask[i]) - (other.hostBytes[i] & mask[i]);
			if (diff != 0) {
                return diff;
            }
		}

        return 0;
    }

    @Override
    @SuppressWarnings({"PMD.NPathComplexity", "PMD.IfStmtsMustUseBraces"})
    public boolean equals(final Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;

        final CidrAddress that = (CidrAddress) o;

        if (netmaskLength != that.netmaskLength) return false;
        if (!Arrays.equals(hostBytes, that.hostBytes)) return false;
        if (!Arrays.equals(maskBytes, that.maskBytes)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        int result = hostBytes != null ? Arrays.hashCode(hostBytes) : 0;
        result = 31 * result + (maskBytes != null ? Arrays.hashCode(maskBytes) : 0);
        result = 31 * result + netmaskLength;
        return result;
    }

    private String getCidrString() {
        return "CidrAddress{" + address.toString() + "/" + netmaskLength + "}";
    }

    @Override
    public String toString() {
        return getCidrString();
    }

    public String getAddressString() {
        return address.toString() + "/" + netmaskLength;
    }
}
