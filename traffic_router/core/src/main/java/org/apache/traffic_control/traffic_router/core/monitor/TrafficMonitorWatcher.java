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

package org.apache.traffic_control.traffic_router.core.monitor;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.net.URI;
import java.net.UnknownHostException;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Properties;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import org.apache.traffic_control.traffic_router.core.config.ConfigHandler;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouterManager;
import org.apache.traffic_control.traffic_router.core.util.AbstractUpdatable;
import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;
import org.apache.traffic_control.traffic_router.core.util.PeriodicResourceUpdater;
import org.springframework.context.ApplicationListener;
import org.springframework.context.event.ApplicationContextEvent;
import org.springframework.context.event.ContextClosedEvent;

@SuppressWarnings("PMD.TooManyFields")
public class TrafficMonitorWatcher implements ApplicationListener<ApplicationContextEvent> {
	private static final Logger LOGGER = LogManager.getLogger(TrafficMonitorWatcher.class);

	private String stateUrl;
	private String configUrl;
	private String monitorHosts;
	private int pollingInterval = 5000;
	private long lastHostAttempt = 0;
	private long reloadPeriod = 60*1000;

	private ConfigHandler configHandler;
	private TrafficRouterManager trafficRouterManager;

	private String statusFile;
	private int statusRefreshPeriod;
	private String configFile;
	private int configRefreshPeriod;

	private String monitorProperties;
	private static boolean bootstrapped = false;
	private static boolean localConfig = false;
	private static List<String> onlineMonitors = new ArrayList<String>();
	private static String[] hosts;
	private static Object hostSync = new Object();
	private static Object monitorSync = new Object();

	private PeriodicResourceUpdater crUpdater;
	private PeriodicResourceUpdater stateUpdater;
	private Path propertiesDirectory;
	private Path databasesDirectory;

	public AbstractUpdatable stateHandler = new AbstractUpdatable() {
		public String toString() {return "status listener";}
		@Override
		public synchronized boolean update(final String jsonStr) {
			try {
				final ObjectMapper mapper = new ObjectMapper();
				return trafficRouterManager.setState(mapper.readTree(jsonStr));
			} catch (JsonProcessingException e) {
				LOGGER.warn("problem with json: ",e);
			} catch (IOException e) {
				LOGGER.warn(e,e);
			}
			return false;
		}
		@Override
		public boolean noChange() {
			try {
				trafficRouterManager.setState(null);
			} catch (UnknownHostException e) {
				LOGGER.warn("UnknownHostException: ",e);
			}
			return false;
		}
	};

	public void destroy() {
		if (crUpdater != null) {
			crUpdater.destroy();
		}

		if (stateUpdater != null) {
			stateUpdater.destroy();
		}
	}

	@SuppressWarnings("PMD.CyclomaticComplexity")
	public void init() {
		final AbstractUpdatable crHandler = new AbstractUpdatable() {
			@Override
			public boolean update(final String configStr) {
				try {
					try {
						return configHandler.processConfig(configStr);
					} catch (JsonUtilsException e) {
						LOGGER.warn(e, e);
					}
				} catch (IOException e) {
					LOGGER.warn("error on config update", e);
				}

				return false;
			}
			public String toString() {return "config listener";}
			@Override
			public boolean noChange() {
				try {
					configHandler.processConfig(null);
				} catch (Exception e) {
					LOGGER.warn(e, e);
				}
				return false;
			}

			@Override
			public void complete() {
				if (!isLocalConfig() && !isBootstrapped()) {
					setBootstrapped(true);
				}
			}

			@Override
			public void cancelUpdate() {
				configHandler.cancelProcessConfig();
			}
		};

		processConfig();

		crUpdater = new PeriodicResourceUpdater(crHandler, new TrafficMonitorResourceUrl(this, configUrl), databasesDirectory.resolve(configFile).toString(), configRefreshPeriod, true);
		crUpdater.init();

		stateUpdater = new PeriodicResourceUpdater(stateHandler, new TrafficMonitorResourceUrl(this, stateUrl), databasesDirectory.resolve(statusFile).toString(), statusRefreshPeriod, true);
		stateUpdater.init();
	}

