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

import sun.security.util.DerInputStream;
import sun.security.util.DerValue;

import java.io.IOException;
import java.math.BigInteger;
import java.security.GeneralSecurityException;
import java.security.spec.KeySpec;
import java.security.spec.RSAMultiPrimePrivateCrtKeySpec;
import java.util.Base64;

public class Pkcs1 extends Pkcs {
	static public final String HEADER = "-----BEGIN RSA PRIVATE KEY-----";
	static public final String FOOTER = "-----END RSA PRIVATE KEY-----";
	static final int SEQUENCE_LENGTH = 9;

	public Pkcs1(final String data) throws IOException, GeneralSecurityException {
		super(data);
	}

	@Override
	public String getHeader() {
		return HEADER;
	}

	@Override
	public String getFooter() {
		return FOOTER;
	}

	@Override
	protected KeySpec decodeKeySpec(final String data) throws IOException, GeneralSecurityException {
		final String pemData = data.replaceAll(HEADER, "").replaceAll(FOOTER, "").replaceAll("\\s", "");

		final DerInputStream derInputStream = new DerInputStream(Base64.getDecoder().decode(pemData));
		final DerValue[] derSequence = derInputStream.getSequence(0);

		// man 3 rsa
		// -- or --
		// http://linux.die.net/man/3/rsa

		if (derSequence.length < SEQUENCE_LENGTH) {
			throw new GeneralSecurityException("Invalid PKCS1 private key! Missing Private Key Data");
		}

		// We don't need the version data at derSequence[0]
		final BigInteger n = derSequence[1].getBigInteger();
		final BigInteger e = derSequence[2].getBigInteger();
		final BigInteger d = derSequence[3].getBigInteger();
		final BigInteger p = derSequence[4].getBigInteger();
		final BigInteger q = derSequence[5].getBigInteger();
		final BigInteger dmp1 = derSequence[6].getBigInteger();
		final BigInteger dmq1 = derSequence[7].getBigInteger();
		final BigInteger iqmp = derSequence[8].getBigInteger();

		return new RSAMultiPrimePrivateCrtKeySpec(n, e, d, p, q, dmp1, dmq1, iqmp, null);
	}
}
