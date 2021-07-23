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

package org.apache.traffic_control.traffic_router.secure;

import org.apache.traffic_control.traffic_router.protocol.RouterNioEndpoint;
import org.apache.traffic_control.traffic_router.shared.CertificateData;
import org.apache.traffic_control.traffic_router.utils.HttpsProperties;
import org.apache.log4j.Logger;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x509.BasicConstraints;
import org.bouncycastle.asn1.x509.ExtendedKeyUsage;
import org.bouncycastle.asn1.x509.KeyPurposeId;
import org.bouncycastle.asn1.x509.KeyUsage;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;
import org.bouncycastle.asn1.x509.Extension;

import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.InputStream;
import java.math.BigInteger;
import java.net.InetAddress;
import java.security.KeyPairGenerator;
import java.security.KeyPair;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.Security;
import java.security.cert.Certificate;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class CertificateRegistry {
	public static final String DEFAULT_SSL_KEY = "default.invalid";
	private static final Logger log = Logger.getLogger(CertificateRegistry.class);
	private CertificateDataConverter certificateDataConverter = new CertificateDataConverter();
	volatile private Map<String, HandshakeData> handshakeDataMap = new HashMap<>();
	private RouterNioEndpoint sslEndpoint = null;
	final private Map<String, CertificateData> previousData = new HashMap<>();
	public String defaultAlias;

	// Recommended Singleton Pattern implementation
	// https://community.oracle.com/docs/DOC-918906
	private CertificateRegistry() {
		try {
			defaultAlias = InetAddress.getLocalHost().getHostName();
		} catch (Exception e) {
			log.error("Error getting hostname");
		}
	}

	public static CertificateRegistry getInstance() {
		return CertificateRegistryHolder.DELIVERY_SERVICE_CERTIFICATES;
	}

	@SuppressWarnings({"PMD.UseArrayListInsteadOfVector", "PMD.AvoidUsingHardCodedIP"})
	private static HandshakeData createDefaultSsl() {
		try {
			final KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
			keyPairGenerator.initialize(2048);
			final KeyPair keyPair = keyPairGenerator.generateKeyPair();

			//Generate self signed certificate
			final X509Certificate[] chain = new X509Certificate[1];

			// Select provider
			Security.addProvider(new BouncyCastleProvider());

			// Generate cert details
			final long now = System.currentTimeMillis();
			final Date startDate = new Date(System.currentTimeMillis());

			final X500Name dnName = new X500Name("C=US; ST=CO; L=Denver; " +
					"O=Apache Traffic Control; OU=Apache Foundation; OU=Hosted by Traffic Control; " +
					"OU=CDNDefault; CN="+DEFAULT_SSL_KEY);
			final BigInteger certSerialNumber = new BigInteger(Long.toString(now));

			final Calendar calendar = Calendar.getInstance();
			calendar.setTime(startDate);
			calendar.add(Calendar.YEAR, 3);

			final Date endDate = calendar.getTime();

			// Build certificate
			final ContentSigner contentSigner = new JcaContentSignerBuilder("SHA1WithRSA").build(keyPair.getPrivate());

			final JcaX509v3CertificateBuilder certBuilder = new JcaX509v3CertificateBuilder(dnName, certSerialNumber, startDate, endDate, dnName, keyPair.getPublic());

			// Attach extensions
			certBuilder.addExtension(Extension.basicConstraints, true, new BasicConstraints(true));
			certBuilder.addExtension(Extension.keyUsage, true, new KeyUsage(KeyUsage.digitalSignature | KeyUsage.keyEncipherment | KeyUsage.keyCertSign));
			certBuilder.addExtension(Extension.extendedKeyUsage, true, new ExtendedKeyUsage(new KeyPurposeId[] {
					KeyPurposeId.id_kp_clientAuth,
					KeyPurposeId.id_kp_serverAuth
			}));

			// Generate final certificate
			final X509CertificateHolder certHolder = certBuilder.build(contentSigner);

			final JcaX509CertificateConverter converter = new JcaX509CertificateConverter();
			converter.setProvider(new BouncyCastleProvider());
			chain[0] = converter.getCertificate(certHolder);

			return new HandshakeData(DEFAULT_SSL_KEY, DEFAULT_SSL_KEY, chain, keyPair.getPrivate());
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

	private HandshakeData createApiDefaultSsl() {
		try {
			final Map<String, String> httpsProperties = (new HttpsProperties()).getHttpsPropertiesMap();

			final KeyStore ks = KeyStore.getInstance("JKS");
			final String selfSignedKeystoreFile = httpsProperties.get("https.certificate.location");
			if (new File(selfSignedKeystoreFile).exists()) {
				final String password = httpsProperties.get("https.password");
				final InputStream readStream = new FileInputStream(selfSignedKeystoreFile);
				ks.load(readStream, password.toCharArray());
				readStream.close();
				final Certificate[] certs = ks.getCertificateChain(defaultAlias);
				final List<X509Certificate> x509certs = new ArrayList<>();

				for (final Certificate cert : certs) {
					final CertificateFactory cf = CertificateFactory.getInstance("X.509");
					final ByteArrayInputStream bais = new ByteArrayInputStream(cert.getEncoded());
					final X509Certificate x509cert = (X509Certificate) cf.generateCertificate(bais);
					x509certs.add(x509cert);
				}

				X509Certificate[] x509CertsArray = new X509Certificate[x509certs.size()];
				x509CertsArray = x509certs.toArray(x509CertsArray);

				final HandshakeData handshakeData = new HandshakeData(defaultAlias, defaultAlias,
						x509CertsArray, (PrivateKey) ks.getKey(defaultAlias, password.toCharArray()));

				return handshakeData;
			}
		} catch (Exception e) {
			log.error("Failed to load default certificate. Received " + e.getClass() + " with message: " + e.getMessage());
			return null;
		}
		return null;
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
		for (final String alias : previousData.keySet()) {
			if (!master.containsKey(alias) && sslEndpoint != null) {
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

		// Check to see if a Default cert has been provided by Traffic Ops
		if (!master.containsKey(DEFAULT_SSL_KEY)){
			// Check to see if a Default cert has been provided/created previously
			if (handshakeDataMap.containsKey(DEFAULT_SSL_KEY)) {
				master.put(DEFAULT_SSL_KEY, handshakeDataMap.get(DEFAULT_SSL_KEY));
			}else{
				// create a new default certificate
				final HandshakeData defaultHd = createDefaultSsl();
				if (defaultHd == null){
					log.error("Failed to initialize the CertificateRegistry because of a problem with the 'default' " +
							"certificate. Returning the Certificate Registry without a default.");
					return;
				}
				master.put(DEFAULT_SSL_KEY, defaultHd);
			}
		}

		if (!master.containsKey(defaultAlias)) {
			if (handshakeDataMap.containsKey(defaultAlias)) {
				master.put(defaultAlias, handshakeDataMap.get(defaultAlias));
			} else {
				final HandshakeData apiDefault = createApiDefaultSsl();
				if (apiDefault == null) {
					log.error("Failed to initialize the API Default certificate.");
				} else {
					master.put(apiDefault.getHostname(), apiDefault);
				}
			}
		}
		handshakeDataMap = master;

		// This will update the SSLHostConfig objects stored in the server
		// if any of those updates fail then we need to be sure to remove them
		// from the previousData list so that we will try to update them again
		// next time we import certificates
		if (sslEndpoint != null && !changes.isEmpty()) {
			final List<String> failedUpdates = sslEndpoint.reloadSSLHosts(changes);
			failedUpdates.forEach(alias-> {
				previousData.remove(alias);
			});
		}
	}

	public CertificateDataConverter getCertificateDataConverter() {
		return certificateDataConverter;
	}

	public void setCertificateDataConverter(final CertificateDataConverter certificateDataConverter) {
		this.certificateDataConverter = certificateDataConverter;
	}
}
