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

package com.comcast.cdn.traffic_control.traffic_router.core.ds;

import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.log4j.Logger;

import java.io.IOException;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class SteeringRegistry {
	private static final Logger LOGGER = Logger.getLogger(SteeringRegistry.class);

	private final Map<String, Steering> registry = new HashMap<String, Steering>();
	private final ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.AvoidDuplicateLiterals"})
	public void update(final String json) {
		Map<String, List<Steering>> m;
		try {
			m = objectMapper.readValue(json, new TypeReference<HashMap<String, List<Steering>>>() { });
		} catch (IOException e) {
			LOGGER.error("Failed consuming Json data to populate steering registry, keeping current data:" + e.getMessage());
			return;
		}

		final List<Steering> steerings = m.values().iterator().next();
		final Map<String, Steering> newSteerings = new HashMap<String, Steering>();

		for (final Steering steering : steerings) {
			for (final SteeringTarget steeringTarget : steering.getTargets()) {
				steeringTarget.generateHashes();
			}
			newSteerings.put(steering.getDeliveryService(), steering);
		}

		registry.clear();
		registry.putAll(newSteerings);
		for (final Steering steering : steerings) {
			for (final SteeringTarget target : steering.getTargets()) {
				if (target.getGeolocation() != null && target.getGeoOrder() != 0) {
					LOGGER.info("Steering " + steering.getDeliveryService() + " target " + target.getDeliveryService() + " now has geolocation [" + target.getLatitude() + ", "  + target.getLongitude() + "] and geoOrder " + target.getGeoOrder());
				} else if (target.getGeolocation() != null && target.getWeight() > 0) {
					LOGGER.info("Steering " + steering.getDeliveryService() + " target " + target.getDeliveryService() + " now has geolocation [" + target.getLatitude() + ", "  + target.getLongitude() + "] and weight " + target.getWeight());
				} else if (target.getGeolocation() != null) {
					LOGGER.info("Steering " + steering.getDeliveryService() + " target " + target.getDeliveryService() + " now has geolocation [" + target.getLatitude() + ", "  + target.getLongitude() + "]");
				} else if (target.getWeight() > 0) {
					LOGGER.info("Steering " + steering.getDeliveryService() + " target " + target.getDeliveryService() + " now has weight " + target.getWeight());
				} else if (target.getOrder() != 0) { // this target has a specific order set
					LOGGER.info("Steering " + steering.getDeliveryService() + " target " + target.getDeliveryService() + " now has order " + target.getOrder());
				} else {
					LOGGER.info("Steering " + steering.getDeliveryService() + " target " + target.getDeliveryService() + " now has weight " + target.getWeight() + " and order " + target.getOrder());
				}
			}
		}
	}

	public boolean verify(final String json) {
		try {
			final ObjectMapper mapper = new ObjectMapper(new JsonFactory());
			mapper.readValue(json, new TypeReference<HashMap<String, List<Steering>>>() { });
		} catch (IOException e) {
			LOGGER.error("Failed consuming Json data to populate steering registry while verifying:" + e.getMessage());
			return false;
		}

		return true;
	}

	public boolean has(final String steeringId) {
		return registry.containsKey(steeringId);
	}

	public Steering get(final String steeringId) {
		return registry.get(steeringId);
	}

	public Collection<Steering> getAll() {
		return registry.values();
	}

	public Collection<Steering> removeAll(final List<String> steeringIds) {
		final List<Steering> removedEntries = new ArrayList<Steering>();

		for (final String steeringId : steeringIds) {
			removedEntries.add(registry.remove(steeringId));
		}

		return removedEntries;
	}
}