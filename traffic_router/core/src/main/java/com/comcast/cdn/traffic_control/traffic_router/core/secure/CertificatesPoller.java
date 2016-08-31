package com.comcast.cdn.traffic_control.traffic_router.core.secure;

import com.comcast.cdn.traffic_control.traffic_router.configuration.ConfigurationListener;
import org.apache.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;

public class CertificatesPoller implements ConfigurationListener {
	private static final Logger LOGGER = Logger.getLogger(CertificatesPoller.class);

	private final ScheduledExecutorService executor;
	private ScheduledFuture future;
	private CertificatesClient certificatesClient;
	private static final long defaultFixedRate = 3600 * 1000L;
	private static final String intervalProperty = "certificates.polling.interval";
	private long pollingInterval = defaultFixedRate;

	@Autowired
	private Environment environment;

	public CertificatesPoller() {
		executor = Executors.newSingleThreadScheduledExecutor();
	}

	public Long getEnvironmentPollingInterval() {
		if (environment == null) {
			LOGGER.warn("Could not find Environment object!");
		}

		try {
			final Long value = environment.getProperty(intervalProperty, Long.class);
			if (value == null) {
				LOGGER.info("No custom value for " + intervalProperty);
			}

			return value;
		} catch (Exception e) {
			LOGGER.warn("Failed to get value of " + intervalProperty + ": " + e.getMessage());
			return null;
		}
	}

	@SuppressWarnings("PMD.AvoidCatchingThrowable")
	public void start() {
		final Runnable runnable = () -> {
			try {
				certificatesClient.refreshData();
			} catch (Throwable t) {
				LOGGER.warn("Failed to refresh certificate data: " + t.getClass().getCanonicalName() + " " + t.getMessage(), t);
			}
		};

		final Long customFixedRate = getEnvironmentPollingInterval();

		if (customFixedRate == null) {
			LOGGER.info("Using default fixed rate polling interval " + pollingInterval + " msec");
		} else {
			LOGGER.info("Using custom fixed rate polling interval " + customFixedRate + " msec");
			pollingInterval = customFixedRate;
		}

		future = executor.scheduleWithFixedDelay(runnable, 0, pollingInterval, TimeUnit.MILLISECONDS);
		LOGGER.info("Polling for certificates every " + pollingInterval + " msec");
	}

	public void stop() {
		if (future != null) {
			future.cancel(false);
		}
	}

	public void destroy() {
		executor.shutdownNow();
	}

	public void setCertificatesClient (final CertificatesClient certificatesClient) {
		this.certificatesClient = certificatesClient;
	}

	private boolean futureIsDone() {
		return future == null || future.isDone() || future.isCancelled();
	}

	public void restart() {
		stop();
		while (!futureIsDone()) {
			try {
				Thread.sleep(250L);
			} catch (InterruptedException e) {
				LOGGER.info("Interrupted sleep while waiting for certificate poller future to finish");
			}
		}

		start();
	}

	public long getPollingInterval() {
		return pollingInterval;
	}

	@Override
	public void configurationChanged() {
		restart();
	}
}
