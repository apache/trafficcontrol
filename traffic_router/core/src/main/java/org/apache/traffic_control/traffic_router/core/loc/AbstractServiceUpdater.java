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

package org.apache.traffic_control.traffic_router.core.loc;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.net.URLConnection;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;
import java.util.Arrays;
import java.util.Date;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.zip.GZIPInputStream;

import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;
import org.apache.commons.io.IOUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import org.apache.traffic_control.traffic_router.core.router.TrafficRouterManager;

import static org.apache.commons.codec.digest.DigestUtils.md5Hex;

@SuppressWarnings({"PMD.CyclomaticComplexity"})
public abstract class AbstractServiceUpdater {
	private static final Logger LOGGER = LogManager.getLogger(AbstractServiceUpdater.class);

	protected String dataBaseURL;
	protected String defaultDatabaseURL;
	protected String databaseName;
	protected ScheduledExecutorService executorService;
	private long pollingInterval;
	protected boolean loaded = false;
	protected ScheduledFuture<?> scheduledService;
	private TrafficRouterManager trafficRouterManager;
	protected Path databasesDirectory;
	private String eTag = null;

	public void destroy() {
		executorService.shutdownNow();
	}

	/**
	 * Gets dataBaseURL.
	 *
	 * @return the dataBaseURL
	 */
	public String getDataBaseURL() {
		return dataBaseURL;
	}

	/**
	 * Gets pollingInterval.
	 *
	 * @return the pollingInterval
	 */
	public long getPollingInterval() {
		if(pollingInterval == 0) { return 10000; }
		return pollingInterval;
	}

	final private Runnable updater = new Runnable() {
		@Override
		@SuppressWarnings("PMD.AvoidCatchingThrowable")
		public void run() {
			try {
				updateDatabase();
			} catch (Throwable t) {
				// Catching Throwable prevents this Service Updater thread from silently dying
				LOGGER.error( "[" + getClass().getSimpleName() +"] Failed updating database!", t);
			}
		}
	};

