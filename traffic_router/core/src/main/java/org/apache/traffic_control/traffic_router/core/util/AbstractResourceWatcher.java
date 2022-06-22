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

package org.apache.traffic_control.traffic_router.core.util;

import org.apache.traffic_control.traffic_router.core.config.WatcherConfig;
import org.apache.traffic_control.traffic_router.core.loc.AbstractServiceUpdater;
import com.fasterxml.jackson.databind.JsonNode;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.io.File;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.net.URL;

public abstract class AbstractResourceWatcher extends AbstractServiceUpdater {
	private static final Logger LOGGER = LogManager.getLogger(AbstractResourceWatcher.class);

	private URL authorizationUrl;
	private String postData;
	private ProtectedFetcher fetcher;
	protected TrafficOpsUtils trafficOpsUtils;
	private int timeout = 15000;

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	public void configure(final JsonNode config) {
		URL authUrl;
		String credentials;

		try {
			authUrl = new URL(trafficOpsUtils.getAuthUrl());
			credentials = trafficOpsUtils.getAuthJSON().toString();
		} catch (Exception e) {
			LOGGER.warn("Failed to update URL for TrafficOps authorization, " +
				"check the api.auth.url, and the TrafficOps username and password configuration setting: " + e.getMessage());
			// All or nothing, don't allow the watcher to be halfway misconfigured
			authUrl = this.authorizationUrl;
			credentials = this.postData;
		}

		if (authUrl == null || credentials == null) {
			LOGGER.warn("[ " + getClass().getSimpleName() + " ] Invalid Traffic Ops authorization URL or credentials data, not updating configuration!");
			return;
		}

		final WatcherConfig watcherConfig = new WatcherConfig(getWatcherConfigPrefix(), config, trafficOpsUtils);
		final String resourceUrl = (watcherConfig.getUrl() != null && !watcherConfig.getUrl().isEmpty()) ? watcherConfig.getUrl() : defaultDatabaseURL;

		final long pollingInterval = (watcherConfig.getInterval() > 0) ? watcherConfig.getInterval() : getPollingInterval();
		final int configTimeout = (watcherConfig.getTimeout() > 0) ? watcherConfig.getTimeout() : this.timeout;

		if (authUrl.equals(this.authorizationUrl) &&
			credentials.equals(this.postData) &&
			resourceUrl.equals(dataBaseURL) &&
			pollingInterval == getPollingInterval() &&
			configTimeout == this.timeout) {
			LOGGER.info("[ " + getClass().getName() + " ] Nothing changed in configuration");
			return;
		}

		// avoid recreating the fetcher if possible
		if (!authUrl.equals(this.authorizationUrl) || !credentials.equals(this.postData) || configTimeout != this.timeout) {
			this.authorizationUrl = authUrl;
			this.postData = credentials;
			this.timeout = configTimeout;
			fetcher = new ProtectedFetcher(authUrl.toString(), credentials, configTimeout);
		}

		setDataBaseURL(resourceUrl, pollingInterval);
	}

	protected boolean useData(final String data) {
		return true;
	}

	abstract protected boolean verifyData(final String data);

	@Override
	public boolean loadDatabase() throws IOException {
		final File existingDB = databasesDirectory.resolve(databaseName).toFile();

		if (!existingDB.exists() || !existingDB.canRead()) {
			return false;
		}

		final char[] jsonData = new char[(int) existingDB.length()];
		final FileReader reader = new FileReader(existingDB);

		try {
			reader.read(jsonData);
		} finally {
			reader.close();
		}

		return useData(new String(jsonData));
	}

	@Override
	public boolean verifyDatabase(final File dbFile) throws IOException {
		if (!dbFile.exists() || !dbFile.canRead()) {
			return false;
		}

		final char[] jsonData = new char[(int) dbFile.length()];
		final FileReader reader = new FileReader(dbFile);

		try {
			reader.read(jsonData);
		} finally {
			reader.close();
		}

		return verifyData(new String(jsonData));
	}


	@Override
	protected File downloadDatabase(final String url, final File existingDb) {
		if ((trafficOpsUtils.getHostname() == null) || trafficOpsUtils.getCdnName() == null) {
			return null;
		}
		final String interpolatedUrl = trafficOpsUtils.replaceTokens(url);
		if (fetcher == null) {
			LOGGER.warn("[" + getClass().getSimpleName() + "] Waiting for configuration to be processed, unable to download from '" + interpolatedUrl + "'");
			return null;
		}

		String jsonData = null;
		try {
			jsonData = fetcher.fetchIfModifiedSince(interpolatedUrl, existingDb.lastModified());
		}
		catch (IOException e) {
			LOGGER.warn("[ " + getClass().getSimpleName() + " ] Failed to fetch data from '" + interpolatedUrl + "': " + e.getMessage());
		}

		if (jsonData == null) {
			return existingDb;
		}

		File databaseFile = null;
		try {
			databaseFile = File.createTempFile(tmpPrefix, tmpSuffix);
			databaseFile.setReadable(true);
			databaseFile.setWritable(true);
			try (FileWriter fw = new FileWriter(databaseFile)) {
				fw.write(jsonData);
				fw.flush();
			}
		}
		catch (IOException e) {
			LOGGER.warn("Failed to create file from data received from '" + interpolatedUrl + "'");
		}

		return databaseFile;
	}

	public void setTrafficOpsUtils(final TrafficOpsUtils trafficOpsUtils) {
		this.trafficOpsUtils = trafficOpsUtils;
	}

	public abstract String getWatcherConfigPrefix();
}
