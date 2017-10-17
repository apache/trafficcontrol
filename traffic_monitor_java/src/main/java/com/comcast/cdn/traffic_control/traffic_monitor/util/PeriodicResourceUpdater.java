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

package com.comcast.cdn.traffic_control.traffic_monitor.util;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.nio.channels.FileLock;
import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;

import org.apache.commons.io.IOUtils;
import org.apache.log4j.Logger;
import org.apache.wicket.model.IModel;
import org.apache.wicket.model.Model;

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;

/** 
 * 
 * @author jlaue
 *
 */
public class PeriodicResourceUpdater {
	private static final Logger LOGGER = Logger.getLogger(PeriodicResourceUpdater.class);
	private boolean isActive = true;
	private boolean running = false;

	protected IModel<Long> pollingInterval;
	protected IModel<String> host;

//	static protected ScheduledExecutorService executorService;
//	protected ScheduledFuture<?> scheduledService;

	public PeriodicResourceUpdater(final IModel<Long> interval) {
		pollingInterval = interval;
	}
	
	private final List<UpdateModel> umList = new ArrayList<UpdateModel>();
	private static class UpdateModel {
		protected List<Model<String>> urlList = new LinkedList<Model<String>>();
		protected String databaseLocation;
		private Updatable listener;
		private boolean hasBeenLoaded = false;

		private void putCurrent() {
			final File existingDB = ConfigHandler.getInstance().getDbFile(databaseLocation);

			if (existingDB.exists()) {
				LOGGER.warn("loading: " + existingDB.getAbsolutePath());
				listener.update(existingDB);
			}

			hasBeenLoaded = true;
		}
	}

	public void add(final Updatable listener, final Model<String> url, final String location) {
		final UpdateModel um = new UpdateModel();
		um.listener = listener;
		um.urlList.add(url);
		um.databaseLocation = location;
		add(um);
	}

	public void add(final Updatable listener, final String[] urla, final String location) {
		final UpdateModel um = new UpdateModel();
		um.listener = listener;
		for(String url : urla) {
			um.urlList.add(new Model<String>(url));
		}
		um.databaseLocation = location;
		add(um);
	}
	
	private void add(final UpdateModel um) {
		um.putCurrent();
		synchronized(umList) {
			umList.add(um);
		}
	}

//	static {
//		executorService = java.util.concurrent.Executors.newSingleThreadScheduledExecutor();
//	}

	public void destroy() {
//		executorService.shutdownNow();
		isActive = false;
		mainThread.interrupt();
		while (running) {
			try {
				Thread.sleep(10);
			} catch (InterruptedException e) { } 
		}
	}

	final private Runnable updater = new Runnable() {
		@Override
		public void run() {
			running = true;
			while(isActive) {
				try {
					synchronized(umList) {
						for(UpdateModel um : umList) {
							if(!isActive) {
								running = false;
								return;
							}
							updateDatabase(um);
						}
					}
				} catch(Exception e) {
					LOGGER.warn("error", e);
				}
				try {
					Thread.sleep(getPollingInterval());
				} catch (InterruptedException e) { } 
			}
			running = false;
		}
	};
	private Thread mainThread;
	public long getPollingInterval() {
		return pollingInterval.getObject().longValue();
	}

	public void init() {
		mainThread = new Thread(updater);
		mainThread.start();
	}

	public void forceUpdate() {
		mainThread.interrupt();
	}
	protected File fetchFile(final String url) throws IOException {
		return Fetcher.downloadFile(url);
	}

	public boolean updateDatabase(final UpdateModel um) {
		final File existingDB = ConfigHandler.getInstance().getDbFile(um.databaseLocation);
		File newDB = null;
		try {
			if (um.hasBeenLoaded) {
				int urlIndex = (int)(Math.random()*um.urlList.size());
				for(int i = 0; i < um.urlList.size(); i++) {
					final String url = um.urlList.get(urlIndex).getObject();
					try {
						newDB = fetchFile(url);
					} catch(Exception e) {
						LOGGER.error("Error with '" + url + "' : " + e);
						urlIndex = (urlIndex+1)%um.urlList.size();
						continue;
					}
					break;
				}
				if(newDB != null && !filesEqual(existingDB, newDB)) {
					if(um.listener.update(newDB)) {
						copyDatabase(existingDB, newDB);
						LOGGER.debug("File saved: "+existingDB.getAbsolutePath());
					} else {
						LOGGER.debug("File rejected: "+existingDB.getAbsolutePath());
						Fetcher.clearTmCookie();
					}
				} else {
					LOGGER.debug("File unchanged: "+existingDB.getAbsolutePath());
				}
				return true;
			} 
		} catch (final Exception e) {
			LOGGER.warn(e.getMessage(), e);
		} finally {
			if (newDB != null) {
				newDB.delete();
			}
		}
		return false;
	}
	static boolean filesEqual(final File a, final File b) throws IOException {
		if(!a.exists() && !b.exists()) { return true; }
		if(!a.exists() || !b.exists()) { return false; }
		if(a.length() != b.length()) { return false; }
		FileInputStream fis = new FileInputStream(a);
		final String md5a = org.apache.commons.codec.digest.DigestUtils.md5Hex(fis);
		fis.close();
		fis = new FileInputStream(b);
		final String md5b = org.apache.commons.codec.digest.DigestUtils.md5Hex(fis);
		fis.close();
		if(md5a.equals(md5b)) { return true; }
		return false;
	}
	static void copyDatabase(final File existingDB, final File newDB) throws IOException {
		final FileInputStream in = new FileInputStream(newDB);
		final FileOutputStream out = new FileOutputStream(existingDB);
		final FileLock lock = out.getChannel().tryLock();
		if (lock != null) {
			LOGGER.info("Updating location database.");
			IOUtils.copy(in, out);
			existingDB.setReadable(true, false);
			existingDB.setWritable(true, false);
			lock.release();
			LOGGER.info("Successfully updated location database.");
		} else {
			LOGGER.info("Location database locked by another process.");
		}
		IOUtils.closeQuietly(in);
		IOUtils.closeQuietly(out);
	}
}