	public void init() {
		final long pollingInterval = getPollingInterval();
		final Date nextFetchDate = new Date(System.currentTimeMillis() + pollingInterval);
		LOGGER.info("[" + getClass().getSimpleName() + "] Fetching external resource " + dataBaseURL + " at interval: " + pollingInterval + " : " + TimeUnit.MILLISECONDS + " next update occurrs at " + nextFetchDate);
		scheduledService = executorService.scheduleWithFixedDelay(updater, pollingInterval, pollingInterval, TimeUnit.MILLISECONDS);
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	public boolean updateDatabase() {
		try {
			if (!Files.exists(databasesDirectory)) {
				Files.createDirectories(databasesDirectory);
			}

		} catch (IOException ex) {
			LOGGER.error(databasesDirectory.toString() + " does not exist and cannot be created!");
			return false;
		}

		final File existingDB = databasesDirectory.resolve(databaseName).toFile();

		if (!isLoaded()) {
			try {
				setLoaded(loadDatabase());
			} catch (Exception e) {
				LOGGER.warn("[" + getClass().getSimpleName() + "] Failed to load existing database! " + e.getMessage());
			}
		} else if (!needsUpdating(existingDB)) {
			LOGGER.info("[" + getClass().getSimpleName() + "] Location database does not require updating.");
			return false;
		}

		File newDB = null;
		boolean isModified = true;

		final String databaseURL = getDataBaseURL();
		if (databaseURL == null) {
			LOGGER.warn("[" + getClass().getSimpleName() + "] Skipping download/update: database URL is null");
			return false;
		}

		try {
			try {
				newDB = downloadDatabase(databaseURL, existingDB);
				trafficRouterManager.trackEvent("last" + getClass().getSimpleName() + "Check");

				// if the remote db's timestamp is less than or equal to ours, the above returns existingDB
				if (newDB == existingDB) {
					isModified = false;
				}
			} catch (Exception e) {
				LOGGER.fatal("[" + getClass().getSimpleName() + "] Caught exception while attempting to download: " + getDataBaseURL(), e);
				return false;
			}

			if (!isModified || newDB == null || !newDB.exists()) {
				return false;
			}

			try {
				if (!verifyDatabase(newDB)) {
					LOGGER.warn("[" + getClass().getSimpleName() + "] " + newDB.getAbsolutePath() + " from " + getDataBaseURL() + " is invalid!");
					return false;
				}
			} catch (Exception e) {
				LOGGER.error("[" + getClass().getSimpleName() + "] Failed verifying database " + newDB.getAbsolutePath() + " : " + e.getMessage());
				return false;
			}

			try {
				if (copyDatabaseIfDifferent(existingDB, newDB)) {
					setLoaded(loadDatabase());
					trafficRouterManager.trackEvent("last" + getClass().getSimpleName() + "Update");
				} else {
					newDB.delete();
				}
			} catch (Exception e) {
				LOGGER.error("[" + getClass().getSimpleName() + "] Failed copying and loading new database " + newDB.getAbsolutePath() + " : " + e.getMessage());
			}

		} finally {
			if (newDB != null && newDB != existingDB && newDB.exists()) {
				LOGGER.info("[" + getClass().getSimpleName() + "] Try to delete downloaded temp file");
				deleteDatabase(newDB);
			}
		}

		return true;
	}

	abstract public boolean verifyDatabase(final File dbFile) throws IOException, JsonUtilsException;
	abstract public boolean loadDatabase() throws IOException, JsonUtilsException;

	public void setDatabaseName(final String databaseName) {
		this.databaseName = databaseName;
	}

	public void stopServiceUpdater() {
		if (scheduledService != null) {
			LOGGER.info("[" + getClass().getSimpleName() + "] Stopping service updater");
			scheduledService.cancel(false);
		}
	}

	public void cancelServiceUpdater() {
		this.stopServiceUpdater();
		pollingInterval = 0;
		dataBaseURL = null;
	}

	public void setDataBaseURL(final String url, final long refresh) {
		if (refresh !=0 && refresh != pollingInterval) {

			this.pollingInterval = refresh;
			LOGGER.info("[" + getClass().getSimpleName() + "] Restarting schedule for " + url + " with interval: "+refresh);
			stopServiceUpdater();
			init();
		}

		if ((url != null) && !url.equals(dataBaseURL) || (refresh!=0 && refresh!=pollingInterval)) {
			this.dataBaseURL = url;
			this.setLoaded(false);
			new Thread(updater).start();
		}
	}

	public void setDatabaseUrl(final String url) {
		this.dataBaseURL = url;
	}

	public void setDefaultDatabaseUrl(final String url) {
		this.defaultDatabaseURL = url;
	}

	/**
	 * Sets executorService.
	 *
	 * @param executorService
	 *            the executorService to set
	 */
	public void setExecutorService(final ScheduledExecutorService executorService) {
		this.executorService = executorService;
	}

	/**
	 * Sets pollingInterval.
	 *
	 * @param pollingInterval
	 *            the pollingInterval to set
	 */
	public void setPollingInterval(final long pollingInterval) {
		this.pollingInterval = pollingInterval;
	}

	boolean filesEqual(final File a, final File b) throws IOException {
		if(!a.exists() && !b.exists()) { return true; }
		if(!a.exists() || !b.exists()) { return false; }
		if (a.isDirectory() && b.isDirectory()) {
			return compareDirectories(a, b);
		}
		return compareFiles(a, b);
	}

	private boolean compareDirectories(final File a, final File b) throws IOException {
		final File[] aFileList = a.listFiles();
		final File[] bFileList = b.listFiles();

		if (aFileList.length != bFileList.length) {
			return false;
		}

		Arrays.sort(aFileList);
		Arrays.sort(bFileList);

		for (int i = 0; i < aFileList.length; i++) {
			if (aFileList[i].length() != bFileList[i].length()) {
				return false;
			}
		}

		return true;
	}

	private String fileMd5(final File file) throws IOException {
		try (FileInputStream stream = new FileInputStream(file)) {
			return md5Hex(stream);
		}
	}

	private boolean compareFiles(final File a, final File b) throws IOException {
		if (a.length() != b.length()) {
			return false;
		}

		return fileMd5(a).equals(fileMd5(b));
	}

	protected boolean copyDatabaseIfDifferent(final File existingDB, final File newDB) throws IOException {
		if (filesEqual(existingDB, newDB)) {
			LOGGER.info("[" + getClass().getSimpleName() + "] database unchanged.");
			existingDB.setLastModified(newDB.lastModified());
			return false;
		}

		if (existingDB.isDirectory() && newDB.isDirectory()) {
			moveDirectory(existingDB, newDB);
			LOGGER.info("[" + getClass().getSimpleName() + "] Successfully updated database " + existingDB);
			return true;
		}

		if (existingDB != null && existingDB.exists()) {
			deleteDatabase(existingDB);
		}

		newDB.setReadable(true);
		newDB.setWritable(true);
		final boolean renamed = newDB.renameTo(existingDB);

		if (!renamed) {
			LOGGER.fatal("[" + getClass().getSimpleName() + "] Unable to rename " + newDB + " to " + existingDB.getAbsolutePath() + "; current working directory is " + System.getProperty("user.dir"));
			return false;
		}

		LOGGER.info("[" + getClass().getSimpleName() + "] Successfully updated database " + existingDB);
		return true;
	}

	private void moveDirectory(final File existingDB, final File newDB) throws IOException {
		LOGGER.info("[" + getClass().getSimpleName() + "] Moving Location database from: " + newDB + ", to: " + existingDB);

		for (final File file : existingDB.listFiles()) {
			file.setReadable(true);
			file.setWritable(true);
			file.delete();
		}

		existingDB.delete();
		Files.move(newDB.toPath(), existingDB.toPath(), StandardCopyOption.ATOMIC_MOVE);
	}

	private void deleteDatabase(final File db) {
		db.setReadable(true);
		db.setWritable(true);

		if (db.isDirectory()) {
			for (final File file : db.listFiles()) {
				file.delete();
			}
			LOGGER.debug("[" + getClass().getSimpleName() + "] Successfully deleted database under: " + db);
		} else {
			db.delete();
		}
	}

	protected boolean sourceCompressed = true;
	protected String tmpPrefix = "loc";
	protected String tmpSuffix = ".dat";

	@SuppressWarnings({"PMD.CyclomaticComplexity"})
	protected File downloadDatabase(final String url, final File existingDb) throws IOException {
		LOGGER.info("[" + getClass().getSimpleName() + "] Downloading database: " + url);
		final URL dbURL = new URL(url);
		final URLConnection conn = dbURL.openConnection();

		final long existingLastModified = existingDb.lastModified();
		if (conn instanceof HttpURLConnection && useModifiedTimestamp(existingDb)) {
			conn.setIfModifiedSince(existingLastModified);
			if (eTag != null) {
				conn.setRequestProperty("If-None-Match", eTag);
			}
		}

		final File outputFile = File.createTempFile(tmpPrefix, tmpSuffix);
		outputFile.setReadable(true);
		outputFile.setWritable(true);
		try (InputStream in = conn.getInputStream();
			 OutputStream out = new FileOutputStream(outputFile)
		) {
			if (conn instanceof HttpURLConnection) {
				eTag = conn.getHeaderField("ETag");
				if (((HttpURLConnection) conn).getResponseCode() == HttpURLConnection.HTTP_NOT_MODIFIED) {
					LOGGER.info("[" + getClass().getSimpleName() + "] " + url + " not modified since our existing database's last update time of " + new Date(existingLastModified));
					outputFile.delete();
					return existingDb;
				}
			} else if (dbURL.getProtocol().equals("file") && conn.getLastModified() > 0 && conn.getLastModified() <= existingLastModified) {
				LOGGER.info("[" + getClass().getSimpleName() + "] " + url + " not modified since our existing database's last update time of " + new Date(existingLastModified));
				outputFile.delete();
				return existingDb;
			}

			IOUtils.copy(sourceCompressed ? new GZIPInputStream(in) : in, out);
		}


		return outputFile;
	}

	private boolean useModifiedTimestamp(final File existingDb) {
		return existingDb != null && existingDb.exists() && existingDb.lastModified() > 0
				&& (!existingDb.isDirectory() || existingDb.listFiles().length > 0);
	}

	protected boolean needsUpdating(final File existingDB) {
		final long now = System.currentTimeMillis();
		final long fileTime = existingDB.lastModified();
		final long pollingIntervalInMS = getPollingInterval();
		return ((fileTime + pollingIntervalInMS) < now);
	}

	public void setLoaded(final boolean loaded) {
		this.loaded = loaded;
	}

	public boolean isLoaded() {
		return loaded;
	}

	public void setTrafficRouterManager(final TrafficRouterManager trafficRouterManager) {
		this.trafficRouterManager = trafficRouterManager;
	}

	public Path getDatabasesDirectory() {
		return databasesDirectory;
	}

	public void setDatabasesDirectory(final Path databasesDirectory) {
		this.databasesDirectory = databasesDirectory;
	}
}
