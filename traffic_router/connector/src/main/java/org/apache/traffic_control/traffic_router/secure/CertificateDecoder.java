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

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.io.ByteArrayInputStream;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.util.ArrayList;
import java.util.Base64;
import java.util.List;

public class CertificateDecoder {
	private static final Logger log = LogManager.getLogger(CertificateDecoder.class);

	private static final String CRT_HEADER = "-----BEGIN CERTIFICATE-----";
	private static final String CRT_FOOTER = "-----END CERTIFICATE-----";
	private static final String PEM_FOOTER_PREFIX = "-----END";

	public List<String> doubleDecode(final String encoded) {
		final byte[] decodedBytes = Base64.getMimeDecoder().decode(encoded.getBytes());

		final List<String> encodedPemItems = new ArrayList<>();

		final String[] lines = new String(decodedBytes).split("\\r?\\n");
		final StringBuilder builder = new StringBuilder();

		for (final String line : lines) {
			builder.append(line);

			if (line.startsWith(PEM_FOOTER_PREFIX)) {
				encodedPemItems.add(builder.toString());
				builder.setLength(0);
			}
		}

		if (encodedPemItems.isEmpty()) {
			if (builder.length() == 0) {
				log.warn("Failed base64 decoding");
			} else {
				encodedPemItems.add(builder.toString());
			}
		}

		return encodedPemItems;
	}

	@SuppressWarnings("PMD.AvoidThrowingRawExceptionTypes")
	public X509Certificate toCertificate(final String encodedCertificate) {
		final byte[] encodedBytes = Base64.getDecoder().decode(encodedCertificate.replaceAll(CRT_HEADER, "").replaceAll(CRT_FOOTER, ""));

		try (ByteArrayInputStream stream = new ByteArrayInputStream(encodedBytes)) {
			return (X509Certificate) CertificateFactory.getInstance("X.509").generateCertificate(stream);
		} catch (Exception e) {
			final String message = "Failed to decode certificate data to X509! " + e.getClass().getSimpleName() + ": " + e.getMessage();
			log.error(message, e);
			throw new RuntimeException(message,e);
		}
	}
}
