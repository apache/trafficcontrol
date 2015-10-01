package com.comcast.cdn.traffic_control.traffic_router.core.util;

import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkNodeException;

import java.net.Inet4Address;
import java.net.Inet6Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Arrays;

public class CidrAddress implements Comparable<CidrAddress> {
    private final byte[] hostBytes;
    private final byte[] maskBytes;
    private final int netmaskLength;
    private final String ipString;

    @SuppressWarnings("PMD.CyclomaticComplexity")
    public CidrAddress(final String cidrString) throws NetworkNodeException {
        this.ipString = cidrString;
        final String[] hostNetworkArray = cidrString.split("/");
		final InetAddress address;

		try {
			address = InetAddress.getByName(hostNetworkArray[0]);
		} catch (UnknownHostException ex) {
			throw new NetworkNodeException(ex);
		}

		final byte[] addressBytes = address.getAddress();

		if (hostNetworkArray.length == 1) {
			netmaskLength = addressBytes.length * 8;
		} else {
			netmaskLength = Integer.parseInt(hostNetworkArray[1]);
		}

		if (address instanceof Inet4Address && (netmaskLength > 32 || netmaskLength < 0)) {
			throw new NetworkNodeException("Rejecting IPv4 subnet with invalid netmask: " + cidrString);
		} else if (address instanceof Inet6Address && (netmaskLength > 128 || netmaskLength < 0)) {
			throw new NetworkNodeException("Rejecting IPv6 subnet with invalid netmask: " + cidrString);
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

    @Override
	public int compareTo(final CidrAddress other) {
		// returns zero if the other address is in the same network as this one.
		byte[] mask = this.maskBytes;
		int len = netmaskLength;

		if (netmaskLength > other.netmaskLength) {
			mask = other.maskBytes;
			len = other.netmaskLength;
		}

		final int numNetmaskBytes = (int) Math.ceil((double) len / 8);

		for(int i = 0; i < numNetmaskBytes; i++) {
			final int diff = (hostBytes[i] & mask[i]) - (other.hostBytes[i] & mask[i]);
			if(diff != 0) { return diff; }
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
        return !(ipString != null ? !ipString.equals(that.ipString) : that.ipString != null);

    }

    @Override
    public int hashCode() {
        int result = hostBytes != null ? Arrays.hashCode(hostBytes) : 0;
        result = 31 * result + (maskBytes != null ? Arrays.hashCode(maskBytes) : 0);
        result = 31 * result + netmaskLength;
        result = 31 * result + (ipString != null ? ipString.hashCode() : 0);
        return result;
    }

    @Override
    public String toString() {
        return "CidrAddress{" + ipString + "}";
    }
}
