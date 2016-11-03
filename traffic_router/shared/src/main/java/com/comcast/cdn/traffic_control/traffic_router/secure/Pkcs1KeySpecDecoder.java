package com.comcast.cdn.traffic_control.traffic_router.secure;

import sun.security.util.DerInputStream;
import sun.security.util.DerValue;

import java.io.IOException;
import java.math.BigInteger;
import java.security.GeneralSecurityException;
import java.security.spec.KeySpec;
import java.security.spec.RSAMultiPrimePrivateCrtKeySpec;
import java.security.spec.RSAPublicKeySpec;
import java.util.Base64;

public class Pkcs1KeySpecDecoder {
	// https://tools.ietf.org/html/rfc3447#appendix-A.1.1

	static public final String HEADER = "-----BEGIN RSA PRIVATE KEY-----";
	static public final String FOOTER = "-----END RSA PRIVATE KEY-----";
	static final int PRIVATE_SEQUENCE_LENGTH = 9;
	static final int PUBLIC_SEQUENCE_LENGTH = 2;

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
			return new RSAPublicKeySpec(n,e);
		}

		// man 3 rsa
		// -- or --
		// http://linux.die.net/man/3/rsa

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
