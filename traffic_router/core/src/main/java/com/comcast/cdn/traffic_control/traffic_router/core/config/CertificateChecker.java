package com.comcast.cdn.traffic_control.traffic_router.core.config;

import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyStoreHelper;
import org.apache.log4j.Logger;
import org.json.JSONArray;
import org.json.JSONObject;

public class CertificateChecker {
	private final static Logger LOGGER = Logger.getLogger(CertificateChecker.class);
	private final KeyStoreHelper keyStoreHelper = KeyStoreHelper.getInstance();

	public String getDeliveryServiceType(final JSONObject deliveryServiceJson) {
		final JSONArray matchsets = deliveryServiceJson.optJSONArray("matchsets");

		for (int i = 0; i < matchsets.length(); i++) {
			final JSONObject matchset = matchsets.optJSONObject(i);
			if (matchset == null) {
				continue;
			}

			final String deliveryServiceType = matchset.optString("protocol", "");
			if (!deliveryServiceType.isEmpty()) {
				return deliveryServiceType;
			}
		}
		return null;
	}

	public boolean certificatesAreValid(final JSONObject deliveryServicesJson) {
		for (final String deliveryServiceId : JSONObject.getNames(deliveryServicesJson)) {
			if (!deliveryServiceHasValidCertificates(deliveryServicesJson, deliveryServiceId)) {
				return false;
			}
		}
		return true;
	}

	private Boolean deliveryServiceHasValidCertificates(final JSONObject deliveryServicesJson, final String deliveryServiceId) {
		final JSONObject deliveryServiceJson = deliveryServicesJson.optJSONObject(deliveryServiceId);
		final JSONObject protocolJson = deliveryServiceJson.optJSONObject("protocol");

		if (!supportsHttps(deliveryServiceJson, protocolJson)) {
			return true;
		}

		final JSONArray domains = deliveryServiceJson.optJSONArray("domains");

		if (domains == null) {
			LOGGER.warn("Delivery Service " + deliveryServiceId + " is not configured with any domains!");
			return true;
		}

		for (int i = 0; i < domains.length(); i++) {
			final String domain = domains.optString(i, "").replaceAll("^\\*\\.", "");
			if (domain == null || domain.isEmpty()) {
				continue;
			}

			if (!keyStoreHelper.hasCertificate(domain)) {
				LOGGER.error("Delivery Service " + deliveryServiceId + " with domain " + domain + " is marked to accept https traffic and does not have a certificate");
				return false;
			}
		}

		return true;
	}

	private boolean supportsHttps(final JSONObject deliveryServiceJson, final JSONObject protocolJson) {
		if (!"HTTP".equals(getDeliveryServiceType(deliveryServiceJson))) {
			return false;
		}

		return protocolJson != null ? protocolJson.optBoolean("acceptHttps", false) : false;
	}
}