	@Override
	public void onApplicationEvent(final ApplicationContextEvent event) {
		if (event instanceof ContextClosedEvent) {
			crUpdater.destroy();
			stateUpdater.destroy();
		}
	}

	public String getStateUrl() {
		return stateUrl;
	}
	public void setStateUrl(final String stateUrl) {
		this.stateUrl = stateUrl;
	}
	public String getConfigUrl() {
		return configUrl;
	}
	public void setConfigUrl(final String configUrl) {
		this.configUrl = configUrl;
	}
	public void setPollingInterval(final int pollingInterval) {
		this.pollingInterval = pollingInterval;
	}
	public int getPollingInterval() {
		return pollingInterval;
	}

	public ConfigHandler getConfigHandler() {
		return configHandler;
	}
	public void setConfigHandler(final ConfigHandler configHandler) {
		this.configHandler = configHandler;
	}

	public String getStatusFile() {
		return statusFile;
	}
	public void setStatusFile(final String statusFile) {
		this.statusFile = statusFile;
	}
	public int getStatusRefreshPeriod() {
		return statusRefreshPeriod;
	}
	public void setStatusRefreshPeriod(final int statusRefreshPeriod) {
		this.statusRefreshPeriod = statusRefreshPeriod;
	}
	public String getConfigFile() {
		return configFile;
	}
	public void setConfigFile(final String configFile) {
		this.configFile = configFile;
	}
	public int getConfigRefreshPeriod() {
		return configRefreshPeriod;
	}
	public void setConfigRefreshPeriod(final int configRefreshPeriod) {
		this.configRefreshPeriod = configRefreshPeriod;
	}
	public TrafficRouterManager getTrafficRouterManager() {
		return trafficRouterManager;
	}
	public void setTrafficRouterManager(final TrafficRouterManager router) {
		this.trafficRouterManager = router;
	}
	public void setMonitorProperties(final String monitorProperties) {
		this.monitorProperties = monitorProperties;
	}
	public void setMonitorHosts(final String monitorHosts) {
		this.monitorHosts = monitorHosts;
	}

	public static void setHosts(final String[] newHosts) {
		synchronized(hostSync) {
			if (hosts == null || hosts.length == 0) {
				hosts = newHosts;
				LOGGER.info("traffic_monitor.bootstrap.hosts: " + Arrays.toString(hosts));
			} else if (!Arrays.asList(hosts).containsAll(Arrays.asList(newHosts))
					|| !Arrays.asList(newHosts).containsAll(Arrays.asList(hosts))) {
				hosts = newHosts;
				LOGGER.info("traffic_monitor.bootstrap.hosts changed to: " + Arrays.toString(hosts));
			}
		}
	}

