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

package org.apache.traffic_control.traffic_router.core.config;

import org.apache.traffic_control.traffic_router.core.util.JsonUtils;
import org.apache.traffic_control.traffic_router.shared.CertificateData;
import com.fasterxml.jackson.databind.JsonNode;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.util.Iterator;
import java.util.List;

public class CertificateChecker {
	private final static Logger LOGGER = LogManager.getLogger(CertificateChecker.class);

	public String getDeliveryServiceType(final JsonNode deliveryServiceJson) {
		final JsonNode matchsets = deliveryServiceJson.get("matchsets");

		for (final JsonNode matchset : matchsets) {
			if (matchset == null) {
				continue;
			}

			final String deliveryServiceType = JsonUtils.optString(matchset, "protocol");
			if (!deliveryServiceType.isEmpty()) {
				return deliveryServiceType;
			}
		}
		return null;
	}

	public boolean certificatesAreValid(final List<CertificateData> certificateDataList, final JsonNode deliveryServicesJson) {

		final Iterator<String> deliveryServiceIdIter = deliveryServicesJson.fieldNames();
		boolean invalidConfig = false;

		while (deliveryServiceIdIter.hasNext()) {
			if (!deliveryServiceHasValidCertificates(certificateDataList, deliveryServicesJson, deliveryServiceIdIter.next())) {
				invalidConfig = true; // individual DS errors are logged when deliveryServiceHasValidCertificates() is called
			}
		}

		if (invalidConfig) {
			return false;
		}

		return true;
	}

	public boolean hasCertificate(final List<CertificateData> certificateDataList, final String deliveryServiceId) {
		return certificateDataList.stream()
			.filter(cd -> cd.getDeliveryservice().equals(deliveryServiceId))
			.findFirst()
			.isPresent();
	}

    @SuppressWarnings("PMD.CyclomaticComplexity")
	private boolean deliveryServiceHasValidCertificates(final List<CertificateData> certificateDataList, final JsonNode deliveryServicesJson, final String deliveryServiceId) {
		final JsonNode deliveryServiceJson = deliveryServicesJson.get(deliveryServiceId);
		final JsonNode protocolJson = deliveryServiceJson.get("protocol");

		if (!supportsHttps(deliveryServiceJson, protocolJson)) {
			return true;
		}

		final JsonNode domains = deliveryServiceJson.get("domains");

		if (domains == null) {
			LOGGER.warn("Delivery Service " + deliveryServiceId + " is not configured with any domains!");
			return true;
		}

		if (domains.size() == 0) {
			return true;
		}

		for (final JsonNode domain : domains) {
			final String domainStr = domain.asText("").replaceAll("^\\*\\.", "");
			if (domainStr == null || domainStr.isEmpty()) {
				continue;
			}

			for (final CertificateData certificateData : certificateDataList) {
				final String certificateDeliveryServiceId = certificateData.getDeliveryservice();
				if ((deliveryServiceId == null) || deliveryServiceId.equals("")) {
					LOGGER.error("Delivery Service name is blank for hostname '" +  certificateData.getHostname() + "', skipping.");
				} else if ((certificateDeliveryServiceId != null) && (deliveryServiceId != null) && (certificateDeliveryServiceId.equals(deliveryServiceId))) {
					LOGGER.debug("Delivery Service " + deliveryServiceId + " has certificate data for https");
					return true;
				}
			}
			LOGGER.error("No certificate data for https " + deliveryServiceId + " domain " + domainStr);
		}

		return false;
	}

	private boolean supportsHttps(final JsonNode deliveryServiceJson, final JsonNode protocolJson) {
		if (!"HTTP".equals(getDeliveryServiceType(deliveryServiceJson))) {
			return false;
		}

		return JsonUtils.optBoolean(protocolJson, "acceptHttps");
	}
}
