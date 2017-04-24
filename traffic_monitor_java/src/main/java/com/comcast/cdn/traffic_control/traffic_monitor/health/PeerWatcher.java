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

import java.util.Map;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.RouterConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.config.Peer;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.util.Fetcher;

public class PeerWatcher {
	private static final Logger LOGGER = Logger.getLogger(PeerWatcher.class);
	private FetchService mainThread;
	private boolean isActive = true;

	final MonitorConfig config = ConfigHandler.getInstance().getConfig();

	public PeerWatcher init() {
		mainThread = new FetchService();
		mainThread.start();
		return this;
	}

	class FetchService extends Thread {
		public FetchService() {
		}

		final Runtime runtime = Runtime.getRuntime();

		public void run() { // run the service
			ExecutorService pool = null;

			while (true) {
				if(!isActive ) {
					return;
				}

				final long time = System.currentTimeMillis();
				final RouterConfig crConfig = RouterConfig.getCrConfig();

				if (crConfig == null) {
					try {
						Thread.sleep(config.getPeerPollingInterval());
					} catch (Exception e) { }
					continue;
				}

				final Map<String, Peer> peerMap = crConfig.getPeerMap();

				try {
					final int poolSize = config.getPeerThreadPool();
					pool = Executors.newFixedThreadPool(poolSize);
					checkPeers(pool, peerMap);
				} catch(Exception e) {
					LOGGER.warn(e,e);
					if(!isActive) { return; }
				}

				try {
					pool.shutdown();
					Thread.sleep(config.getPeerPollingInterval());
				} catch (Exception e) { }

				try {
					while(!pool.awaitTermination(1, TimeUnit.SECONDS)) {
						LOGGER.warn("Pool did not terminate");
					}
				} catch (Exception e) { }

				final long mytime = System.currentTimeMillis()-time;
			}
		}

	}

	private void checkPeers(final ExecutorService pool, final Map<String, Peer> peerMap) {
		Map<String, Peer> myMap = null;
		synchronized(PeerWatcher.this) {
			myMap = peerMap;
		}

		final String urlPattern = config.getPeerUrl(); // http://${hostname}/publish/CrStates?raw

		if (myMap == null || myMap.isEmpty()) {
			return;
		}

		for (String key : myMap.keySet()) {
			if (!isActive) {
				return;
			}

			final Peer peer = myMap.get(key);
			final PeerState peerState = PeerState.getOrCreate(peer);

			pool.execute(getHandler(peerState, urlPattern));
		}
	}

	private Runnable getHandler(final PeerState peerState, final String urlPattern) {
		return new Runnable(){
			@Override
			public void run() {
				if (!isActive) {
					return;
				}

				final Peer peer = peerState.getPeer();
				final String url = urlPattern.replace("${hostname}", peer.getIpAddress()).
							replace("${port}", peer.getPortString());
				final String prettyUrl = urlPattern.replace("${hostname}", peer.getFqdn()).
							replace("${port}", peer.getPortString());

				try {
					final String result = Fetcher.fetchContent(url, peer.getHeaderMap(), config.getConnectionTimeout());
					final JSONObject jr = new JSONObject(result);
					final JSONObject cacheStates = jr.getJSONObject("caches");

					peerState.setReachable(true);
					peerState.prepareStatisticsForUpdate();

					for (String id : JSONObject.getNames(cacheStates)) {
						final JSONObject cache = cacheStates.getJSONObject(id);
						peerState.putDataPoint(id, cache.optString(AbstractState.IS_AVAILABLE_STR));
					}
				} catch (Exception e) {
					peerState.setReachable(false, e.getMessage());
					LOGGER.warn(e + " to " + prettyUrl);
				}

				final MonitorConfig config = ConfigHandler.getInstance().getConfig();
				final RouterConfig crConfig = RouterConfig.getCrConfig();

				if (crConfig != null && config.getPeerOptimistic()) {
					for (Cache cache : crConfig.getCacheList()) {
						if (!cache.isAvailable()) {
							PeerState.logOverride(cache);
						} else {
							PeerState.clearOverride(cache);
						}
					}
				}
			}
		};
	}

	public void destroy() {
		LOGGER.warn("PeerWatcher: shutting down ");
		isActive  = false;
		final long time = System.currentTimeMillis();

		mainThread.interrupt();

		while (mainThread.isAlive()) {
			try {
				Thread.sleep(10);
			} catch (InterruptedException e) {
			}
		}

		LOGGER.warn("Stopped: Termination time: "+(System.currentTimeMillis() - time));
	}
}