	public String[] getHosts() {
		processConfig();

		return hosts;
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	private void processConfig() {
		final long now = System.currentTimeMillis();

		if (now < (lastHostAttempt+reloadPeriod)) {
			return;
		}

		lastHostAttempt = now;

		try {
			File trafficMonitorConfigFile;

			if (monitorProperties.matches("^\\w+:.*")) {
				trafficMonitorConfigFile = new File(new URI(monitorProperties));
			} else {
				trafficMonitorConfigFile = new File(monitorProperties);
			}

			final Properties props = new Properties();

			if (trafficMonitorConfigFile.exists()) {
				LOGGER.info("Loading properties from " + trafficMonitorConfigFile.getAbsolutePath());
				try (FileInputStream configStream = new FileInputStream(trafficMonitorConfigFile)) {
					props.load(configStream);
				}
			} else {
				LOGGER.warn("Cannot load traffic monitor properties file " + trafficMonitorConfigFile.getAbsolutePath() + " file not found!");
			}

			boolean localConfig = Boolean.parseBoolean(props.getProperty("traffic_monitor.bootstrap.local", "false"));

			String localEnvString = System.getenv("TRAFFIC_MONITOR_BOOTSTRAP_LOCAL");

			if (localEnvString != null) {
				localEnvString = localEnvString.toLowerCase();
			}

			if ("true".equals(localEnvString) || "false".equals(localEnvString)) {
				localConfig = Boolean.parseBoolean(localEnvString);
			}

			if (localConfig != isLocalConfig()) {
				LOGGER.info("traffic_monitor.bootstrap.local changed to: " + localConfig);
				setLocalConfig(localConfig);
			}

			if (localConfig || !isBootstrapped()) {
				String hostList = System.getenv("TRAFFIC_MONITOR_HOSTS");

				if (hostList != null && !hostList.isEmpty()) {
					LOGGER.warn("hostlist initialized to '" + hostList + "' from env var 'TRAFFIC_MONITOR_HOSTS'");
				}

				if (hostList == null || hostList.isEmpty()) {
					hostList = props.getProperty("traffic_monitor.bootstrap.hosts");
				}

				if (hostList == null || hostList.isEmpty()) {
					if (!trafficMonitorConfigFile.exists()) {
						LOGGER.fatal(trafficMonitorConfigFile.getAbsolutePath() + " does not exist and the environment variable 'TRAFFIC_MONITOR_HOSTS' was not found");
					} else {
						LOGGER.error("Cannot determine Traffic Monitor hosts from property 'traffic_monitor.bootstrap.hosts' in config file " + trafficMonitorConfigFile.getAbsolutePath());
					}
				} else {
					setHosts(hostList.contains(";") ? hostList.split(";") : new String[]{hostList});
				}

			} else if (!isLocalConfig() && isBootstrapped()) {
				synchronized(monitorSync) {
					if (!onlineMonitors.isEmpty()) {
						setHosts(onlineMonitors.toArray(new String[0]));
					}
				}
			}

			final String reloadPeriodStr = props.getProperty("traffic_monitor.properties.reload.period");

			if (reloadPeriodStr != null) {
				final long newReloadPeriod = Integer.parseInt(reloadPeriodStr);

				if (newReloadPeriod != reloadPeriod) {
					reloadPeriod = newReloadPeriod;
					LOGGER.info("traffic_monitor.properties.reload.period changed to: " + reloadPeriod);
				}
			}
		} catch (Exception e) {
			LOGGER.warn(e,e);
		}

		if (hosts==null) {
			hosts = monitorHosts.split(";");
		}
	}

	public static boolean isBootstrapped() {
		return bootstrapped;
	}

	private static void setBootstrapped(final boolean bootstrapped) {
		TrafficMonitorWatcher.bootstrapped = bootstrapped;
	}

	public static boolean isLocalConfig() {
		return localConfig;
	}

	private static void setLocalConfig(final boolean localConfig) {
		TrafficMonitorWatcher.localConfig = localConfig;
	}

	public static List<String> getOnlineMonitors() {
		return onlineMonitors;
	}

	public static void setOnlineMonitors(final List<String> onlineMonitors) {
		synchronized(monitorSync) {
			if (isLocalConfig()) {
				return;
			}

			TrafficMonitorWatcher.onlineMonitors = onlineMonitors;
			setBootstrapped(true);
			setHosts(onlineMonitors.toArray(new String[0]));
		}
	}

	public Path getPropertiesDirectory() {
		return propertiesDirectory;
	}

	public void setPropertiesDirectory(final Path propertiesDirectory) {
		this.propertiesDirectory = propertiesDirectory;
	}

	public Path getDatabasesDirectory() {
		return databasesDirectory;
	}

	public void setDatabasesDirectory(final Path databasesDirectory) {
		this.databasesDirectory = databasesDirectory;
	}
}
