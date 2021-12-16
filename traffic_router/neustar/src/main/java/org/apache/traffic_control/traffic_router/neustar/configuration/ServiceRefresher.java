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

package org.apache.traffic_control.traffic_router.neustar.configuration;

import org.apache.traffic_control.traffic_router.neustar.NeustarGeolocationService;
import org.apache.traffic_control.traffic_router.neustar.data.NeustarDatabaseUpdater;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;

public class ServiceRefresher implements Runnable {
	private final Logger logger = LogManager.getLogger(ServiceRefresher.class);

	@Autowired
	NeustarDatabaseUpdater neustarDatabaseUpdater;

	@Autowired
	NeustarGeolocationService neustarGeolocationService;

	@Override
	public void run() {
		try {
			if (neustarDatabaseUpdater.update() || !neustarGeolocationService.isInitialized()) {
				neustarGeolocationService.reloadDatabase();
			}
		} catch (Exception e) {
			logger.error("Failed to refresh Neustar Geolocation Service:" + e.getMessage());
		}
	}
}
