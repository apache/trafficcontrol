package com.comcast.cdn.traffic_control.traffic_router.neustar.configuration;

import com.comcast.cdn.traffic_control.traffic_router.neustar.NeustarGeolocationService;
import com.comcast.cdn.traffic_control.traffic_router.neustar.data.NeustarDatabaseUpdater;
import org.apache.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;

public class ServiceRefresher implements Runnable {
	private final Logger logger = Logger.getLogger(ServiceRefresher.class);

	@Autowired
	NeustarDatabaseUpdater neustarDatabaseUpdater;

	@Autowired
	NeustarGeolocationService neustarGeolocationService;

	@Override
	public void run() {
		try {
			if (neustarDatabaseUpdater.update() || !neustarGeolocationService.isInitialized()) {
				neustarGeolocationService.reloadDatabase();
			}
		} catch (Exception e) {
			logger.error("Failed to refresh Neustar Geolocation Service:" + e.getMessage());
		}
	}
}
