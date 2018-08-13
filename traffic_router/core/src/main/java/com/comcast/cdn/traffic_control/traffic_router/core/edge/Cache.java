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

package com.comcast.cdn.traffic_control.traffic_router.core.edge;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.List;
import java.util.ArrayList;

import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;

import com.comcast.cdn.traffic_control.traffic_router.core.config.ParseException;

public class Cache extends Node {
	private final Map<String, DeliveryServiceReference> deliveryServices = new HashMap<String, DeliveryServiceReference>();
	private final Geolocation geolocation;

	public Cache(final String id, final String hashId, final int hashCount, final Geolocation geolocation) {
		super(id, hashId, hashCount);
		this.geolocation = geolocation;
	}

	public Cache(final String id, final String hashId, final int hashCount) {
		this(id, hashId, hashCount, null);
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

	public DeliveryServiceReference getDeliveryService(final String dsid) {
		return deliveryServices.get(dsid);
	}

	public Geolocation getGeolocation() {
		return geolocation;
	}

	@Override
	public int hashCode() {
		return new HashCodeBuilder(1, 31)
		.append(getId())
		.toHashCode();
	}

	public void setDeliveryServices(final Collection<DeliveryServiceReference> deliveryServices) {
		this.deliveryServices.clear();
		for (final DeliveryServiceReference deliveryServiceReference : deliveryServices) {
			this.deliveryServices.put(deliveryServiceReference.getDeliveryServiceId(), deliveryServiceReference);
		}
	}

	public boolean hasDeliveryService(final String deliveryServiceId) {
		return deliveryServices.containsKey(deliveryServiceId);
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
		private final String host;
		private final String domain;
		private List<String> aliases = new ArrayList<>();

		public DeliveryServiceReference(final String deliveryServiceId, final String fqdn) throws ParseException {
			if (fqdn.split("\\.", 2).length != 2) {
				throw new ParseException("Invalid FQDN (" + fqdn + ") on delivery service " + deliveryServiceId + "; please verify the HOST regex(es) in Traffic Ops");
			}

			this.deliveryServiceId = deliveryServiceId;
			this.fqdn = fqdn;
			final String[] parts = fqdn.split("\\.", 2);
			host = parts[0];
			domain = parts[1];
		}

		public DeliveryServiceReference(final String deliveryServiceId, final String fqdn,
				final List<String> dsNames) throws ParseException {
			this(deliveryServiceId, fqdn);

			if (dsNames != null) {
				aliases = dsNames;
			}
		}

		public String getDeliveryServiceId() {
			return deliveryServiceId;
		}

		public String getFqdn() {
			return fqdn;
		}

		/*
		 * Aliases are only used to create and update the statTracker
		 */
		public List<String> getAliases() {
			return aliases;
		}

		public void setAliases(final List<String> aliasesList) {
			if (aliasesList == null) {
				return;
			}

			aliases = aliasesList;
		}

		public String getHost() {
			return host;
		}

		public String getDomain() {
			return domain;
		}
	}
}
