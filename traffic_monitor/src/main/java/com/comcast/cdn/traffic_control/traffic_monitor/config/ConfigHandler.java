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
import java.io.FileNotFoundException;
import java.io.FileReader;

import org.apache.commons.io.IOUtils;
import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONObject;

public class ConfigHandler {
	private static final Logger LOGGER = Logger.getLogger(ConfigHandler.class);
	public static final String CONFIG_FILEPATH =   "/opt/traffic_monitor/conf/traffic_monitor_config.js";
	public static final String VAR_FILEPATH = "/opt/traffic_monitor/var";
	public static final String DB_FILEPATH = "/opt/traffic_monitor/db/";

	private final Object lok = new Object();
	private String confDir = null;
	private String confFile = null;
	private MonitorConfig config = null;
	private boolean shutdown;
	private final File configFile = new File(CONFIG_FILEPATH);
	private File varDirectory = new File(VAR_FILEPATH);

	// Recommended Singleton Pattern implementation
	// https://community.oracle.com/docs/DOC-918906

	private ConfigHandler() { }

	public static ConfigHandler getInstance() {
		return ConfigHandlerHolder.CONFIG_HANDLER;
	}

	private static class ConfigHandlerHolder {
		private static final ConfigHandler CONFIG_HANDLER = new ConfigHandler();
	}

	public void destroy() {
		shutdown = true;
	}

	public MonitorConfig getConfig() {
		if (shutdown) {
			return null;
		}

		synchronized (lok) {
			if (config != null) {
				return config;
			}

			final String confFile = getConfFile();

			if (confFile == null || confFile.isEmpty()) {
				config = new MonitorConfig();
				return config;
			}

			try {
				final String str = IOUtils.toString(new FileReader(getConfFile()));
				final JSONObject o = new JSONObject(str);
				config = new MonitorConfig(o.getJSONObject("traffic_monitor_config"));
			} catch (FileNotFoundException e) {
				LOGGER.error("Failed to find traffic monitor configuration file " + CONFIG_FILEPATH);
			} catch (Exception e) {
				LOGGER.warn(e, e);
			}

			if (config == null) {
				config = new MonitorConfig();
			}
		}

		return config;
	}

	public String getDbDir() {
		synchronized (lok) {
			if (confDir != null) {
				return confDir;
			}

			if (varDirectory.exists()) {
				confDir = DB_FILEPATH;
			}

			return confDir;
		}
	}

	public String getConfFile() {
		synchronized (lok) {
			if (confFile != null) {
				return confFile;
			}

			if (configFile.exists()) {
				confFile = CONFIG_FILEPATH;
			}

			return confFile;
		}
	}
}
