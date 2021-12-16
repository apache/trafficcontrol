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

package org.apache.traffic_control.traffic_router.core.router;

import java.io.IOException;
import java.net.UnknownHostException;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

import org.apache.traffic_control.traffic_router.core.ds.SteeringRegistry;
import org.apache.traffic_control.traffic_router.core.loc.FederationRegistry;
import com.fasterxml.jackson.databind.JsonNode;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import org.apache.traffic_control.traffic_router.core.edge.CacheRegister;
import org.apache.traffic_control.traffic_router.core.dns.NameServer;
import org.apache.traffic_control.traffic_router.geolocation.GeolocationService;
import org.apache.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationListener;
import org.springframework.context.event.ContextRefreshedEvent;
import org.apache.traffic_control.traffic_router.core.loc.AnonymousIpDatabaseService;

public class TrafficRouterManager implements ApplicationListener<ContextRefreshedEvent> {
	private static final Logger LOGGER = LogManager.getLogger(TrafficRouterManager.class);

	public static final int DEFAULT_API_PORT = 3333;
	public static final int DEFAULT_SECURE_API_PORT = 0; // Must be set through server.xml properties

	private JsonNode state;
	private TrafficRouter trafficRouter;
	private GeolocationService geolocationService;
	private GeolocationService geolocationService6;
	private AnonymousIpDatabaseService anonymousIpService;
	private StatTracker statTracker;
	private static final Map<String, Long> timeTracker = new ConcurrentHashMap<String, Long>();
	private NameServer nameServer;
	private TrafficOpsUtils trafficOpsUtils;
	private FederationRegistry federationRegistry;
	private SteeringRegistry steeringRegistry;
	private ApplicationContext applicationContext;
	private int apiPort = DEFAULT_API_PORT;
	private int secureApiPort = DEFAULT_SECURE_API_PORT;

	public NameServer getNameServer() {
		return nameServer;
	}

	public static Map<String, Long> getTimeTracker() {
		return timeTracker;
	}

	public void trackEvent(final String event) {
		timeTracker.put(event, System.currentTimeMillis());
	}

	public void setNameServer(final NameServer nameServer) {
		this.nameServer = nameServer;
	}

	public boolean setState(final JsonNode jsonObject) throws UnknownHostException {
		trackEvent("lastCacheStateCheck");

		if (jsonObject == null) {
			return false;
		}

		trackEvent("lastCacheStateChange");

		synchronized(this) {
			this.state = jsonObject;

			if (trafficRouter != null) {
				trafficRouter.setState(state);
			}

			return true;
		}
	}

	public TrafficRouter getTrafficRouter() {
		return trafficRouter;
	}

	public void setCacheRegister(final CacheRegister cacheRegister) throws IOException {
		trackEvent("lastConfigCheck");

		if (cacheRegister == null) {
			return;
		}

		final TrafficRouter tr = new TrafficRouter(cacheRegister, geolocationService, geolocationService6, anonymousIpService, statTracker, trafficOpsUtils, federationRegistry, this);
		tr.setSteeringRegistry(steeringRegistry);
		synchronized(this) {
			if (state != null) {
				try {
					tr.setState(state);
				} catch (UnknownHostException e) {
					LOGGER.warn(e,e);
				}
			}

			this.trafficRouter = tr;
			if (applicationContext != null) {
				this.trafficRouter.setApplicationContext(applicationContext);
			}
		}

		trackEvent("lastConfigChange");
	}

	public void setGeolocationService(final GeolocationService geolocationService) {
		this.geolocationService = geolocationService;
	}

	public void setGeolocationService6(final GeolocationService geolocationService) {
		this.geolocationService6 = geolocationService;
	}

	public void setAnonymousIpService(final AnonymousIpDatabaseService anonymousIpService) {
		this.anonymousIpService = anonymousIpService;
	}

	public void setStatTracker(final StatTracker statTracker) {
		this.statTracker = statTracker;
	}

	public void setTrafficOpsUtils(final TrafficOpsUtils trafficOpsUtils) {
		this.trafficOpsUtils = trafficOpsUtils;
	}

	public void setFederationRegistry(final FederationRegistry federationRegistry) {
		this.federationRegistry = federationRegistry;
	}

	public void setSteeringRegistry(final SteeringRegistry steeringRegistry) {
		this.steeringRegistry = steeringRegistry;
	}

	@Override
	public void onApplicationEvent(final ContextRefreshedEvent event) {
		applicationContext = event.getApplicationContext();
		if (trafficRouter != null) {
			trafficRouter.setApplicationContext(applicationContext);
			trafficRouter.configurationChanged();
		}
	}

	public void setApiPort(final int apiPort) {
		this.apiPort = apiPort;
	}

	public int getApiPort() {
		return apiPort;
	}

	public int getSecureApiPort() {
		return secureApiPort;
	}

	public void setSecureApiPort(final int secureApiPort) {
		this.secureApiPort = secureApiPort;
	}
}
