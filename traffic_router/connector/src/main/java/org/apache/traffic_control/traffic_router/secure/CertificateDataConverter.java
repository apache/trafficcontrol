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

import org.apache.traffic_control.traffic_router.shared.CertificateData;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.bouncycastle.jcajce.provider.asymmetric.rsa.BCRSAPrivateCrtKey;

import java.math.BigInteger;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.cert.CertificateExpiredException;
import java.security.cert.CertificateNotYetValidException;
import java.security.cert.X509Certificate;
import java.security.spec.RSAPrivateCrtKeySpec;
import java.security.spec.RSAPublicKeySpec;
import java.util.ArrayList;
import java.util.List;

@SuppressWarnings({"PMD.CyclomaticComplexity"})
public class CertificateDataConverter {
	private static final Logger log = LogManager.getLogger(CertificateDataConverter.class);

	private PrivateKeyDecoder privateKeyDecoder = new PrivateKeyDecoder();
	private CertificateDecoder certificateDecoder = new CertificateDecoder();

	@SuppressWarnings({"PMD.CyclomaticComplexity"})
	public HandshakeData toHandshakeData(final CertificateData certificateData) {
		try {
			final PrivateKey privateKey = privateKeyDecoder.decode(certificateData.getCertificate().getKey());
			final List<String> encodedCertificates = certificateDecoder.doubleDecode(certificateData.getCertificate().getCrt());

			final List<X509Certificate> x509Chain = new ArrayList<>();
			boolean hostMatch = false;
			boolean modMatch = false;
			for (final String encodedCertificate : encodedCertificates) {
				final X509Certificate certificate = certificateDecoder.toCertificate(encodedCertificate);
				certificate.checkValidity();
				if (!hostMatch && verifySubject(certificate, certificateData.alias())) {
					hostMatch = true;
				}
				if (!modMatch && verifyModulus(privateKey, certificate)) {
					modMatch = true;
				}
				x509Chain.add(certificate);
			}
			if (hostMatch && modMatch) {
				return new HandshakeData(certificateData.getDeliveryservice(), certificateData.getHostname(),
						x509Chain.toArray(new X509Certificate[0]), privateKey);
			}
			else if (!hostMatch) {
				log.warn("Service name doesn't match the subject of the certificate = "+certificateData.getHostname());
			}
			else if (!modMatch) {
				log.warn("Modulus of the private key does not match the public key modulus for certificate host: "+certificateData.getHostname());
			}

		} catch (CertificateNotYetValidException er) {
			log.warn("Failed to convert certificate data for delivery service = " + certificateData.getHostname()
							+ ", because the certificate is not valid yet. This certificate will not be used by " +
					"Traffic Router.");
		} catch (CertificateExpiredException ex ) {
			log.warn("Failed to convert certificate data for delivery service = " + certificateData.getHostname()
					+ ", because the certificate has expired. This certificate will not be used by Traffic Router.");
		} catch (Exception e) {
			log.warn("Failed to convert certificate data (delivery service = " + certificateData.getDeliveryservice()
					+ ", hostname = " + certificateData.getHostname() + ") from traffic ops to handshake data! This " +
					"certificate will not be used by Traffic Router. "
					+ e.getClass().getSimpleName() + ": " + e.getMessage(), e);
		}
		return null;
	}

	public boolean verifySubject(final X509Certificate certificate, final String hostAlias ) {
		final String host = certificate.getSubjectDN().getName();
		if (hostCompare(hostAlias,host)) {
			return true;
		}

		try {
			// This approach is probably the only one that is JDK independent
			if (certificate.getSubjectAlternativeNames() != null) {
				for (final List<?> altName : certificate.getSubjectAlternativeNames()) {
					if (hostCompare(hostAlias, (String) altName.get(1))) {
						return true;
					}
				}
			}
		}
		catch (Exception e) {
			log.error("Encountered an error while validating the certificate subject for service: "+hostAlias+", " +
					"error: "+e.getClass().getSimpleName()+": " + e.getMessage(), e);
			return false;
		}

		return false;
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity"})
	private boolean hostCompare(final String hostAlias, final String subject) {
		if (hostAlias.contains(subject) || subject.contains(hostAlias)) {
			return true;
		}

		// Parse subjectName out of Common Name
		// If no CN= present, then subjectName is a SAN and needs only wildcard removal
		String subjectName = subject;
		if (subjectName.contains("CN=")) {
			final String[] chopped = subjectName.split("CN=", 2);
			if (chopped != null && chopped.length > 1) {
				final String chop = chopped[1];
				subjectName = chop.split(",", 2)[0];
			}
		}

		subjectName = subjectName.replaceFirst("\\*\\.", ".");
		if (subjectName.length() > 0 && (hostAlias.contains(subjectName) || subjectName.contains(hostAlias))) {
			return true;
		}

		return false;
	}

	public boolean verifyModulus(final PrivateKey privateKey, final X509Certificate certificate) {
		BigInteger privModulus = null;
		if (privateKey instanceof BCRSAPrivateCrtKey) {
			privModulus = ((BCRSAPrivateCrtKey) privateKey).getModulus();
		} else if (privateKey instanceof RSAPrivateCrtKeySpec) {
			privModulus = ((RSAPrivateCrtKeySpec) privateKey).getModulus();
		} else {
			return false;
		}
		BigInteger pubModulus = null;
		final PublicKey publicKey = certificate.getPublicKey();
		if ((publicKey instanceof RSAPublicKeySpec)) {
			pubModulus = ((RSAPublicKeySpec) publicKey).getModulus();
		} else {
			final String[] keyparts = publicKey.toString().split(System.getProperty("line.separator"));
			for (final String part : keyparts) {
				final int start = part.indexOf("modulus: ") + 9;
				if (start < 9) {
					continue;
				} else {
					pubModulus = new BigInteger(part.substring(start));
					break;
				}
			}
		}
		if (privModulus.equals(pubModulus)) {
			return true;
		}
		return false;
	}

	public PrivateKeyDecoder getPrivateKeyDecoder() {
		return privateKeyDecoder;
	}

	public void setPrivateKeyDecoder(final PrivateKeyDecoder privateKeyDecoder) {
		this.privateKeyDecoder = privateKeyDecoder;
	}

	public CertificateDecoder getCertificateDecoder() {
		return certificateDecoder;
	}

	public void setCertificateDecoder(final CertificateDecoder certificateDecoder) {
		this.certificateDecoder = certificateDecoder;
	}
}
