/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.config;

import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import org.apache.log4j.Logger;
import org.json.JSONObject;

import java.net.URL;

public class WatcherConfig {
	private static final Logger LOGGER = Logger.getLogger(WatcherConfig.class);
	private final URL url;
	private final long interval;
	private final int timeout;

	public WatcherConfig(final String prefix, final JSONObject config, final TrafficOpsUtils trafficOpsUtils) {
		// Make PMD happy by using temporary variable :(
		URL u = null;
		try{
			u = new URL(trafficOpsUtils.getUrl(prefix + ".polling.url"));
		} catch (Exception e) {
			LOGGER.warn("Invalid Federation Polling URL, check the " + prefix + ".polling.url configuration: " + e.getMessage());
		}

		url = u;

		interval = config.optLong(prefix + ".polling.interval", -1L);
		if (interval == -1L) {
			LOGGER.warn("Failed getting valid value for " + prefix + ".polling.interval");
		}

		timeout = config.optInt(prefix + ".polling.timeout", -1);
		if (timeout == -1) {
			LOGGER.warn("Failed getting valid value for " + prefix + ".polling.timeout");
		}
	}

	public long getInterval() {
		return interval;
	}

	public URL getUrl() {
		return url;
	}

	public int getTimeout() {
		return timeout;
	}
}
