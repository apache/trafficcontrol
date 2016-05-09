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

package com.comcast.cdn.traffic_control.traffic_router.core.util;

import com.comcast.cdn.traffic_control.traffic_router.core.config.WatcherConfig;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.AbstractServiceUpdater;
import org.apache.log4j.Logger;
import org.json.JSONObject;

import java.io.File;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.net.URL;

public abstract class AbstractResourceWatcher extends AbstractServiceUpdater {
	private static final Logger LOGGER = Logger.getLogger(AbstractResourceWatcher.class);

	private URL authorizationUrl;
	private String postData;
	private ProtectedFetcher fetcher;
	protected TrafficOpsUtils trafficOpsUtils;

	public void configure(final JSONObject config) {
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

		final WatcherConfig watcherConfig = new WatcherConfig(getWatcherConfigPrefix(), config, trafficOpsUtils);

		int timeout =  watcherConfig.getTimeout();
		if (timeout == -1) {
			timeout = 15 * 1000;
		}

		if (authUrl != null && credentials != null && watcherConfig.getUrl() != null && watcherConfig.getInterval() != -1L) {
			configure(authUrl, credentials, watcherConfig.getUrl(), watcherConfig.getInterval(), timeout);
		} else {
			LOGGER.warn("Not updating configuration - did get following from cr-config url '" + watcherConfig.getInterval() + "' interval '" + watcherConfig.getInterval() + "' timeout '" + watcherConfig.getTimeout() + "'");
		}
	}


	public void configure(final URL authorizationUrl, final String postData, final URL resourceUrl, final long pollingInterval, final int timeout) {
		if (authorizationUrl.equals(this.authorizationUrl) && postData.equals(this.postData) &&
			resourceUrl.toString().equals(dataBaseURL) && pollingInterval == getPollingInterval()) {
			LOGGER.warn("Nothing changed in configuration");
			return;
		}

		// avoid recreating the fetcher if possible
		if (!authorizationUrl.equals(this.authorizationUrl) || !postData.equals(this.postData)) {
			this.authorizationUrl = authorizationUrl;
			this.postData = postData;
			fetcher = new ProtectedFetcher(authorizationUrl.toString(), postData, timeout);
		}

		setDataBaseURL(resourceUrl.toString(), pollingInterval);
	}

	protected boolean useData(final String data) {
		return true;
	}

	@Override
	public boolean loadDatabase() throws IOException, org.apache.wicket.ajax.json.JSONException {
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
	protected File downloadDatabase(final String url, final File existingDb) {
		if (fetcher == null) {
			LOGGER.warn("Waiting for federations configuration to be processed, unable download federations");
			return null;
		}

		String jsonData = null;
		try {
			jsonData = fetcher.fetchIfModifiedSince(url, existingDb.lastModified());
		}
		catch (IOException e) {
			LOGGER.warn("Failed to fetch data from '" + url + "': " + e.getMessage());
		}

		if (jsonData == null) {
			return existingDb;
		}

		File databaseFile = null;
		FileWriter fw;
		try {
			databaseFile = File.createTempFile(tmpPrefix, tmpSuffix);
			fw = new FileWriter(databaseFile);
			fw.write(jsonData);
			fw.flush();
			fw.close();
		}
		catch (IOException e) {
			LOGGER.warn("Failed to create file from data received from '" + url + "'");
		}

		return databaseFile;
	}

	public void setTrafficOpsUtils(final TrafficOpsUtils trafficOpsUtils) {
		this.trafficOpsUtils = trafficOpsUtils;
	}

	public abstract String getWatcherConfigPrefix();
}
