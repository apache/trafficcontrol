/*
 * Copyright 2016 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.secure;

import com.comcast.cdn.traffic_control.traffic_router.configuration.ConfigurationListener;
import com.comcast.cdn.traffic_control.traffic_router.core.config.CertificateChecker;
import com.comcast.cdn.traffic_control.traffic_router.shared.CertificateData;
import org.apache.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;

import java.util.List;
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
	private CertificateChecker certificateChecker;

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
				List<CertificateData> certificateDataList = certificatesClient.refreshData();
				if (certificateDataList != null) {
					certificateChecker.setCertificateDataList(certificateDataList);
				}
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

	public CertificateChecker getCertificateChecker() {
		return certificateChecker;
	}

	public void setCertificateChecker(final CertificateChecker certificateChecker) {
		this.certificateChecker = certificateChecker;
	}

	@Override
	public void configurationChanged() {
		restart();
	}
}
