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

import org.apache.traffic_control.traffic_router.core.loc.Federation;
import org.apache.traffic_control.traffic_router.core.loc.FederationMapping;
import org.apache.traffic_control.traffic_router.core.loc.FederationRegistry;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class FederationExporter {

	private FederationRegistry federationRegistry;

	public List<Object> getMatchingFederations(final CidrAddress cidrAddress) {
		final List<Object> federationsList = new ArrayList<Object>();

		for (final Federation federation : federationRegistry.findFederations(cidrAddress)) {
			final List<Map<String, Object>> filteredFederationMappings = new ArrayList<Map<String, Object>>();

			for (final FederationMapping federationMapping : federation.getFederationMappings()) {
				filteredFederationMappings.add(getMappingProperties(cidrAddress, federationMapping));
			}

			final Map<String, Object> federationProperties = new HashMap<String, Object>();
			federationProperties.put("deliveryService", federation.getDeliveryService());
			federationProperties.put("federationMappings", filteredFederationMappings);

			federationsList.add(federationProperties);
		}
		return federationsList;
	}

	private Map<String, Object> getMappingProperties(final CidrAddress cidrAddress, final FederationMapping federationMapping) {
		final FederationMapping filteredMapping = federationMapping.createFilteredMapping(cidrAddress);
		final Map<String, Object> properties = new HashMap<String, Object>();

		properties.put("cname", filteredMapping.getCname());
		properties.put("ttl", filteredMapping.getTtl());

		addAddressProperties("resolve4", filteredMapping.getResolve4(), properties);
		addAddressProperties("resolve6", filteredMapping.getResolve6(), properties);

		return properties;
	}

	private Map<String, Object> addAddressProperties(final String propertyName, final ComparableTreeSet<CidrAddress> cidrAddresses, final Map<String, Object> properties) {
		final List<String> addressStrings = new ArrayList<String>();

		if (cidrAddresses == null || cidrAddresses.isEmpty()) {
			return properties;
		}

		for (final CidrAddress cidrAddress : cidrAddresses) {
			addressStrings.add(cidrAddress.getAddressString());
		}

		properties.put(propertyName, addressStrings);
		return properties;
	}

	public void setFederationRegistry(final FederationRegistry federationRegistry) {
		this.federationRegistry = federationRegistry;
	}
}
