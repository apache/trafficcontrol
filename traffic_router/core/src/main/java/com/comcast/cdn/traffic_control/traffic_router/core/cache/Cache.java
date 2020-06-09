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

package com.comcast.cdn.traffic_control.traffic_router.core.cache;

import java.net.Inet6Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.comcast.cdn.traffic_control.traffic_router.core.hash.DefaultHashable;
import com.comcast.cdn.traffic_control.traffic_router.core.hash.Hashable;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtils;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;
import com.fasterxml.jackson.databind.JsonNode;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.log4j.Logger;

import com.comcast.cdn.traffic_control.traffic_router.core.config.ParseException;

public class Cache implements Comparable<Cache>, Hashable<Cache> {
	private static final Logger LOGGER = Logger.getLogger(Cache.class);
	private static final int REPLICAS = 1000;

	public enum IpVersions {
		IPV4ONLY, IPV6ONLY, BOTH
	}
	private IpVersions ipAvailableVersions;
	private final String id;
	private String fqdn;
	private List<InetRecord> ipAddresses;
	private List<InetRecord> unavailableIpAddresses;
	private int port;
	private final Map<String, DeliveryServiceReference> deliveryServices = new HashMap<String, DeliveryServiceReference>();
	private final Geolocation geolocation;
	private final Hashable hashable = new DefaultHashable();
	private int httpsPort = 443;

	public Cache(final String id, final String hashId, final int hashCount, final Geolocation geolocation) {
		this.id = id;
		hashable.generateHashes(hashId, hashCount > 0 ? hashCount : REPLICAS);
		this.geolocation = geolocation;
	}

	public Cache(final String id, final String hashId, final int hashCount) {
		this(id, hashId, hashCount, null);
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
		return deliveryServices.values();
	}

	public Geolocation getGeolocation() {
		return geolocation;
	}

	public String getFqdn() {
		return fqdn;
	}

	public String getId() {
		return id;
	}

	public List<InetRecord> getIpAddresses(final JsonNode ttls, final Resolver resolver) {
		return getIpAddresses(ttls, resolver, true);
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	public List<InetRecord> getIpAddresses(final JsonNode ttls, final Resolver resolver, final boolean ip6RoutingEnabled) {
		if(ipAddresses == null || ipAddresses.isEmpty()) {
			ipAddresses = resolver.resolve(this.getFqdn()+".");
		}
		if(ipAddresses == null) { return null; }
		final List<InetRecord> ret = new ArrayList<InetRecord>();
		for (final InetRecord ir : ipAddresses) {
			if (ir.isInet6() && !ip6RoutingEnabled) {
				continue;
			}

			long ttl = 0;

			if(ttls == null) {
				ttl = -1;
			} else if(ir.isInet6()) {
				ttl = JsonUtils.optLong(ttls, "AAAA");
			} else {
				ttl = JsonUtils.optLong(ttls, "A");

			}

			ret.add(new InetRecord(ir.getAddress(), ttl));
		}
		return ret;
	}

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
		for (final DeliveryServiceReference deliveryServiceReference : deliveryServices) {
			this.deliveryServices.put(deliveryServiceReference.getDeliveryServiceId(), deliveryServiceReference);
		}
	}

	public boolean hasDeliveryService(final String deliveryServiceId) {
		return deliveryServices.containsKey(deliveryServiceId);
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
	 * Contains a reference to a DeliveryService ID and the FQDN that should be used if this Cache
	 * is used when supporting the DeliveryService.
	 */
	public static class DeliveryServiceReference {
		private final String deliveryServiceId;
		private final String fqdn;

		public DeliveryServiceReference(final String deliveryServiceId, final String fqdn) throws ParseException {
			if (fqdn.split("\\.", 2).length != 2) {
				throw new ParseException("Invalid FQDN (" + fqdn + ") on delivery service " + deliveryServiceId + "; please verify the HOST regex(es) in Traffic Ops");
			}

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
		this.unavailableIpAddresses = new ArrayList<InetRecord>();

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

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	public void setState(final JsonNode state) {
		final boolean isAvailable = JsonUtils.optBoolean(state, "isAvailable", true);
		final boolean ipv4Available = JsonUtils.optBoolean(state, "ipv4Available", true);
		final boolean ipv6Available = JsonUtils.optBoolean(state, "ipv6Available", true);

		final List<InetRecord> newlyAvailable = new ArrayList<>();
		final List<InetRecord> newlyUnavailable = new ArrayList<>();

		if (ipv4Available && !ipv6Available) {
			this.ipAvailableVersions = IpVersions.IPV4ONLY;
		} else if (ipv6Available && !ipv4Available) {
			this.ipAvailableVersions = IpVersions.IPV6ONLY;
		} else {
			this.ipAvailableVersions = IpVersions.BOTH;
		}

		for (final InetRecord record : ipAddresses) {
			if (record.getAddress().equals(ip4) && !ipv4Available) {
				newlyUnavailable.add(record);
			}
			if (record.getAddress().equals(ip6) && !ipv6Available) {
				newlyUnavailable.add(record);
			}
		}

		for (final InetRecord record : unavailableIpAddresses) {
			if (record.getAddress().equals(ip4) && ipv4Available) {
				newlyAvailable.add(record);
			}
			if (record.getAddress().equals(ip6) && ipv6Available) {
				newlyAvailable.add(record);
			}
		}

		ipAddresses.addAll(newlyAvailable);
		ipAddresses.removeAll(newlyUnavailable);
		unavailableIpAddresses.addAll(newlyUnavailable);
		unavailableIpAddresses.removeAll(newlyAvailable);

		this.setIsAvailable(isAvailable);


	}

	@Override
	public Hashable<Cache> generateHashes(final String hashId, final int hashCount) {
		hashable.generateHashes(hashId, hashCount);
		return this;
	}

	@Override
	public double getClosestHash(final double hash) {
		return hashable.getClosestHash(hash);
	}

	@Override
	public List<Double> getHashValues() {
		return hashable.getHashValues();
	}

	public int getHttpsPort() {
		return httpsPort;
	}

	public void setHttpsPort(final int httpsPort) {
		this.httpsPort = httpsPort;
	}

	@Override
	public boolean hasHashes() {
		return hashable.hasHashes();
	}

	@Override
	public int getOrder() {
		return hashable.getOrder();
	}

	@Override
	public void setOrder(final int order) {
		hashable.setOrder(order);
	}

	public IpVersions getIpAvailableVersions() {
		return ipAvailableVersions;
	}

	public void setIpAvailableVersions(final IpVersions ipAvailableVersions) {
		this.ipAvailableVersions = ipAvailableVersions;
	}
}
