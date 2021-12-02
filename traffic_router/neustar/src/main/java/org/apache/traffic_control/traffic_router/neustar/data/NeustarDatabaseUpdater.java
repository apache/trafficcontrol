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

package org.apache.traffic_control.traffic_router.neustar.data;

import org.apache.traffic_control.traffic_router.neustar.files.FilesMover;
import com.quova.bff.reader.io.GPDatabaseReader;
import org.apache.http.Header;
import org.apache.http.HttpResponse;
import org.apache.http.client.config.RequestConfig;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;

import java.io.File;
import java.io.IOException;
import java.net.URI;
import java.nio.file.Files;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.zip.GZIPInputStream;

import static java.lang.Long.parseLong;

public class NeustarDatabaseUpdater {
	private final Logger LOGGER = LogManager.getLogger(NeustarDatabaseUpdater.class);

	@Autowired
	private Integer neustarPollingTimeout;

	@Autowired
	private String neustarDataUrl;

	@Autowired
	private File neustarDatabaseDirectory;

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

	private File createTmpDir(File directory) {
		try {
			return Files.createTempDirectory(directory.toPath(), "neustar-").toFile();
		} catch (IOException e) {
			System.out.println("Failed to create temporary directory in " + directory.getAbsolutePath() + ": " + e.getMessage());
		}

		return null;
	}

	public boolean verifyNewDatabase(File directory) {
		try {
			new GPDatabaseReader.Builder(directory).build();
			return true;
		} catch (Exception e) {
			LOGGER.error("Database Directory " + directory + " is not a valid Neustar database. " + e.getMessage());
			return false;
		}
	}

	public boolean update() {
		if (neustarDataUrl == null || neustarDataUrl.isEmpty()) {
			LOGGER.error("Cannot get latest neustar data 'neustar.polling.url' needs to be set in environment or properties file");
			return false;
		}

		File tmpDir = createTmpDir(neustarDatabaseDirectory);
		if (tmpDir == null) {
			return false;
		}

		try (CloseableHttpResponse response = getRemoteDataResponse(URI.create(neustarDataUrl))) {
			if (response.getStatusLine().getStatusCode() == 304) {
				LOGGER.info("Neustar database unchanged at " + neustarDataUrl);
				return false;
			}

			if (response.getStatusLine().getStatusCode() != 200) {
				LOGGER.error("Failed downloading remote neustar database from " + neustarDataUrl + " " + response.getStatusLine().getReasonPhrase());
			}

			if (!enoughFreeSpace(tmpDir, response, neustarDataUrl)) {
				return false;
			}

			try (GZIPInputStream gzipStream = new GZIPInputStream(response.getEntity().getContent())) {
				if (!tarExtractor.extractTo(tmpDir, gzipStream)) {
					LOGGER.error("Failed to decompress remote content from " + neustarDataUrl);
					return false;
				}
			}

			LOGGER.info("Replacing neustar files in " + neustarDatabaseDirectory.getAbsolutePath() + " with those in " + tmpDir.getAbsolutePath());

			if (!filesMover.updateCurrent(neustarDatabaseDirectory, tmpDir, neustarOldDatabaseDirectory)) {
				LOGGER.error("Failed updating neustar files");
				return false;
			}

			if (!verifyNewDatabase(tmpDir)) {
				return false;
			}
		} catch (Exception e) {
			LOGGER.error("Failed getting remote neustar data: " + e.getMessage());
		} finally {
			httpClient.close();

			if (!filesMover.purgeDirectory(tmpDir) || !tmpDir.delete()) {
				LOGGER.error("Failed purging temporary directory " + tmpDir.getAbsolutePath());
			}
		}

		return true;
	}

	public CloseableHttpResponse getRemoteDataResponse(URI uri) {
		HttpGet httpGet = new HttpGet(uri);
		httpGet.setConfig(RequestConfig.custom().setSocketTimeout(neustarPollingTimeout).build());
		Date buildDate = getDatabaseBuildDate();

		if (buildDate != null) {
			httpGet.setHeader("If-Modified-Since", new SimpleDateFormat("EEE, dd MMM yyyy HH:mm:ss Z").format(buildDate));
		}

		return httpClient.execute(httpGet);
	}

	public boolean enoughFreeSpace(File destination, HttpResponse response, String request) {
		Header contentLengthHeader = response.getFirstHeader("Content-Length");

		if (contentLengthHeader == null) {
			LOGGER.warn("Unable to determine size of data from " + request);
			return true;
		}

		long contentLength = parseLong(contentLengthHeader.getValue());
		long freespace = destination.getFreeSpace();

		if (freespace < contentLength) {
			LOGGER.error("Not enough space in " + destination + " to save " + request + "(Free: " + freespace + ", Need: " + contentLength);
			return false;
		}

		return true;
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
