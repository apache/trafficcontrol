package com.comcast.cdn.traffic_control.traffic_router.neustar.configuration;

import com.comcast.cdn.traffic_control.traffic_router.configuration.ConfigurationListener;
import org.apache.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;

import javax.annotation.PostConstruct;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;

public class TrafficRouterConfigurationListener implements ConfigurationListener {
	private final Logger logger = Logger.getLogger(TrafficRouterConfigurationListener.class);

	@Autowired
	private Environment environment;

	@Autowired
	ScheduledExecutorService scheduledExecutorService;

	@Autowired
	ServiceRefresher serviceRefresher;

	private ScheduledFuture<?> scheduledFuture;

	@PostConstruct
	@Override
	public void configurationChanged() {
		if (scheduledFuture != null) {
			scheduledFuture.cancel(true);

			while (!scheduledFuture.isDone()) {
				try {
					Thread.sleep(100L);
				} catch (InterruptedException e) {
					// ignore
				}
			}
		}

		Long fixedRate = environment.getProperty("neustar.polling.interval", Long.class, 86400000L);
		scheduledFuture = scheduledExecutorService.scheduleAtFixedRate(serviceRefresher, 0L, fixedRate, TimeUnit.MILLISECONDS);
		logger.warn("Restarting Neustar remote database refresher at rate " + fixedRate + " msec");
	}
}
