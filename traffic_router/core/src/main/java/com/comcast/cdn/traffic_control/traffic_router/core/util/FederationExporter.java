package com.comcast.cdn.traffic_control.traffic_router.core.util;

import com.comcast.cdn.traffic_control.traffic_router.core.loc.Federation;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationMapping;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationRegistry;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Component
public class FederationExporter {

	@Autowired
	private FederationRegistry federationRegistry;

	public List<Object> getMatchingFederations(CidrAddress cidrAddress) {
		List<Object> federationsList = new ArrayList<Object>();

		for (Federation federation : federationRegistry.findFederations(cidrAddress)) {
			List<Map<String, Object>> filteredFederationMappings = new ArrayList<Map<String, Object>>();

			for (FederationMapping federationMapping : federation.getFederationMappings()) {
				filteredFederationMappings.add(getMappingProperties(cidrAddress, federationMapping));
			}

			Map<String, Object> federationProperties = new HashMap<String, Object>();
			federationProperties.put("deliveryService", federation.getDeliveryService());
			federationProperties.put("federationMappings", filteredFederationMappings);

			federationsList.add(federationProperties);
		}
		return federationsList;
	}

	private Map<String, Object> getMappingProperties(CidrAddress cidrAddress, FederationMapping federationMapping) {
		final FederationMapping filteredMapping = federationMapping.createFilteredMapping(cidrAddress);
		Map<String, Object> properties = new HashMap<String, Object>();

		properties.put("cname", filteredMapping.getCname());
		properties.put("ttl", filteredMapping.getTtl());

		addstuff("resolve4", filteredMapping.getResolve4(), properties);
		addstuff("resolve6", filteredMapping.getResolve6(), properties);

		return properties;
	}

	private Map<String, Object> addstuff(String propertyName, ComparableTreeSet<CidrAddress> cidrAddresses, Map<String, Object> properties) {
		List<String> addressStrings = new ArrayList<String>();

		if (cidrAddresses.isEmpty()) {
			return properties;
		}

		for (CidrAddress cidrAddress4 : cidrAddresses) {
			addressStrings.add(cidrAddress4.getAddressString());
		}

		properties.put(propertyName, addressStrings);
		return properties;
	}

	public void setFederationRegistry(FederationRegistry federationRegistry) {
		this.federationRegistry = federationRegistry;
	}
}
