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

package com.comcast.cdn.traffic_control.traffic_router.neustar.data;

import com.comcast.cdn.traffic_control.traffic_router.neustar.files.FilesMover;
import com.quova.bff.reader.io.GPDatabaseReader;
import org.apache.http.client.config.RequestConfig;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.net.URI;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.zip.GZIPInputStream;

public class NeustarDatabaseUpdater {
	private final Logger LOGGER = Logger.getLogger(NeustarDatabaseUpdater.class);

	@Autowired
	private Integer neustarPollingTimeout;

	@Autowired
	private String neustarDataUrl;

	@Autowired
	private File neustarDatabaseDirectory;

	@Autowired
	private File neustarTempDatabaseDirectory;

	@Autowired
	private File neustarOldDatabaseDirectory;

	@Autowired
	private FilesMover filesMover;

	private HttpClient httpClient = new HttpClient();

	@Autowired
	private TarExtractor tarExtractor;

	public void setHttpClient(HttpClient httpClient) {
		this.httpClient = httpClient;
	}

	public File extractRemoteContent(InputStream inputStream) {
		if (!neustarTempDatabaseDirectory.exists() && !neustarTempDatabaseDirectory.mkdirs()) {
			LOGGER.error("Cannot save remote content from " + neustarDataUrl + " to disk: " + neustarTempDatabaseDirectory.getAbsolutePath() + " does not exist and cannot be created");
			return null;
		}

		return tarExtractor.extractTgzTo(neustarTempDatabaseDirectory, inputStream);
	}

	public boolean verifyNewDatabase() {
		try {
			new GPDatabaseReader.Builder(neustarTempDatabaseDirectory).build();
			return true;
		} catch (Exception e) {
			LOGGER.error("Database Directory " + neustarTempDatabaseDirectory + " is not a valid Neustar database. " + e.getMessage());
			return false;
		}
	}

	public boolean update() {
		if (neustarDataUrl == null || neustarDataUrl.isEmpty()) {
			LOGGER.error("Cannot get latest neustar data 'neustar.polling.url' needs to be set in environment or properties file");
			return false;
		}

		URI uri;
		try {
			uri = URI.create(neustarDataUrl);
		} catch (Exception e) {
			LOGGER.error("Cannot get latest neustar data 'neustar.polling.url' value '" + neustarDataUrl + "' is not valid");
			return false;
		}

		HttpGet httpGet = new HttpGet(uri);
		httpGet.setConfig(RequestConfig.custom().setSocketTimeout(neustarPollingTimeout).build());
		Date buildDate = getDatabaseBuildDate();

		if (buildDate != null) {
			httpGet.setHeader("If-Modified-Since", new SimpleDateFormat("EEE, dd MMM yyyy HH:mm:ss Z").format(buildDate));
		}

		CloseableHttpResponse response = httpClient.execute(httpGet);

		if (response == null) {
			return false;
		}

		if (response.getStatusLine().getStatusCode() != 200) {
			if (response.getStatusLine().getStatusCode() != 304) {
				LOGGER.warn("Failed downloading Neustar Database from " + neustarDataUrl + " " + response.getStatusLine().getReasonPhrase());
			}

			try {
				response.close();
			} catch (IOException e) {
				LOGGER.warn("Failed to close http response for " + neustarDataUrl + " : " + e.getMessage());
			}
			httpClient.close();
			return false;
		}

		try {
			tarExtractor.extractTo(neustarTempDatabaseDirectory, new GZIPInputStream(response.getEntity().getContent()));
		} catch (IOException e) {
			LOGGER.error("Failed to decompress remote content from " + neustarDataUrl + " : " + e.getMessage());
			return false;
		} finally {
			try {
				response.close();
			} catch (IOException e) {
				LOGGER.warn("Failed to close http response for " + neustarDataUrl);
			}
			httpClient.close();
		}

		if (!verifyNewDatabase()) {
			filesMover.purgeDirectory(neustarTempDatabaseDirectory);
			return false;
		}

		return filesMover.updateCurrent(neustarDatabaseDirectory, neustarTempDatabaseDirectory, neustarOldDatabaseDirectory);
	}

	public Date getDatabaseBuildDate() {
		File[] neustarDatabaseFiles = neustarDatabaseDirectory.listFiles();

		if (neustarDatabaseFiles == null || neustarDatabaseFiles.length == 0) {
			return null;
		}

		long modifiedTimestamp = 0;

		for (File file : neustarDatabaseFiles) {
			if (file.isDirectory()) {
				continue;
			}

			if (modifiedTimestamp == 0 || file.lastModified() < modifiedTimestamp) {
				modifiedTimestamp = file.lastModified();
			}
		}

		return modifiedTimestamp > 0 ? new Date(modifiedTimestamp) : null;
	}

	public void setNeustarDataUrl(String neustarDataUrl) {
		this.neustarDataUrl = neustarDataUrl;
	}

	public void setNeustarPollingTimeout(int neustarPollingTimeout) {
		this.neustarPollingTimeout = neustarPollingTimeout;
	}
}
