/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.cache;

import java.net.Inet6Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.SortedSet;
import java.util.TreeSet;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.log4j.Logger;
import org.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_router.core.hash.HashFunction;
import com.comcast.cdn.traffic_control.traffic_router.core.hash.MD5HashFunction;

public class Cache implements Comparable<Cache> {
	private static final Logger LOGGER = Logger.getLogger(Cache.class);
	private static final int REPLICAS = 1000;

	/*
	 * Configuration Attributes
	 */
	private final String id;
	private String fqdn;
	private List<InetRecord> ipAddresses;
	private int port;
	private Collection<DeliveryServiceReference> deliveryServices = new ArrayList<DeliveryServiceReference>();
	final private List<Double> hashValues;

	final private int replicas;

	/**
	 * Creates a new {@link Cache}.
	 * 
	 * @param id
	 *            the id of the new cache
	 * @param hashCount 
	 */
	public Cache(final String id, final String hashId, final int hashCount) {
		this.id = id;

		final SortedSet<Double> sorter = new TreeSet<Double>();
		final HashFunction hash = new MD5HashFunction();
		replicas = (hashCount==0)? REPLICAS : hashCount;
		for (int i = 0; i < replicas; i++) {
			sorter.add(hash.hash(hashId + "--" + i));
			// hashValues.add(hash.hash(id + i));
		}
		hashValues = new ArrayList<Double>(sorter);
	}

	@Override
	public int compareTo(final Cache o) {
		return getId().compareTo(o.getId());
	}

	@Override
	public boolean equals(final Object obj) {
		if (this == obj) {
			return true;
		} else if (obj instanceof Cache) {
			final Cache rhs = (Cache) obj;
			return new EqualsBuilder()
			.append(getId(), rhs.getId())
			.isEquals();
		} else {
			return false;
		}
	}

	public Collection<DeliveryServiceReference> getDeliveryServices() {
		return deliveryServices;
	}

	public String getFqdn() {
		return fqdn;
	}

	public List<Double> getHashValues() {
		return hashValues;
	}

	public double getClosestHash(final double hash) {
		// assume hashValues sorted
		int hi = hashValues.size() -1;
		int lo = 0;
		int i = (hi-lo)/2;

		// you can tell a match if it's closer to hash than it's neighbors 
		for(int j = 0; j < replicas; j++) { // j is just for an escape hatch, should be found O(log(REPLICAS))
			final int r = match(hashValues, i, hash);
			if(r==0) {
				return hashValues.get(i);
			}
			if(r < 0) {
				hi = i-1;
			} else {
				lo = i+1;
			}
			i = (hi+lo)/2;
		}
		return 0;
	}

	private int match(final List<Double> a, final int i, final double hash) {
		// you can tell a match if it's closer to hash than it's neighbors 
		final double v = a.get(i).doubleValue();
		if(i+1 < a.size() && Math.abs(hash - a.get(i+1).doubleValue() ) < Math.abs(hash-v)) {
			return 1; // closer to hi neighbor
		}
		if(i-1 >= 0 && Math.abs(hash - a.get(i-1).doubleValue() ) < Math.abs(hash-v)) {
			return -1; // closer to lo neighbor
		}
		return 0; // match!
	}

	public String getId() {
		return id;
	}

	public List<InetRecord> getIpAddresses(final JSONObject ttls, final Resolver resolver) {
		return getIpAddresses(ttls, resolver, true);
	}

	public List<InetRecord> getIpAddresses(final JSONObject ttls, final Resolver resolver, final boolean ip6RoutingEnabled) {
		if(ipAddresses == null || ipAddresses.isEmpty()) {
			ipAddresses = resolver.resolve(this.getFqdn()+".");
		}
		if(ipAddresses == null) { return null; }
		final List<InetRecord> ret = new ArrayList<InetRecord>();
		for(InetRecord ir : ipAddresses) {
			if (ir.isInet6() && !ip6RoutingEnabled) {
				continue;
			}

			long ttl = 0;

			if(ttls == null) {
				ttl = -1;
			} else if(ir.isInet6()) {
				ttl = ttls.optLong("AAAA");
			} else {
				ttl = ttls.optLong("A");
			}

			ret.add(new InetRecord(ir.getAddress(), ttl));
		}
		return ret;
	}

//	static Resolver resolver = new Resolver();
//	private static Resolver getResolver() {
//		return resolver;
//	}
//	public static void setResolver(final Resolver r) {
//		resolver = r;
//	}

	public int getPort() {
		return port;
	}

	@Override
	public int hashCode() {
		return new HashCodeBuilder(1, 31)
		.append(getId())
		.toHashCode();
	}

	public void setDeliveryServices(final Collection<DeliveryServiceReference> deliveryServices) {
		this.deliveryServices = deliveryServices;
	}

	public void setFqdn(final String fqdn) {
		this.fqdn = fqdn;
	}

	public void setIpAddresses(final List<InetRecord> ipAddresses) {
		this.ipAddresses = ipAddresses;
	}

	public void setPort(final int port) {
		this.port = port;
	}

	@Override
	public String toString() {
		return "Cache [id=" + id + "] ";
	}

	/**
	 * Status enumeration for administratively reported status.
	 */
//	public enum AdminStatus {
//		ONLINE, OFFLINE, REPORTED, ADMIN_DOWN
//	}

	/**
	 * Contains a reference to a DeliveryService ID and the FQDN that should be used if this Cache
	 * is used when supporting the DeliveryService.
	 */
	public static class DeliveryServiceReference {
		private final String deliveryServiceId;
		private final String fqdn;

		public DeliveryServiceReference(final String deliveryServiceId, final String fqdn) {
			this.deliveryServiceId = deliveryServiceId;
			this.fqdn = fqdn;
		}

		public String getDeliveryServiceId() {
			return deliveryServiceId;
		}

		public String getFqdn() {
			return fqdn;
		}
	}

	boolean isAvailable = false;
	boolean hasAuthority = false;
	public void setIsAvailable(final boolean isAvailable) {
		this.hasAuthority = true;
		this.isAvailable = isAvailable;
	}
	public boolean hasAuthority() {
		return hasAuthority;
	}
	public boolean isAvailable() {
		return isAvailable;
	}
	InetAddress ip4;
	InetAddress ip6;
	public void setIpAddress(final String ip, final String ip6, final long ttl) throws UnknownHostException {
		this.ipAddresses = new ArrayList<InetRecord>();

		if (ip != null && !ip.isEmpty()) {
			this.ip4 = InetAddress.getByName(ip);
			ipAddresses.add(new InetRecord(ip4, ttl));
		} else {
			LOGGER.error(getFqdn() + " - no IPv4 address configured!");
		}

		if (ip6 != null && !ip6.isEmpty()) {
			final String ip6addr = ip6.replaceAll("/.*", "");
			this.ip6 = Inet6Address.getByName(ip6addr);
			ipAddresses.add(new InetRecord(this.ip6, ttl));
		} else {
			LOGGER.error(getFqdn() + " - no IPv6 address configured!");
		}
	}
	public InetAddress getIp4() {
		return ip4;
	}
	public InetAddress getIp6() {
		return ip6;
	}

	public void setState(final JSONObject state) {
		boolean isAvailable = true;
		if(state != null && state.has("isAvailable")) {
			isAvailable = state.optBoolean("isAvailable");
		}
		this.setIsAvailable(isAvailable);
	}
}
