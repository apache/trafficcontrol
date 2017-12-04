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

package com.comcast.cdn.traffic_control.traffic_monitor.config;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileReader;

import org.apache.commons.io.IOUtils;
import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONObject;

public class ConfigHandler {
	private static final Logger LOGGER = Logger.getLogger(ConfigHandler.class);
	private static final String CONFIG_FILEPATH = "/opt/traffic_monitor/conf/traffic_monitor_config.js";
	private static final String DB_FILEPATH = "/opt/traffic_monitor/db";
	public static final String CONFIG_PROPERTY = "traffic_monitor.path.config";
	public static final String DB_PROPERTY = "traffic_monitor.path.db";

	private final Object lok = new Object();
	private MonitorConfig config = null;
	private boolean shutdown;
	private final File configFile = new File(getFilePath(CONFIG_PROPERTY, CONFIG_FILEPATH));
	private final String dbPath = getFilePath(DB_PROPERTY, DB_FILEPATH);

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

			if (!configFileExists()) {
				config = new MonitorConfig();
				return config;
			}

			try {
				final JSONObject jsonConfig = new JSONObject(IOUtils.toString(new FileReader(configFile)));
				config = new MonitorConfig(jsonConfig.getJSONObject("traffic_monitor_config"));
			} catch (FileNotFoundException e) {
				LOGGER.error("Failed to find traffic monitor configuration file " + configFile.toString());
			} catch (Exception e) {
				LOGGER.warn(e, e);
			}

			if (config == null) {
				config = new MonitorConfig();
			}
		}

		return config;
	}

	public File getDbFile(final String filename) {
		return new File(dbPath, filename);
	}

	public File getConfigFile() {
		return configFile;
	}

	public boolean configFileExists() {
		return configFile.exists();
	}

	private String getFilePath(final String property, final String staticFilePath) {
		if (property != null && System.getProperty(property) != null) {
			return System.getProperty(property);
		}

		return staticFilePath;
	}
}
