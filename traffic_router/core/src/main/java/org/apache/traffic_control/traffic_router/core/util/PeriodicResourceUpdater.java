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

import java.io.BufferedReader;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.FileReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.StringReader;
import java.net.URI;
import java.net.URISyntaxException;
import java.nio.channels.FileLock;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.zip.GZIPInputStream;

import org.apache.commons.io.IOUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import org.asynchttpclient.AsyncCompletionHandler;
import org.asynchttpclient.AsyncHttpClient;
import org.asynchttpclient.DefaultAsyncHttpClient;
import org.asynchttpclient.DefaultAsyncHttpClientConfig;
import org.asynchttpclient.Request;
import org.asynchttpclient.Response;

import static org.apache.commons.codec.digest.DigestUtils.md5Hex;

/**
 * 
 * @author jlaue
 *
 */
public class PeriodicResourceUpdater {
	private static final Logger LOGGER = LogManager.getLogger(PeriodicResourceUpdater.class);

	private AsyncHttpClient asyncHttpClient;
	protected String databaseLocation;
	protected final ResourceUrl urls;
	protected ScheduledExecutorService executorService = Executors.newSingleThreadScheduledExecutor();
	protected long pollingInterval;

	private static final String GZIP_ENCODING_STRING = "gzip";

	protected ScheduledFuture<?> scheduledService;

	public PeriodicResourceUpdater(final AbstractUpdatable listener, final ResourceUrl urls, final String location, final int interval, final boolean pauseTilLoaded) {
		this.listener = listener;
		this.urls = urls;
		databaseLocation = location;
		pollingInterval = interval;
		this.pauseTilLoaded = pauseTilLoaded;
	}

	public void destroy() {
		executorService.shutdownNow();

		while (!asyncHttpClient.isClosed()) {
			try {
				asyncHttpClient.close();
			} catch (IOException e) {
				LOGGER.error(e.getMessage());
			}
		}
	}

	/**
	 * Gets pollingInterval.
	 * 
	 * @return the pollingInterval
	 */
	public long getPollingInterval() {
		if(pollingInterval == 0) { return 66000; }
		return pollingInterval;
	}

	final private Runnable updater = new Runnable() {
		@Override
		public void run() {
			updateDatabase();
		}
	};

	private boolean hasBeenLoaded = false;

	final private AbstractUpdatable listener;
	final private boolean pauseTilLoaded;

	public void init() {
		asyncHttpClient = newAsyncClient();
		putCurrent();
		LOGGER.info("Starting schedule with interval: "+getPollingInterval() + " : "+TimeUnit.MILLISECONDS);
		scheduledService = executorService.scheduleWithFixedDelay(updater, 0, getPollingInterval(), TimeUnit.MILLISECONDS);
		// wait here until something is loaded
		final File existingDB = new File(databaseLocation);
		if(pauseTilLoaded ) {
			while(!existingDB.exists()) {
				LOGGER.info("Waiting for valid: " + databaseLocation);
				try {
					Thread.sleep(getPollingInterval());
				} catch (InterruptedException e) {
				}
			}
		}
	}

	private AsyncHttpClient newAsyncClient() {
		return new DefaultAsyncHttpClient(
				new DefaultAsyncHttpClientConfig.Builder()
						.setFollowRedirect(true)
							.setConnectTimeout(10000)
								.build());

	}

	private synchronized void putCurrent() {
		final File existingDB = new File(databaseLocation);
		if(existingDB.exists()) {
			try {
				listener.update(IOUtils.toString(new FileReader(existingDB)));
			} catch (Exception e) {
				LOGGER.warn(e,e);
			}
		}
	}

	public synchronized boolean updateDatabase() {
		final File existingDB = new File(databaseLocation);
		try {
			if (!hasBeenLoaded || needsUpdating(existingDB)) {
				final Request request = getRequest(urls.nextUrl());
				if (request != null) {
					request.getHeaders().add("Accept-Encoding", GZIP_ENCODING_STRING);
					if ((asyncHttpClient!=null) && (!asyncHttpClient.isClosed())) {
						asyncHttpClient.executeRequest(request, new UpdateHandler(request)); // AsyncHandlers are NOT thread safe; one instance per request
					}
					return true;
				}
			} else {
				LOGGER.info("Database " + existingDB.getAbsolutePath() + " does not require updating.");
			}
		} catch (final Exception e) {
			LOGGER.warn(e.getMessage(), e);
		}
		return false;
	}

