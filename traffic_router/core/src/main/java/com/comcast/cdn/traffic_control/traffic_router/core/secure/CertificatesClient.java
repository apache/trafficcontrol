package com.comcast.cdn.traffic_control.traffic_router.core.secure;

import com.comcast.cdn.traffic_control.traffic_router.core.util.ProtectedFetcher;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyStoreHelper;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.log4j.Logger;

import java.nio.file.Files;
import java.nio.file.Paths;
import java.nio.file.attribute.FileTime;
import java.util.ArrayList;
import java.util.Base64;
import java.util.List;

public class CertificatesClient {
	private static final Logger LOGGER = Logger.getLogger(CertificatesClient.class);
	private TrafficOpsUtils trafficOpsUtils;
	private static final String PEM_HEADER_PREFIX = "-----BEGIN";
	private static final String PEM_FOOTER_PREFIX = "-----END";

	public void refreshData() {
		final String jsonData = fetchRawData();

		if (jsonData == null) {
			return;
		}

		final List<CertificateData> certificateDataList = getCertificateData(jsonData);

		if (certificateDataList.isEmpty()) {
			return;
		}

		persistCertificates(certificateDataList);
	}

	public String fetchRawData() {
		final String certificatesUrl = trafficOpsUtils.getUrl("certificate.api.url", "https://${toHostname}/api/1.2/cdns/name/${cdnName}/sslkeys.json");

		try {
			final ProtectedFetcher fetcher = new ProtectedFetcher(trafficOpsUtils.getAuthUrl(), trafficOpsUtils.getAuthJSON().toString(), 15000);
			final FileTime fileTime = Files.getLastModifiedTime(Paths.get(getKeystorePath()));
			return fetcher.fetchIfModifiedSince(certificatesUrl, fileTime.toMillis());
		} catch (Exception e) {
			LOGGER.warn("Failed to fetch data for certificates from " + certificatesUrl + " : " + e.getMessage());
		}

		return null;
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
				continue;
			}

			if (line.startsWith(PEM_HEADER_PREFIX)) {
				continue;
			}

			builder.append(line);
		}

		return encodedPemItems.toArray(new String[encodedPemItems.size()]);
	}

	public boolean persistCertificates(final List<CertificateData> certificateDataList) {
		boolean allCertificatesPersisted = true;
		final KeyStoreHelper keyStoreHelper = KeyStoreHelper.getInstance();

		keyStoreHelper.clearCertificates();

		for (final CertificateData certificateData : certificateDataList) {
			final String alias = certificateData.getDeliveryservice();
			final String key = doubleDecode(certificateData.getCertificate().getKey())[0];
			final String[] chain = doubleDecode(certificateData.getCertificate().getCrt());

			if (!keyStoreHelper.importCertificateChain(alias, key, chain)) {
				allCertificatesPersisted = false;
			} else {
				LOGGER.info("Persisted certificate for alias '" + alias + "'");
			}
		}

		keyStoreHelper.save();
		return allCertificatesPersisted;
	}

	public String getKeystorePath() {
		return System.getProperty("deploy.dir", "/opt/traffic_router") + "/db/.keystore";
	}

	public void setTrafficOpsUtils(final TrafficOpsUtils trafficOpsUtils) {
		this.trafficOpsUtils = trafficOpsUtils;
	}
}
