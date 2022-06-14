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

package org.apache.traffic_control.traffic_router.core.secure;

import org.apache.traffic_control.traffic_router.core.router.TrafficRouterManager;
import org.apache.traffic_control.traffic_router.core.util.ProtectedFetcher;
import org.apache.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import org.apache.traffic_control.traffic_router.shared.CertificateData;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.net.HttpURLConnection;
import java.util.ArrayList;
import java.util.Base64;
import java.util.List;

public class CertificatesClient {
	private static final Logger LOGGER = LogManager.getLogger(CertificatesClient.class);
	private TrafficOpsUtils trafficOpsUtils;
	private static final String PEM_FOOTER_PREFIX = "-----END";
	private long lastValidfetchTimestamp = 0L;
	private boolean shutdown = false;
	private TrafficRouterManager trafficRouterManager;

	public List<CertificateData> refreshData() {
		final StringBuilder stringBuilder = new StringBuilder();
		trafficRouterManager.trackEvent("lastHttpsCertificatesFetchAttempt");
		int status = fetchRawData(stringBuilder);

		while (status != HttpURLConnection.HTTP_NOT_MODIFIED && status != HttpURLConnection.HTTP_OK) {
			trafficRouterManager.trackEvent("lastHttpsCertificatesFetchFail");
			try {
				Thread.sleep(trafficOpsUtils.getConfigLongValue("certificates.retry.interval", 30 * 1000L));
			} catch (InterruptedException e) {
				if (!shutdown) {
					LOGGER.warn("Interrupted while pausing to fetch certificates from traffic ops", e);
				} else {
					return null;
				}
			}

			trafficRouterManager.trackEvent("lastHttpsCertificatesFetchAttempt");
			status = fetchRawData(stringBuilder);
		}

		if (status == HttpURLConnection.HTTP_NOT_MODIFIED) {
			return null;
		}

		lastValidfetchTimestamp = System.currentTimeMillis();
		trafficRouterManager.trackEvent("lastHttpsCertificatesFetchSuccess");
		return getCertificateData(stringBuilder.toString());
	}

	public int fetchRawData(final StringBuilder stringBuilder) {
		while (trafficOpsUtils == null || trafficOpsUtils.getHostname() == null || trafficOpsUtils.getHostname().isEmpty()) {
			LOGGER.error("No traffic ops hostname yet!");
			try {
				Thread.sleep(5000L);
			} catch (Exception e) {
				LOGGER.info("Interrupted while pausing for check of traffic ops config");
			}
		}

		final String certificatesUrl = trafficOpsUtils.getUrl("certificate.api.url", "https://${toHostname}/api/"+TrafficOpsUtils.TO_API_VERSION+"/cdns/name/${cdnName}/sslkeys");

		try {
			final ProtectedFetcher fetcher = new ProtectedFetcher(trafficOpsUtils.getAuthUrl(), trafficOpsUtils.getAuthJSON().toString(), 15000);
			return fetcher.getIfModifiedSince(certificatesUrl, 0L, stringBuilder);
		} catch (Exception e) {
			LOGGER.warn("Failed to fetch data for certificates from " + certificatesUrl + "(" + e.getClass().getSimpleName() + ") : " + e.getMessage(), e);
		}

		return -1;
	}

	public List<CertificateData> getCertificateData(final String jsonData) {
		try {
			LOGGER.debug("Certificates successfully updated @ "+lastValidfetchTimestamp);
			return ((CertificatesResponse) new ObjectMapper().readValue(jsonData, new TypeReference<CertificatesResponse>() { })).getResponse();
		} catch (Exception e) {
			LOGGER.warn("Failed parsing json data: " + e.getMessage());
		}

		return new ArrayList<>();
	}

	public String[] doubleDecode(final String encoded) {
		final byte[] decodedBytes = Base64.getMimeDecoder().decode(encoded.getBytes());

		final List<String> encodedPemItems = new ArrayList<>();

		final String[] lines = new String(decodedBytes).split("\\r?\\n");
		final StringBuilder builder = new StringBuilder();

		for (final String line : lines) {
			if (line.startsWith(PEM_FOOTER_PREFIX)) {
				encodedPemItems.add(builder.toString());
				builder.setLength(0);
			}

			builder.append(line);
		}

		if (encodedPemItems.isEmpty()) {
			if (builder.length() == 0) {
				LOGGER.warn("Failed base64 decoding");
			 } else {
				encodedPemItems.add(builder.toString());
			}
		}

		return encodedPemItems.toArray(new String[0]);
	}

	public void setTrafficOpsUtils(final TrafficOpsUtils trafficOpsUtils) {
		this.trafficOpsUtils = trafficOpsUtils;
	}

	public void setShutdown(final boolean shutdown) {
		this.shutdown = true;
	}

	public TrafficRouterManager getTrafficRouterManager() {
		return trafficRouterManager;
	}

	public void setTrafficRouterManager(final TrafficRouterManager trafficRouterManager) {
		this.trafficRouterManager = trafficRouterManager;
	}
}