	public boolean updateDatabase(final String newDB) {
		final File existingDB = new File(databaseLocation);
		try {
			if (newDB != null && !filesEqual(existingDB, newDB)) {
				listener.cancelUpdate();
				if (listener.update(newDB)) {
					copyDatabase(existingDB, newDB);
					LOGGER.info("updated " + existingDB.getAbsolutePath());
					listener.setLastUpdated(System.currentTimeMillis());
					listener.complete();
				} else {
					LOGGER.warn("File rejected: " + existingDB.getAbsolutePath());
				}
			} else {
				listener.noChange();
			}
			hasBeenLoaded = true;
			return true;
		} catch (final Exception e) {
			LOGGER.warn(e.getMessage(), e);
		}
		return false;
	}

	public void setDatabaseLocation(final String databaseLocation) {
		this.databaseLocation = databaseLocation;
	}

	/**
	 * Sets executorService.
	 * 
	 * @param es
	 *            the executorService to set
	 */
	public void setExecutorService(final ScheduledExecutorService es) {
		executorService = es;
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

	private String fileMd5(final File file) throws IOException {
		try (FileInputStream stream = new FileInputStream(file)) {
			return md5Hex(stream);
		}
	}

	boolean filesEqual(final File a, final String newDB) throws IOException {
		if (!a.exists()) {
			return newDB == null;
		}

		if (newDB == null) {
			return false;
		}

		if (a.length() != newDB.length()) {
			return false;
		}

		try (InputStream newDBStream = IOUtils.toInputStream(newDB)) {
			return fileMd5(a).equals(md5Hex(newDBStream));
		}
	}

	protected synchronized void copyDatabase(final File existingDB, final String newDB) throws IOException {
		try (StringReader in = new StringReader(newDB);
			FileOutputStream out = new FileOutputStream(existingDB);
			FileLock lock = out.getChannel().tryLock()) {

			if (lock == null) {
				LOGGER.error("Database " + existingDB.getAbsolutePath() + " locked by another process.");
				return;
			}

			IOUtils.copy(in, out);
			existingDB.setReadable(true, false);
			existingDB.setWritable(true, true);
			lock.release();
		}
	}

	protected boolean needsUpdating(final File existingDB) {
		final long now = System.currentTimeMillis();
		final long fileTime = existingDB.lastModified();
		final long pollingIntervalInMS = getPollingInterval();
		return ((fileTime + pollingIntervalInMS) < now);
	}

	private class UpdateHandler extends AsyncCompletionHandler<Object> {
		final Request request;
		public UpdateHandler(final Request request) {
			this.request = request;
		}

		@Override
		public Integer onCompleted(final Response response) throws IOException {
			// Do something with the Response
			final int code = response.getStatusCode();

			if (code != 200) {
				if (code >= 400) {
					LOGGER.warn("failed to GET " + response.getUri() + " - returned status code " + code);
				}
				return code;
			}

			final String responseBody;
			if (GZIP_ENCODING_STRING.equals(response.getHeader("Content-Encoding"))) {
				final StringBuilder stringBuilder = new StringBuilder();
				try (GZIPInputStream zippedInputStream =  new GZIPInputStream(response.getResponseBodyAsStream());
					 BufferedReader r = new BufferedReader(new InputStreamReader(zippedInputStream))) {
					String line;
					while ((line = r.readLine()) != null) {
						stringBuilder.append(line);
					}
				}
				responseBody = stringBuilder.toString();
			} else {
				responseBody = response.getResponseBody();
			}

			updateDatabase(responseBody);

			return code;
		}

		@Override
		public void onThrowable(final Throwable t){
			LOGGER.warn("Failed request " + request.getUrl() + ": " + t, t);
		}
	};

	private Request getRequest(final String url) {
		try {
			new URI(url);
			return asyncHttpClient.prepareGet(url).setFollowRedirect(true).build();
		} catch (URISyntaxException e) {
			LOGGER.fatal("Cannot update database from Bad URI - " + url);
			return null;
		}
	}
}
