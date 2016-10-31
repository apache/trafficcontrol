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

import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.util.AbstractResourceWatcher;
import org.apache.log4j.Logger;

import java.util.ArrayList;
import java.util.List;

public class SteeringWatcher extends AbstractResourceWatcher {
	private static final Logger LOGGER = Logger.getLogger(SteeringWatcher.class);
	private SteeringRegistry steeringRegistry;
	private CacheRegister cacheRegister;

	public static final String DEFAULT_STEERING_DATA_URL = "https://${toHostname}/internal/api/1.2/steering.json";

	public SteeringWatcher() {
		setDatabaseUrl(DEFAULT_STEERING_DATA_URL);
	}

	public void setCacheRegister(final CacheRegister cacheRegister) {
		this.cacheRegister = cacheRegister;
	}

	@Override
	public boolean useData(final String data) {
		try {
			steeringRegistry.update(data);
			if (cacheRegister != null) {
				final List<String> invalidOnes = new ArrayList<String>();

				for (final Steering steering : steeringRegistry.getAll()) {
					if (cacheRegister.getDeliveryService(steering.getDeliveryService()) == null) {
						LOGGER.warn("Steering data from " + dataBaseURL + " contains delivery service id reference '" + steering.getDeliveryService() + "' that's not in cr-config");
						invalidOnes.add(steering.getDeliveryService());
					}
				}

				steeringRegistry.removeAll(invalidOnes);
			}

			return true;
		} catch (Exception e) {
			LOGGER.warn("Failed updating steering registry with data from " + dataBaseURL);
		}

		return false;
	}

	@Override
	protected boolean verifyData(final String data) {
		try {
			return steeringRegistry.verify(data);
		} catch (Exception e) {
			LOGGER.warn("Failed to build steering data while verifying");
		}

		return false;
	}

	@Override
	public String getWatcherConfigPrefix() {
		return "steeringmapping";
	}

	public void setSteeringRegistry(final SteeringRegistry steeringRegistry) {
		this.steeringRegistry = steeringRegistry;
	}
}
