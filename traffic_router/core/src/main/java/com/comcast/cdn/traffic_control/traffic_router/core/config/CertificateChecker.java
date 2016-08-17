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
			final JSONObject deliveryServiceJson = deliveryServicesJson.optJSONObject(deliveryServiceId);
			final JSONObject protocolJson = deliveryServiceJson.optJSONObject("protocol");

			if (protocolJson == null) {
				continue;
			}

			if (!"HTTP".equals(getDeliveryServiceType(deliveryServiceJson))) {
				continue;
			}

			if (!protocolJson.optBoolean("acceptHttps", false)) {
				continue;
			}

			if (!keyStoreHelper.hasCertificate(deliveryServiceId)) {
				LOGGER.error("Delivery Service " + deliveryServiceId + " is marked to accept https traffic and does not have a certificate");
				return false;
			}
		}
		return true;
	}
}
