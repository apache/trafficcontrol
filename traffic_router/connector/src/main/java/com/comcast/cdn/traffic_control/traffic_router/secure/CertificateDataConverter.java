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

import java.security.PrivateKey;
import java.security.cert.X509Certificate;
import java.util.List;
import java.util.stream.Collectors;

public class CertificateDataConverter {
	private final static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(CertificateDataConverter.class);

	private PrivateKeyDecoder privateKeyDecoder = new PrivateKeyDecoder();
	private CertificateDecoder certificateDecoder = new CertificateDecoder();

	public HandshakeData toHandshakeData(final CertificateData certificateData) {
		try {
			final PrivateKey privateKey = privateKeyDecoder.decode(certificateData.getCertificate().getKey());
			final List<String> encodedCertificates = certificateDecoder.doubleDecode(certificateData.getCertificate().getCrt());

			final List<X509Certificate> x509Chain = encodedCertificates.stream()
				.map(encodedCertificate -> certificateDecoder.toCertificate(encodedCertificate))
				.collect(Collectors.toList());

			return new HandshakeData(certificateData.getDeliveryservice(), certificateData.getHostname(),
				x509Chain.toArray(new X509Certificate[x509Chain.size()]), privateKey);

		} catch (Exception e) {
			log.error("Failed to convert certificate data from traffic ops to handshake data! " + e.getClass().getSimpleName() + ": " + e.getMessage(), e);
		}
		return null;
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
