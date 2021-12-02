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

import org.apache.traffic_control.traffic_router.configuration.ConfigurationListener;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;

import javax.annotation.PostConstruct;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;

public class TrafficRouterConfigurationListener implements ConfigurationListener {
	private final Logger logger = LogManager.getLogger(TrafficRouterConfigurationListener.class);

	@Autowired
	private Environment environment;

	@Autowired
	ScheduledExecutorService scheduledExecutorService;

	@Autowired
	ServiceRefresher serviceRefresher;

	private ScheduledFuture<?> scheduledFuture;

	@Override
	public void configurationChanged() {
		boolean restarting = false;
		if (scheduledFuture != null) {
			restarting = true;
			scheduledFuture.cancel(true);

			while (!scheduledFuture.isDone()) {
				try {
					Thread.sleep(100L);
				} catch (InterruptedException e) {
					// ignore
				}
			}
		}

		Long fixedRate = environment.getProperty("neustar.polling.interval", Long.class, 86400000L);
		scheduledFuture = scheduledExecutorService.scheduleAtFixedRate(serviceRefresher, 0L, fixedRate, TimeUnit.MILLISECONDS);

		String prefix = restarting ? "Restarting" : "Starting";
		logger.warn(prefix + " Neustar remote database refresher at rate " + fixedRate + " msec");
	}
}
