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

package com.comcast.cdn.traffic_control.traffic_router.core.util;

import java.net.InetAddress;
import java.net.UnknownHostException;

import org.apache.log4j.Logger;
import org.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;

public class LanguidState {
	private static final Logger LOGGER = Logger.getLogger(LanguidState.class);
	private boolean ready = false;
	private TrafficRouterManager trafficRouterManager;
	private int port = 0;
	private int apiPort = 0;
	private int securePort = 0;

	public void init() {
		if (trafficRouterManager == null || trafficRouterManager.getTrafficRouter() == null) {
			return;
		}

		final TrafficRouter tr = trafficRouterManager.getTrafficRouter();

		if (tr.getCacheRegister() == null) {
			return;
		}

		final String hostname;

		try {
			hostname = InetAddress.getLocalHost().getHostName().replaceAll("\\..*", "");
		} catch (UnknownHostException e) {
			LOGGER.error("Cannot lookup hostname of this traffic router!: " + e.getMessage());
			return;
		}

		final JSONObject routers = tr.getCacheRegister().getTrafficRouters();

		for (final String key : JSONObject.getNames(routers)) {
			final JSONObject routerJson = routers.optJSONObject(key);

			if (! hostname.equalsIgnoreCase(key)) {
				continue;
			}

			initPorts(routerJson);
			break;
		}

		setReady(true);
	}

	private void initPorts(final JSONObject routerJson) {
		if (routerJson.has("port")) {
			setPort(routerJson.optInt("port"));
		}

		if (routerJson.has("api.port")) {
			setApiPort(routerJson.optInt("api.port"));
			trafficRouterManager.setApiPort(apiPort);
		}

		if (routerJson.has("secure.port")) {
			setSecurePort(routerJson.optInt("secure.port"));
		}
	}

	public boolean isReady() {
		return ready;
	}

	public void setReady(final boolean ready) {
		this.ready = ready;
	}

	public int getPort() {
		return port;
	}

	public void setPort(final int port) {
		this.port = port;
	}

	public int getApiPort() {
		return apiPort;
	}

	public void setApiPort(final int apiPort) {
		this.apiPort = apiPort;
	}

	public TrafficRouterManager getTrafficRouterManager() {
		return trafficRouterManager;
	}

	public void setTrafficRouterManager(final TrafficRouterManager trafficRouterManager) {
		this.trafficRouterManager = trafficRouterManager;
	}

	public int getSecurePort() {
		return securePort;
	}

	public void setSecurePort(final int securePort) {
		this.securePort = securePort;
	}
}
