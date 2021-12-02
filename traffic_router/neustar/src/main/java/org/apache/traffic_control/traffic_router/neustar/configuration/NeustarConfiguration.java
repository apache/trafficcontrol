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
import org.apache.traffic_control.traffic_router.neustar.data.TarExtractor;
import org.apache.traffic_control.traffic_router.neustar.files.FilesMover;

import org.apache.traffic_control.traffic_router.configuration.ConfigurationListener;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;
import org.springframework.core.env.Environment;

import java.io.File;
import java.nio.file.Path;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;

@Configuration
@PropertySource(value = "neustar.properties", ignoreResourceNotFound = true)
public class NeustarConfiguration {
	private final Logger LOGGER = LogManager.getLogger(NeustarConfiguration.class);
	@Autowired
	private Environment environment;

	@Autowired
	private Path databasesDir;

	NeustarDatabaseUpdater neustarDatabaseUpdater;
	NeustarGeolocationService neustarGeolocationService;

	private File neustarDatabaseDirectory;
	private File neustarOldDatabaseDirectory;

	private File checkDirectory(File directory) {
		if (!directory.exists() && !directory.mkdirs()) {
			LOGGER.error(directory.getAbsolutePath() + " does not exist and cannot be created");
		}
		return directory;
	}

	@Bean
	public File neustarDatabaseDirectory() {
		if (neustarDatabaseDirectory == null) {
			neustarDatabaseDirectory = checkDirectory(databasesDir.resolve(environment.getProperty("neustar.subdirectory", "neustar")).toFile());
		}

		return neustarDatabaseDirectory;
	}

	@Bean
	public File neustarOldDatabaseDirectory() {
		if (neustarOldDatabaseDirectory == null) {
			neustarOldDatabaseDirectory = checkDirectory(new File(neustarDatabaseDirectory(), "/old"));
		}

		return neustarOldDatabaseDirectory;
	}

	@Bean
	public FilesMover filesMover() {
		return new FilesMover();
	}

	@Bean
	public TarExtractor tarExtractor() {
		return new TarExtractor();
	}

	@Bean
	public String neustarRemoteSource() {
		String pollingUri = environment.getProperty("neustar.polling.url");
		if (pollingUri == null || pollingUri.isEmpty()) {
			LOGGER.error("'neustar.polling.url' must be set in the environment or file 'neustar.properties' on the classpath");
		}

		LOGGER.info("Using " + pollingUri + " for 'neustar.polling.url'");
		return pollingUri;
	}

	@Bean
	public Integer neustarPollingTimeout() {
		return environment.getProperty("neustar.polling.timeout", Integer.class, 30000);
	}

	@Bean
	public NeustarDatabaseUpdater neustarDatabaseUpdater() {
		if (neustarDatabaseUpdater == null) {
			neustarDatabaseUpdater = new NeustarDatabaseUpdater();
		}
		return neustarDatabaseUpdater;
	}

	@Bean
	public NeustarGeolocationService neustarGeolocationService() {
		if (neustarGeolocationService == null) {
			neustarGeolocationService = new NeustarGeolocationService();
		}
		return neustarGeolocationService;
	}

	@Bean
	ServiceRefresher serviceRefresher() {
		return new ServiceRefresher();
	}

	@Bean
	ScheduledExecutorService scheduledExecutorService() {
		return Executors.newSingleThreadScheduledExecutor();
	}

	@Bean
	ConfigurationListener trafficRouterConfigurationListener() {
		return new TrafficRouterConfigurationListener();
	}
}
