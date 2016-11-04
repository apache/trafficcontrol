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

package com.comcast.cdn.traffic_control.traffic_router.secure;

import com.comcast.cdn.traffic_control.traffic_router.shared.CertificateData;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class CertificateRegistry {
	protected static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(CertificateRegistry.class);

	private CertificateDataConverter certificateDataConverter = new CertificateDataConverter();
	private Map<String, HandshakeData>	handshakeDataMap = new HashMap<>();

	// Recommended Singleton Pattern implementation
	// https://community.oracle.com/docs/DOC-918906
	private CertificateRegistry() {
	}

	public static CertificateRegistry getInstance() {
		return CertificateRegistryHolder.DELIVERY_SERVICE_CERTIFICATES;
	}

	public List<String> getAliases() {
		synchronized (handshakeDataMap) {
			return new ArrayList<>(handshakeDataMap.keySet());
		}
	}

	public HandshakeData getHandshakeData(final String alias) {
		synchronized (handshakeDataMap) {
			return handshakeDataMap.get(alias);
		}
	}

	@SuppressWarnings("PMD.AccessorClassGeneration")
	private static class CertificateRegistryHolder {
		private static final CertificateRegistry DELIVERY_SERVICE_CERTIFICATES = new CertificateRegistry();
	}

	public void importCertificateDataList(final List<CertificateData> certificateDataList) {
		final Map<String, HandshakeData> map = new HashMap<>();
		for (final CertificateData certificateData : certificateDataList) {
			final HandshakeData handshakeData = certificateDataConverter.toHandshakeData(certificateData);
			final String alias = handshakeData.getHostname().replaceFirst("\\*\\.", "");
			log.warn("Imported handshake data with alias " + alias);
			map.put(alias, handshakeData);
		}

		synchronized (handshakeDataMap) {
			handshakeDataMap = map;
		}
	}

	public CertificateDataConverter getCertificateDataConverter() {
		return certificateDataConverter;
	}

	public void setCertificateDataConverter(final CertificateDataConverter certificateDataConverter) {
		this.certificateDataConverter = certificateDataConverter;
	}
}
