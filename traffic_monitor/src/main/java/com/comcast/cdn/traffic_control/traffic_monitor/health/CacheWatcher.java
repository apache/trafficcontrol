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

package com.comcast.cdn.traffic_control.traffic_monitor.health;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;

import com.comcast.cdn.traffic_control.traffic_monitor.wicket.models.CacheDataModel;
import org.apache.log4j.Logger;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.RouterConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;

import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus.OFFLINE;
import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus.ONLINE;

public class CacheWatcher {
	private static final Logger LOGGER = Logger.getLogger(CacheWatcher.class);

	private HealthDeterminer myHealthDeterminer;

	private static final List<CacheDataModel> list = new ArrayList<CacheDataModel>();
	private static final CacheDataModel itercount = new CacheDataModel("Iteration Count");
	private static final CacheDataModel fetchCount = new CacheDataModel("Fetch Count");
	private static final CacheDataModel errorCount = new CacheDataModel("Error Count");
	private static final CacheDataModel queryInterval = new CacheDataModel("Last Query Interval");
	private static final CacheDataModel queryIntervalActual = new CacheDataModel("Query Interval Actual");
	private static final CacheDataModel queryIntervalTarget = new CacheDataModel("Query Interval Target");
	private static final CacheDataModel queryIntervalDelta = new CacheDataModel("Query Interval Delta");
	private static final CacheDataModel freeMem = new CacheDataModel("Free Memory (MB)");
	private static final CacheDataModel totalMem = new CacheDataModel("Total Memory (MB)");
	private static final CacheDataModel maxMemory = new CacheDataModel("Max Memory (MB)");
	final MonitorConfig config = ConfigHandler.getInstance().getConfig();
	private final List<CacheStateUpdater> cacheStateUpdaters = new ArrayList<CacheStateUpdater>();
	boolean isActive = true;

	private FetchService mainThread;

	private final CacheStateRegistry cacheStateRegistry = CacheStateRegistry.getInstance();
	private final CacheStatisticsClient cacheStatisticsClient = new CacheStatisticsClient();

	public static List<CacheDataModel> getProps() {
		return list;
	}

	public CacheWatcher init(final HealthDeterminer hd) {
		myHealthDeterminer = hd;
		list.add(itercount);
		list.add(fetchCount);
		list.add(errorCount);
		list.add(queryInterval);
		list.add(queryIntervalActual);
		list.add(queryIntervalTarget);
		list.add(queryIntervalDelta);
		list.add(totalMem);
		list.add(freeMem);
		list.add(maxMemory);

		mainThread = new FetchService();
		mainThread.start();
		return this;
	}

	class FetchService extends Thread {
		public FetchService() {
		}

		final Runtime runtime = Runtime.getRuntime();

		private List<CacheState> checkCaches(final RouterConfig crConfig, final AtomicInteger failCount) {
			maxMemory.set(runtime.maxMemory() / (1024 * 1024));
			totalMem.set(runtime.totalMemory() / (1024 * 1024));
			freeMem.set(runtime.freeMemory() / (1024 * 1024));

			final List<CacheState> cacheStates = new ArrayList<CacheState>();

			for (Cache cache : crConfig.getCacheList()) {

				if (!isActive) {
					// destroy was called, do stop fetching
					return cacheStates;
				}

				if (!myHealthDeterminer.shouldMonitor(cache)) {
					continue;
				}

				final CacheState state = cacheStateRegistry.update(cache);

				cacheStates.add(state);

				if (!shouldFetchStats(cache)) {
					cache.setState(state, myHealthDeterminer);
					continue;
				}

				state.prepareStatisticsForUpdate();

				fetchCount.inc();
				state.putDataPoint("_queryUrl_", cache.getStatisticsUrl());
				state.setHistoryTime(cache.getHistoryTime());

				final long requestTimeout = System.currentTimeMillis() + myHealthDeterminer.getConnectionTimeout(cache, 2000);

				final CacheStateUpdater cacheStateUpdater = new CacheStateUpdater(state, errorCount).update(myHealthDeterminer, failCount, requestTimeout);
				cacheStateUpdaters.add(cacheStateUpdater);
				cacheStatisticsClient.fetchCacheStatistics(cache, cacheStateUpdater);

				cacheTimePad();
			}

			return cacheStates;
		}

		private void cacheTimePad() {
			if (config == null) {
				return;
			}

			final int t = config.getCacheTimePad();

			if (t == 0) {
				return;
			}

			try {
				Thread.sleep(t);
			} catch (InterruptedException e) {
				// Ignore
			}
		}

		public void run() {
			while (true) {
				if (!isActive) {
					LOGGER.warn("Not active");
					return;
				}

				try {
					final long time = System.currentTimeMillis();
					final RouterConfig crConfig = RouterConfig.getCrConfig();

					if (crConfig == null && config != null) {
						try {
							Thread.sleep(config.getHealthPollingInterval());
						} catch (InterruptedException e) {
							// Ignore
						}

						LOGGER.warn("No router config available, skipping health check");
						continue;
					}

					final AtomicInteger failCount = new AtomicInteger(0);
					final List<CacheState> states = checkCaches(crConfig, failCount);

					boolean waitForFinish = true;
					final AtomicInteger cancelCount = new AtomicInteger(0);

					while (waitForFinish) {
						waitForFinish = false;

						for (CacheStateUpdater updater : cacheStateUpdaters) {
							waitForFinish |= !updater.completeFetchStatistics(cancelCount);
						}
					}

					cacheStateUpdaters.clear();
					cacheStateRegistry.removeAllBut(states);
					final long completedTime = System.currentTimeMillis();

					try {
						Thread.sleep(Math.max(config.getHealthPollingInterval() - (completedTime - time), 0));
					} catch (InterruptedException e) {
						// Ignore
					}

					itercount.inc();

					final long mytime = System.currentTimeMillis() - time;

					queryInterval.set(mytime);
					queryIntervalTarget.set(config.getHealthPollingInterval());
					queryIntervalActual.set(completedTime - time);
					queryIntervalDelta.set((completedTime - time) - config.getHealthPollingInterval());

					LOGGER.debug("Check time of " + states.size() + " caches elapsed: " + mytime + " msec, (Active time was " + (completedTime - time) + ") msec, " + cancelCount.get() + " checks were cancelled, " + failCount.get() + " failed");
				} catch (Exception e) {
					LOGGER.warn(e, e);

					try {
						Thread.sleep(100);
					} catch (InterruptedException ex) {
						// Ignore
					}
				}
			}
		}
	}


	public void destroy() {
		LOGGER.warn("CacheWatcher: shutting down ");

		isActive = false;
		final long time = System.currentTimeMillis();

		mainThread.interrupt();
		cacheStatisticsClient.shutdown();

		while (mainThread.isAlive()) {
			try {
				Thread.sleep(10);
			} catch (InterruptedException e) {
				// Ignore
			}
		}
		LOGGER.warn("Stopped: Termination time: " + (System.currentTimeMillis() - time));
	}

	public long getCycleCount() {
		return itercount.getRawValue();
	}

	public boolean shouldFetchStats(final Cache cache) {
		HealthDeterminer.AdminStatus adminStatus = HealthDeterminer.AdminStatus.valueOf(cache.getStatus());
		return (adminStatus != OFFLINE && adminStatus != ONLINE);
	}
}
