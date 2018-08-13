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

import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtils;
import com.comcast.cdn.traffic_control.traffic_router.shared.CertificateData;
import com.fasterxml.jackson.databind.JsonNode;
import org.apache.log4j.Logger;

import java.util.Iterator;
import java.util.List;

public class CertificateChecker {
	private final static Logger LOGGER = Logger.getLogger(CertificateChecker.class);

	public String getDeliveryServiceType(final JsonNode deliveryServiceJson) {
		final JsonNode matchsets = deliveryServiceJson.get("matchsets");
		return getDeliveryServiceTypeFromJson(matchsets);
	}

	public String getDeliveryServiceType(final DeliveryService deliveryService) {
		final JsonNode matchsets = deliveryService.getMatchsets();
		return getDeliveryServiceTypeFromJson(matchsets);
	}

	private String getDeliveryServiceTypeFromJson(final JsonNode matchsets) {

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

	public boolean certificatesAreValid(final List<CertificateData> certificateDataList, final List<DeliveryService> deliveryServices) {

		boolean validConfig = true;
		final Iterator<DeliveryService> deliveryServiceIdIter = deliveryServices.iterator();
		while (deliveryServiceIdIter.hasNext()) {
			if (!deliveryServiceHasValidCertificates(certificateDataList, deliveryServiceIdIter.next())) {
				validConfig = false; // individual DS errors are logged in deliveryServiceHasValidCertificates
			}
		}
		return validConfig;
	}

	public boolean hasCertificate(final List<CertificateData> certificateDataList, final String deliveryServiceId) {
		return certificateDataList.stream()
			.filter(cd -> cd.getDeliveryservice().equals(deliveryServiceId))
			.findFirst()
			.isPresent();
	}

	@SuppressWarnings("PMD.CyclomaticComplexity")
	private Boolean deliveryServiceHasValidCertificates(final List<CertificateData> certificateDataList, final DeliveryService deliveryService) {
		final String deliveryServiceId = deliveryService.getId();

		if (!supportsHttps(deliveryService)) {
			return true;
		}

		final String domain = deliveryService.getDomain();

		if (domain == null) {
			LOGGER.warn("Delivery service " + deliveryServiceId + " is not configured with any domains!");
			return true;
		}

		for (final CertificateData certificateData : certificateDataList) {
			final String certificateDeliveryServiceId = certificateData.getDeliveryservice();
			if ((deliveryServiceId == null) || deliveryServiceId.equals("")) {
				LOGGER.error("DeliveryService name is blank for hostname '" +  certificateData.getHostname() + "', skipping.");
			} else if ((certificateDeliveryServiceId != null) && (deliveryServiceId != null) && (certificateDeliveryServiceId.equals(deliveryServiceId))) {
				LOGGER.debug("DService " + deliveryServiceId + " has certificate data for https");
				return true;
			}
		}

		LOGGER.error("No certificate data for https " + deliveryServiceId + " domain " + domain);
		return false;
	}

	private boolean supportsHttps(final DeliveryService deliveryService) {
		if (!"HTTP".equals(getDeliveryServiceType(deliveryService))) {
			return false;
		}

		return deliveryService.isAcceptHttps();
	}
}
