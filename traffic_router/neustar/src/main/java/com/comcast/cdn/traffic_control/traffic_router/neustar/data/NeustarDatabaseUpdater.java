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
import java.nio.file.Files;
import java.nio.file.Path;
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
	private File neustarOldDatabaseDirectory;

	@Autowired
	private FilesMover filesMover;

	private HttpClient httpClient = new HttpClient();

	@Autowired
	private TarExtractor tarExtractor;

	public void setHttpClient(HttpClient httpClient) {
		this.httpClient = httpClient;
	}

	private File createTmpDir(File directory) throws IOException {
		return Files.createTempDirectory(directory.toPath(), null).toFile();
	}

	public File extractRemoteContent(InputStream inputStream) throws IOException {
		return tarExtractor.extractTgzTo(createTmpDir(neustarDatabaseDirectory), inputStream);
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

		File tmpDir = null;
		try {

			tmpDir = createTmpDir(neustarDatabaseDirectory);
			tarExtractor.extractTo(tmpDir, new GZIPInputStream(response.getEntity().getContent()));
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

		if (!verifyNewDatabase(tmpDir)) {
			filesMover.purgeDirectory(tmpDir);
			return false;
		}

		LOGGER.info("Replacing files in " + neustarDatabaseDirectory.getAbsolutePath() + " with those in " + tmpDir.getAbsolutePath());
		if (!filesMover.updateCurrent(neustarDatabaseDirectory, tmpDir, neustarOldDatabaseDirectory)) {
			LOGGER.warn("Failed replacing files, not purging " + tmpDir.getAbsolutePath());
			return false;
		}

		if (filesMover.purgeDirectory(tmpDir)) {
			tmpDir.delete();
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
