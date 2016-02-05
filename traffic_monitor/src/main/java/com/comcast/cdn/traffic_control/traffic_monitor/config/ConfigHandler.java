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

package com.comcast.cdn.traffic_control.traffic_monitor.config;

import java.io.File;
import java.io.FileReader;

import org.apache.commons.io.IOUtils;
import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONObject;

public class ConfigHandler {
	private static final Logger LOGGER = Logger.getLogger(ConfigHandler.class);

	private static final Object lok = new Object();
	private static String confDir = null;
	private static String confFile = null;
	private static MonitorConfig config = null;
	private static boolean shutdown;

	public static void destroy() {
		shutdown = true;
	}

	public static MonitorConfig getConfig() {
		if (shutdown) {
			return null;
		}

		synchronized (lok) {
			if (config != null) {
				return config;
			}

			try {
				final String str = IOUtils.toString(new FileReader(getConfFile()));
				final JSONObject o = new JSONObject(str);
				config = new MonitorConfig(o.getJSONObject("traffic_monitor_config"));
			} catch (Exception e) {
				LOGGER.warn(e, e);
			}

			if (config == null) {
				config = new MonitorConfig();
			}
		}

		return config;
	}

	public static String getDbDir() {
		synchronized (lok) {
			if (confDir != null) {
				return confDir;
			}

			confDir = "target/test-classes/var/";

			if (new File("/opt/traffic_monitor/var").exists()) {
				confDir = "/opt/traffic_monitor/db/";
			}

			return confDir;
		}
	}

	public static String getConfFile() {
		synchronized (lok) {
			if (confFile != null) {
				return confFile;
			}

			confFile = "target/test-classes/conf/traffic_monitor_config.js";

			if (new File("/opt/traffic_monitor/conf/traffic_monitor_config.js").exists()) {
				confFile = "/opt/traffic_monitor/conf/traffic_monitor_config.js";
			}

			return confFile;
		}
	}
}
