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
import org.bouncycastle.asn1.ASN1Integer;
import org.bouncycastle.asn1.ASN1Sequence;
import org.bouncycastle.asn1.ASN1SequenceParser;

import java.io.IOException;
import java.math.BigInteger;
import java.security.GeneralSecurityException;
import java.security.spec.KeySpec;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.RSAPublicKeySpec;
import java.util.Base64;

import static org.apache.traffic_control.traffic_router.secure.Pkcs1.FOOTER;
import static org.apache.traffic_control.traffic_router.secure.Pkcs1.HEADER;

public class Pkcs1KeySpecDecoder {
	// https://tools.ietf.org/html/rfc3447#appendix-A.1.1


	static final int PRIVATE_SEQUENCE_LENGTH = 9;
	static final int PUBLIC_SEQUENCE_LENGTH = 2;
	private static final Logger LOGGER = LogManager.getLogger(Pkcs1KeySpecDecoder.class);

	public KeySpec decode(final String data) throws IOException, GeneralSecurityException {
		final String pemData = data.replaceAll(HEADER, "").replaceAll(FOOTER, "").replaceAll("\\s", "");
		final ASN1Sequence asn1Sequence = ASN1Sequence.getInstance(Base64.getDecoder().decode(pemData));
		final int sequenceLength = asn1Sequence.toArray().length;
		if(sequenceLength != PUBLIC_SEQUENCE_LENGTH && sequenceLength != PRIVATE_SEQUENCE_LENGTH) {
			throw new GeneralSecurityException("Invalid PKCS1 key! Missing Key Data, incorrect number of DER values for either public or private key");
		}
		if (asn1Sequence.toArray().length == PUBLIC_SEQUENCE_LENGTH) {
			final ASN1SequenceParser asn1Parser = asn1Sequence.parser();
			final BigInteger n = ((ASN1Integer) asn1Parser.readObject()).getValue();
			final BigInteger e = ((ASN1Integer) asn1Parser.readObject()).getValue();
			return new RSAPublicKeySpec(n, e);
		}

		// man 3 rsa
		// -- or --
		// http://linux.die.net/man/3/rsa
		//Convert to PKCS8 since OpenSSL doesn't support PKCS1.  This works because of the BouncyCastle security provider.
		try {
			return new PKCS8EncodedKeySpec(Base64.getDecoder().decode((data.getBytes())));
		} catch (Exception e) {
			LOGGER.error("Error converting to PKCS8 Encoded Key Spec " + e.getClass().getCanonicalName() + ": " + e.getMessage(), e);
		}
		return null;
	}
}