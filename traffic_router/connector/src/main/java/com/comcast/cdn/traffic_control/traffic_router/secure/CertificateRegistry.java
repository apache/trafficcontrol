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
import sun.security.tools.keytool.CertAndKeyGen;
import sun.security.util.ObjectIdentifier;
import sun.security.x509.CertificateExtensions;
import sun.security.x509.BasicConstraintsExtension;
import sun.security.x509.KeyUsageExtension;
import sun.security.x509.ExtendedKeyUsageExtension;
import sun.security.x509.X500Name;

import java.security.PrivateKey;
import java.security.cert.X509Certificate;
import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Vector;

public class CertificateRegistry {
	public static final String DEFAULT_SSL_KEY = "default.invalid";
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
			final Map<String, HandshakeData> handshakeDataMap =
					CertificateRegistryHolder.DELIVERY_SERVICE_CERTIFICATES.getHandshakeData();
			final HandshakeData defaultHd = createDefaultSsl();
			if (defaultHd == null) {
				log.error("Failed to initialize the CertificateRegistry.");
				return null;
			}
			handshakeDataMap.putIfAbsent(DEFAULT_SSL_KEY, defaultHd);
		}
		return CertificateRegistryHolder.DELIVERY_SERVICE_CERTIFICATES;
	}

	@SuppressWarnings("PMD.UseArrayListInsteadOfVector")
	private static HandshakeData createDefaultSsl() {
		try {
			/*final String DEFAULT_CERT =
					"    {\n" +
							"      \"deliveryservice\": \""+DEFAULT_SSL_KEY+"\",\n" +
							"      \"certificate\": {\n" +
							"        \"comment\" : \"The following is valid default \",\n" +
							"        \"key\": " +
							"\"LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRQzNNYmt3ZTFyLy9HMVQKaEVSR3JMWFNvK1B3Uisrc1h0Zzk0dk00VWhnelRUQzl6bWd5Y2NPK3RTTERhOHRmOVlsSEdvb25NZVBnUndvNgpYRWZnQmdjMllyR0ErYnJyU08vVVBUNzdhUkN5ajB4UUhmQXdqZVNlYStNdEN2Nk5CdFdiV1Bua24yalpGRk1sCkJyRlZ2RkpHc08zYjdHbVI2NjlWSkIrSXl0dlpkWW1kek9zSUlzKzQ2ZFNsNm5tcVdjWWZZa1NlYld0b0Y0SEYKODJoSGlkOStNUGl1Q2ozd2NqbExMNlNVbWtiWmdEODFYbGVuQUZXQlR1eFV0QzJWVzVsVElMMUliSVM1SSt5MwpxelA2dFg1ZktVZlJGVjBBM0hPSGpZcVJWR2Z4QXh3Qy9qanpJQlYvSllwamNGQWk2Qzh4Q0pGMzdwVWtlOVExCldjb0VGeU12QWdNQkFBRUNnZ0VCQUk0aHJ2dlZpU241RUYxdXpvWkM4NkxrOHpGMnJwWiswNmxZVHJwUXYyUDIKTEszbTJlTGhieXlrWHI5ZC8rR0lvQ1NoaTdTVE9hakZsVUxvVy8rTXpjVzlWdGlwYVFPcGlDR1VEeXlDUEtrOQpFc2xLSVJPYTAxaXlmZ1J4ZGtPMm5MNDFqMVI0OWFFTzZ0OWNUUFFtODNMVFRRaUhhUFVFOWZqSjJRbUowbjdwCkowbHFsTFpnaUxFQlloWWtOR1ZMbGppZzN6U3puWUNtSnpDRmRVYjhzTFBGeTlLdXZvdWoxSnRnWC9BUjFSQncKcDduVkpaK3N3MG5FZk9EVE0yZ2ZwNVZzaGV5d3c0ZzFIci9JbXdhTFdjSGpaWlYrbURsQ3hIZDZaL3dTcW9uQQpBTGZacnBhbi9KTm1yamE3NHN5UUNRQjQ4dUYxNmorSE5taFg0RUpPckNFQ2dZRUE4czc4YmZ6Mzg0aG1GTzJSCnRianhOTVFlbUtwWGxPN2R3Wm1DTHZ2amRMN09Qb2ZPQnNEdjV1eTI5KzlVN0JDV053RzY5bjVvNGp0cE5QREgKV252VklQVnEyWHhwY01oWlZLaGpDWTM1a1Q0TWRWSXhlODArdm1yMFV3TGExa0ozajVrZlRlWmZnRnZFNFhUWQo1MS9qYVB2dHJUOWttQlh3QTJPVXdSUEZTaEVDZ1lFQXdTV2REdWVMckJoZ3UyWE9ZNEFsemc4V3M5Vzk4MmxuCk9abzRqZG5yenBRZFpXcnpmamZvL1FVQk1HUkdJUm5MRUF4aUdONytBemtNeEhMa2xPVHEwWWtTN2ZkSmRQNncKWGs2WGRBMDJzUEYxVUtRVFhRMGRtd3B0eWxEa1hML3loL3FOSUliZHh4YzFGMmdXcng2Z1JnVm9DL0tjNlJFWQpyL0NnUU5rcVdUOENnWUFUTEVNRWtHQW42OUpidnJLdHpjL0dJZUprbmJiU3ZOWG43cTQzOVkzdGJ3K3NJbDhqCmEyTEdNbFQwV1FLMHJVNmZRMVMzR0I1Q0Z2emt3RXFObTQrbHpadEZWeXlnU2tHN2pKeGRhY2VXTDNjZVlJSWwKeTN3ejN4QXg2ZHpMNUcyNmVoWGR1ZDQ2cllSclpTV25oNHZXZzJZdU12NUhnQnYydUl0TGY3c3BjUUtCZ0dxUQoxVElIQU9JbjVSOGdFWnFHZHRWVkwrSnpHTVczTHhQeUNpZ0J4NEFINnB3dFFVRXZtZVlZSDhyU1dIc2szd3Z3CnVTTWR6YXA3akpiTENXRTVXSEhabms4YmREVVAzTUY3dlVaemoreGFuSzZzaUY1N3dRenMyUnlhT3hVTmRzUWQKc2tYekEyUTRZcnVTVzRtdGJTS1ZFdzRjZ3dSNHdWVTVmMEdvVUJ4REFvR0FEVzRWaldjZWRTV1RvUWFMVzNYTgpuTXF1dEx1MG4rVEZMeWg5KzZzbHFvKytKZmQ0aWJMMXRGZ0FaOVQ5dkx3Z1FPZncwVldhYWdyRDNUKzFiYW1iClJzeFFDWDZlTFh4VkZaTmsrb3F6ZVF0L0JBQWI0MHZJY2R2V1V4dk1FbW5taHdWVTlRWUlMaVVxa0JPRnJMVzYKSlJKSWJoZGd6MGdHdnVoR3R6UXJ6Y009Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K\",\n" +
							"        \"crt\": " +
							"\"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURjakNDQWxvQ0NRRENDbW5uQ0dtd2N6QU5CZ2txaGtpRzl3MEJBUXNGQURCN01Rc3dDUVlEVlFRR0V3SlYKVXpFTE1Ba0dBMVVFQ0F3Q1EwOHhEREFLQmdOVkJBY01BMFJGVGpFUE1BMEdBMVVFQ2d3R1FYQmhZMmhsTVFzdwpDUVlEVlFRTERBSlVRekVWTUJNR0ExVUVBd3dNWkdWbVlYVnNkQzVqWlhKME1Sd3dHZ1lKS29aSWh2Y05BUWtCCkZnMTBZMEJoY0dGamFHVXViM0puTUI0WERURTVNRE14T1RFM01EVXlNbG9YRFRJNU1ETXhOakUzTURVeU1sb3cKZXpFTE1Ba0dBMVVFQmhNQ1ZWTXhDekFKQmdOVkJBZ01Ba05QTVF3d0NnWURWUVFIREFORVJVNHhEekFOQmdOVgpCQW9NQmtGd1lXTm9aVEVMTUFrR0ExVUVDd3dDVkVNeEZUQVRCZ05WQkFNTURHUmxabUYxYkhRdVkyVnlkREVjCk1Cb0dDU3FHU0liM0RRRUpBUllOZEdOQVlYQmhZMmhsTG05eVp6Q0NBU0l3RFFZSktvWklodmNOQVFFQkJRQUQKZ2dFUEFEQ0NBUW9DZ2dFQkFMY3h1VEI3V3YvOGJWT0VSRWFzdGRLajQvQkg3NnhlMkQzaTh6aFNHRE5OTUwzTwphREp4dzc2MUlzTnJ5MS8xaVVjYWlpY3g0K0JIQ2pwY1IrQUdCelppc1lENXV1dEk3OVE5UHZ0cEVMS1BURkFkCjhEQ041SjVyNHkwSy9vMEcxWnRZK2VTZmFOa1VVeVVHc1ZXOFVrYXc3ZHZzYVpIcnIxVWtINGpLMjlsMWlaM00KNndnaXo3anAxS1hxZWFwWnhoOWlSSjV0YTJnWGdjWHphRWVKMzM0dytLNEtQZkJ5T1VzdnBKU2FSdG1BUHpWZQpWNmNBVllGTzdGUzBMWlZibVZNZ3ZVaHNoTGtqN0xlck0vcTFmbDhwUjlFVlhRRGNjNGVOaXBGVVovRURIQUwrCk9QTWdGWDhsaW1Od1VDTG9MekVJa1hmdWxTUjcxRFZaeWdRWEl5OENBd0VBQVRBTkJna3Foa2lHOXcwQkFRc0YKQUFPQ0FRRUFET3hPbGgzWXl5NHFNQ0Q1YzZRaXZ0SnN6ZkZhcDV3eG9weTJwS0tpdFRXWVBEUVRKcXJ0dnNPZwphQ3d5L3NHVElNcHF6SFdwamxUUDhWdHJwSllLdVp0dXNwcDIzd0xTU0JaVUJobXEveW9OZTM0WVRJQXZncVhGCm5QZWtkWFFUdDVrUU9uUGgyS2N5WVBLdWxQMkRmc1JYWXJOVSsweEpZUFQ1bGdIcXFteklONzVYRTdIWDErMWcKSUxtdjFibmdzQmxDRFk2Unl3YmRrak9obUg0OU0rbGFCdXJCbDlDdjNqTEhvampncWNiOTg5RkZQRjdpdGtqNwpBSkxHK0NiV2tkckZmY0pwUURrVC9TaDdvRkhiZmhLRWoyR2ZaZmkrc3QxNlhTaWJxaXVZOWdRclhLRFVtMk9QCnBuK0xudllNT3N2Uld4Q0pLQXVjMERta3FLOWlJUT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K\"\n"+
							"      },\n" +
							"      \"hostname\": \""+DEFAULT_SSL_KEY+"\"\n" +
							"    }";

			final CertificateDataConverter certificateDataConverter = new CertificateDataConverter();
			final CertificateData	certificateData = ((CertificateData) new ObjectMapper().readValue(DEFAULT_CERT,
						new TypeReference<CertificateData>() { }));
			return certificateDataConverter.toHandshakeData(certificateData);
			*/

			final CertificateExtensions extensions = new CertificateExtensions();
			final KeyUsageExtension keyUsageExtension = new KeyUsageExtension();
			keyUsageExtension.set(KeyUsageExtension.DIGITAL_SIGNATURE, true);
			keyUsageExtension.set(KeyUsageExtension.KEY_ENCIPHERMENT, true);
			keyUsageExtension.set(KeyUsageExtension.KEY_CERTSIGN, true);
			extensions.set(keyUsageExtension.getExtensionId().toString(), keyUsageExtension);
			final Vector<ObjectIdentifier> objectIdentifiers = new Vector<>();
			objectIdentifiers.add(new ObjectIdentifier("1.3.6.1.5.5.7.3.1"));
			objectIdentifiers.add(new ObjectIdentifier("1.3.6.1.5.5.7.3.2"));
			final ExtendedKeyUsageExtension extendedKeyUsageExtension = new ExtendedKeyUsageExtension( true,
					objectIdentifiers);
			extensions.set(extendedKeyUsageExtension.getExtensionId().toString(), extendedKeyUsageExtension);
			extensions.set(BasicConstraintsExtension.NAME, new BasicConstraintsExtension(true,
					new BasicConstraintsExtension(true,-1).getExtensionValue()));
			final CertAndKeyGen certGen = new CertAndKeyGen("RSA", "SHA1WithRSA", null);
			certGen.generate(1024);

			//Generate self signed certificate
			final X509Certificate[] chain = new X509Certificate[1];
			chain[0] = certGen.getSelfCertificate(new X500Name("C=US; ST=PA; L=Philadelphia; " +
					"O=Comcast Corporation; OU=Comcast; OU=Hosted by Comcast Corporation; " +
					"OU=CDNDefault; CN="+DEFAULT_SSL_KEY), new Date(System.currentTimeMillis() - 1000L * 60 ),
					(long) 3 * 365 * 24 * 3600, extensions);
			final PrivateKey serverPrivateKey = certGen.getPrivateKey();
			return new HandshakeData(DEFAULT_SSL_KEY, DEFAULT_SSL_KEY, chain, serverPrivateKey);
		}
		catch (Exception e) {
			log.error("Could not generate the default certificate: "+e.getMessage(),e);
			return null;
		}
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

		master.putIfAbsent(DEFAULT_SSL_KEY, handshakeDataMap.get(DEFAULT_SSL_KEY));
		handshakeDataMap = master;

		if (sslEndpoint != null) {
			sslEndpoint.replaceSSLHosts(changes);
		}
	}

	public CertificateDataConverter getCertificateDataConverter() {
		return certificateDataConverter;
	}

	public void setCertificateDataConverter(final CertificateDataConverter certificateDataConverter) {
		this.certificateDataConverter = certificateDataConverter;
	}
}
