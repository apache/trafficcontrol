package com.comcast.cdn.traffic_control.traffic_router.core.secure;

import com.comcast.cdn.traffic_control.traffic_router.core.util.ProtectedFetcher;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.comcast.cdn.traffic_control.traffic_router.shared.CertificateData;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.log4j.Logger;

import java.net.HttpURLConnection;
import java.util.ArrayList;
import java.util.Base64;
import java.util.List;

public class CertificatesClient {
	private static final Logger LOGGER = Logger.getLogger(CertificatesClient.class);
	private TrafficOpsUtils trafficOpsUtils;
	private static final String PEM_FOOTER_PREFIX = "-----END";
	private long lastValidfetchTimestamp = 0L;

	public List<CertificateData> refreshData() {
		final StringBuilder stringBuilder = new StringBuilder();
		int status = fetchRawData(stringBuilder);

		while (status != HttpURLConnection.HTTP_NOT_MODIFIED && status != HttpURLConnection.HTTP_OK) {
			try {
				Thread.sleep(trafficOpsUtils.getConfigLongValue("certificates.retry.interval", 30 * 1000L));
			} catch (InterruptedException e) {
				LOGGER.warn("Interrupted while pausing to fetch certificates from traffic ops", e);
			}
			status = fetchRawData(stringBuilder);
		}

		if (status == HttpURLConnection.HTTP_NOT_MODIFIED) {
			return null;
		}

		lastValidfetchTimestamp = System.currentTimeMillis();
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

		final String certificatesUrl = trafficOpsUtils.getUrl("certificate.api.url", "https://${toHostname}/api/1.2/cdns/name/${cdnName}/sslkeys.json");

		try {
			final ProtectedFetcher fetcher = new ProtectedFetcher(trafficOpsUtils.getAuthUrl(), trafficOpsUtils.getAuthJSON().toString(), 15000);
			return fetcher.getIfModifiedSince(certificatesUrl, lastValidfetchTimestamp, stringBuilder);
		} catch (Exception e) {
			LOGGER.warn("Failed to fetch data for certificates from " + certificatesUrl + "(" + e.getClass().getSimpleName() + ") : " + e.getMessage(), e);
		}

		return -1;
	}

	public List<CertificateData> getCertificateData(final String jsonData) {
		try {
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

		return encodedPemItems.toArray(new String[encodedPemItems.size()]);
	}

	public void setTrafficOpsUtils(final TrafficOpsUtils trafficOpsUtils) {
		this.trafficOpsUtils = trafficOpsUtils;
	}
}
