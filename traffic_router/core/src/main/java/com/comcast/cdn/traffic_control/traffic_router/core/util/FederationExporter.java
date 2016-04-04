package com.comcast.cdn.traffic_control.traffic_router.core.util;

import com.comcast.cdn.traffic_control.traffic_router.core.loc.Federation;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationMapping;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationRegistry;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class FederationExporter {

	private FederationRegistry federationRegistry;

	public List<Object> getMatchingFederations(final CidrAddress cidrAddress) {
		final List<Object> federationsList = new ArrayList<Object>();

		for (Federation federation : federationRegistry.findFederations(cidrAddress)) {
			final List<Map<String, Object>> filteredFederationMappings = new ArrayList<Map<String, Object>>();

			for (FederationMapping federationMapping : federation.getFederationMappings()) {
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

		for (CidrAddress cidrAddress : cidrAddresses) {
			addressStrings.add(cidrAddress.getAddressString());
		}

		properties.put(propertyName, addressStrings);
		return properties;
	}

	public void setFederationRegistry(final FederationRegistry federationRegistry) {
		this.federationRegistry = federationRegistry;
	}
}
