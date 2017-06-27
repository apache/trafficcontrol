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

import org.apache.log4j.Logger;
import sun.security.util.DerInputStream;
import sun.security.util.DerValue;

import java.io.IOException;
import java.math.BigInteger;
import java.security.GeneralSecurityException;
import java.security.spec.KeySpec;
import java.security.spec.PKCS8EncodedKeySpec;
//import java.security.spec.RSAMultiPrimePrivateCrtKeySpec;
import java.security.spec.RSAPublicKeySpec;
import java.util.Base64;

import static com.comcast.cdn.traffic_control.traffic_router.secure.Pkcs1.FOOTER;
import static com.comcast.cdn.traffic_control.traffic_router.secure.Pkcs1.HEADER;

public class Pkcs1KeySpecDecoder {
	// https://tools.ietf.org/html/rfc3447#appendix-A.1.1


	static final int PRIVATE_SEQUENCE_LENGTH = 9;
	static final int PUBLIC_SEQUENCE_LENGTH = 2;
	private static final Logger LOGGER = Logger.getLogger(Pkcs1KeySpecDecoder.class);

	public KeySpec decode(final String data) throws IOException, GeneralSecurityException {
		final String pemData = data.replaceAll(HEADER, "").replaceAll(FOOTER, "").replaceAll("\\s", "");

		final DerInputStream derInputStream = new DerInputStream(Base64.getDecoder().decode(pemData));
		final DerValue[] derSequence = derInputStream.getSequence(0);

		if (derSequence.length != PUBLIC_SEQUENCE_LENGTH && derSequence.length != PRIVATE_SEQUENCE_LENGTH) {
			throw new GeneralSecurityException("Invalid PKCS1 key! Missing Key Data, incorrect number of DER values for either public or private key");
		}

		if (derSequence.length == PUBLIC_SEQUENCE_LENGTH) {
			final BigInteger n = derSequence[0].getBigInteger();
			final BigInteger e = derSequence[1].getBigInteger();
			return new RSAPublicKeySpec(n, e);
		}

		// man 3 rsa
		// -- or --
		// http://linux.die.net/man/3/rsa

		// We don't need the version data at derSequence[0]
//		final BigInteger n = derSequence[1].getBigInteger();
//		final BigInteger e = derSequence[2].getBigInteger();
//		final BigInteger d = derSequence[3].getBigInteger();
//		final BigInteger p = derSequence[4].getBigInteger();
//		final BigInteger q = derSequence[5].getBigInteger();
//		final BigInteger dmp1 = derSequence[6].getBigInteger();
//		final BigInteger dmq1 = derSequence[7].getBigInteger();
//		final BigInteger iqmp = derSequence[8].getBigInteger();
//
//		return new RSAMultiPrimePrivateCrtKeySpec(n, e, d, p, q, dmp1, dmq1, iqmp, null);

		//Convert to PKCS8 since OpenSSL doesn't support PKCS1.  This works because of the BouncyCastle security provider.
		try {
			return new PKCS8EncodedKeySpec(Base64.getDecoder().decode((data.getBytes())));
		} catch (Exception e) {
			LOGGER.error("Error converting to PKCS8 Encoded Key Spec " + e.getClass().getCanonicalName() + ": " + e.getMessage(), e);
		}
		return null;
	}
}