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

package com.comcast.cdn.traffic_control.traffic_monitor.health;

import java.util.List;

import org.apache.log4j.Logger;

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.RouterConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;

public class DsWatcher {
	private static final Logger LOGGER = Logger.getLogger(DsWatcher.class);

	private HealthDeterminer myHealthDeterminer;

	final MonitorConfig config = ConfigHandler.getConfig();
	
	boolean isActive = true;

	private FetchService mainThread;

	public DsWatcher init(final HealthDeterminer hd) {
		myHealthDeterminer = hd;
		mainThread = new FetchService();
		mainThread.start();
		return this;
	}

	class FetchService extends Thread {
		public FetchService() {
		}

		public void run() { // run the service
			while(true) {
				try {
					final long time = System.currentTimeMillis();
					final RouterConfig crConfig = RouterConfig.getCrConfig();
					if(crConfig == null) {
						try {
							Thread.sleep(config.getHealthPollingInterval());
						} catch (InterruptedException e) { }
						continue;
					}
	
					final List<CacheState> states = CacheState.getCacheStates();
					DsState.startUpdateAll();
					DsState.completeAll(states, myHealthDeterminer, crConfig.getDsList(), time-config.getDsCacheLeniency());
					try {
						Thread.sleep(Math.max(config.getHealthDsInterval()-(System.currentTimeMillis()-time),0));
					} catch (InterruptedException e) { }
					final long mytime = System.currentTimeMillis()-time;
					LOGGER.debug("Pool time elapsed: "+mytime);
				} catch (Exception e) {
					LOGGER.warn(e,e);
					try {
						Thread.sleep(100);
					} catch (InterruptedException ex) { }
				}
				if(!isActive) { return; }
			}
		}

	}

	public void destroy() {
		LOGGER.warn("CacheWatcher: shutting down ");
		isActive  = false;
		final long time = System.currentTimeMillis();
		mainThread.interrupt();
		CacheState.shutdown();
		while(mainThread.isAlive()) {
			try {
				Thread.sleep(10);
			} catch (InterruptedException e) {
			}
		}
		LOGGER.warn("Stopped: Termination time: "+(System.currentTimeMillis() - time));	
	}

}



