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

package org.apache.traffic_control.traffic_router.core.ds;

import org.apache.traffic_control.traffic_router.core.util.AbstractResourceWatcher;
import org.apache.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class SteeringWatcher extends AbstractResourceWatcher {
	private static final Logger LOGGER = LogManager.getLogger(SteeringWatcher.class);
	private SteeringRegistry steeringRegistry;

	public static final String DEFAULT_STEERING_DATA_URL = "https://${toHostname}/api/"+TrafficOpsUtils.TO_API_VERSION+"/steering";

	public SteeringWatcher() {
		setDatabaseUrl(DEFAULT_STEERING_DATA_URL);
		setDefaultDatabaseUrl(DEFAULT_STEERING_DATA_URL);
	}

	@Override
	public boolean useData(final String data) {
		try {
			// NOTE: it is likely that the steering data will contain xml_ids for delivery services
			// that haven't been added to the CRConfig yet. This is okay because the SteeringRegistry
			// will only be queried for Delivery Service xml_ids that exist in CRConfig
			steeringRegistry.update(data);

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
