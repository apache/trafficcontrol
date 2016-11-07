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

package com.comcast.cdn.traffic_control.traffic_router.core.config;

import com.comcast.cdn.traffic_control.traffic_router.shared.CertificateData;
import org.apache.log4j.Logger;
import org.json.JSONArray;
import org.json.JSONObject;

import java.util.List;

public class CertificateChecker {
	private final static Logger LOGGER = Logger.getLogger(CertificateChecker.class);

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

	public boolean certificatesAreValid(final List<CertificateData> certificateDataList, final JSONObject deliveryServicesJson) {
		for (final String deliveryServiceId : JSONObject.getNames(deliveryServicesJson)) {
			if (!deliveryServiceHasValidCertificates(certificateDataList, deliveryServicesJson, deliveryServiceId)) {
				return false;
			}
		}
		return true;
	}

	public boolean hasCertificate(final List<CertificateData> certificateDataList, final String deliveryServiceId) {
		return certificateDataList.stream()
			.filter(cd -> cd.getDeliveryservice().equals(deliveryServiceId))
			.findFirst()
			.isPresent();
	}

	private Boolean deliveryServiceHasValidCertificates(final List<CertificateData> certificateDataList, final JSONObject deliveryServicesJson, final String deliveryServiceId) {
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

		if (domains.length() == 0) {
			return true;
		}

		boolean hasValidCertificates = false;

		for (int i = 0; i < domains.length(); i++) {
			final String domain = domains.optString(i, "").replaceAll("^\\*\\.", "");
			if (domain == null || domain.isEmpty()) {
				continue;
			}

			for (final CertificateData certificateData : certificateDataList) {
				if (certificateData.getDeliveryservice().equals(deliveryServiceId)) {
					hasValidCertificates = true;
				}
			}
			LOGGER.error("No certificate data for https " + deliveryServiceId + " domain " + domain);
		}

		return hasValidCertificates;
	}

	private boolean supportsHttps(final JSONObject deliveryServiceJson, final JSONObject protocolJson) {
		if (!"HTTP".equals(getDeliveryServiceType(deliveryServiceJson))) {
			return false;
		}

		return protocolJson != null ? protocolJson.optBoolean("acceptHttps", false) : false;
	}
}
