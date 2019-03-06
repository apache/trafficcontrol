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

import com.comcast.cdn.traffic_control.traffic_router.protocol.RouterNioEndpoint;
import com.comcast.cdn.traffic_control.traffic_router.shared.CertificateData;
import org.apache.log4j.Logger;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class CertificateRegistry {
	private static final Logger log = Logger.getLogger(CertificateRegistry.class);
	private CertificateDataConverter certificateDataConverter = new CertificateDataConverter();
	volatile private Map<String, HandshakeData>	handshakeDataMap = new HashMap<>();
	private RouterNioEndpoint sslEndpoint = null;
	final private Map<String, CertificateData> previousData = new HashMap<>();

	// Recommended Singleton Pattern implementation
	// https://community.oracle.com/docs/DOC-918906
	private CertificateRegistry() {
	}

	public static CertificateRegistry getInstance() {
		synchronized (CertificateRegistryHolder.DELIVERY_SERVICE_CERTIFICATES) {
			Map<String, HandshakeData> handshakeDataMap =
					CertificateRegistryHolder.DELIVERY_SERVICE_CERTIFICATES.getHandshakeData();
			handshakeDataMap.putIfAbsent("_default_", createDefaultSsl());
		}
		return CertificateRegistryHolder.DELIVERY_SERVICE_CERTIFICATES;
	}

	private static HandshakeData createDefaultSsl() {
		return new HandshakeData("", "", null,null);
	}

	public List<String> getAliases() {
		return new ArrayList<>(handshakeDataMap.keySet());
	}

	public HandshakeData getHandshakeData(final String alias) {
		return handshakeDataMap.get(alias);
	}

	public Map<String, HandshakeData> getHandshakeData() {
	    return handshakeDataMap;
    }

	public void setEndPoint(final RouterNioEndpoint routerNioEndpoint) {
		sslEndpoint = routerNioEndpoint;
	}

	@SuppressWarnings("PMD.AccessorClassGeneration")
	private static class CertificateRegistryHolder {
		private static final CertificateRegistry DELIVERY_SERVICE_CERTIFICATES = new CertificateRegistry();
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.AvoidDeeplyNestedIfStmts", "PMD.NPathComplexity"})
	synchronized public void importCertificateDataList(final List<CertificateData> certificateDataList) {
		final Map<String, HandshakeData> changes = new HashMap<>();
		final Map<String, HandshakeData> master = new HashMap<>();

		// find CertificateData which has changed
		for (final CertificateData certificateData : certificateDataList) {
			try {
				final String alias = certificateData.alias();

				if (!master.containsKey(alias)) {
					final HandshakeData handshakeData = certificateDataConverter.toHandshakeData(certificateData);
					if (handshakeData != null) {
						master.put(alias, handshakeData);
						if (!certificateData.equals(previousData.get(alias))) {
							changes.put(alias, handshakeData);
							log.warn("Imported handshake data with alias " + alias);
						}
					}
				}
				else {
					log.error("An TLS certificate already exists in the registry for host: "+alias+" There can be " +
							"only one!" );
				}
			} catch (Exception e) {
				log.error("Failed to import certificate data for delivery service: '" + certificateData.getDeliveryservice() + "', hostname: '" + certificateData.getHostname() + "'");
			}
		}

		// find CertificateData which has been removed
		for (final String alias : previousData.keySet())
		{
			if (!master.containsKey(alias) && sslEndpoint != null)
			{
				final String hostname = previousData.get(alias).getHostname();
				sslEndpoint.removeSslHostConfig(hostname);
			    log.warn("Removed handshake data with hostname " + hostname);
			}
		}

		// store the result for the next import
		previousData.clear();
		for (final CertificateData certificateData : certificateDataList) {
			final String alias = certificateData.alias();
			if (!previousData.containsKey(alias) && master.containsKey(alias)) {
				previousData.put(alias, certificateData);
			}
		}

		handshakeDataMap = master;

		if (sslEndpoint != null) {
			sslEndpoint.reloadSSLHosts(changes);
		}

	}

	public CertificateDataConverter getCertificateDataConverter() {
		return certificateDataConverter;
	}

	public void setCertificateDataConverter(final CertificateDataConverter certificateDataConverter) {
		this.certificateDataConverter = certificateDataConverter;
	}
}
