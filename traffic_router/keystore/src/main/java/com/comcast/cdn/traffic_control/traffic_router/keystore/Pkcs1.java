package com.comcast.cdn.traffic_control.traffic_router.keystore;

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
	protected KeySpec decodeKeySpec(String data) throws IOException, GeneralSecurityException {
		String pemData = data;
		pemData = pemData.replaceAll(HEADER, "");
		pemData = pemData.replaceAll(FOOTER, "");

		DerInputStream derInputStream = new DerInputStream(Base64.getDecoder().decode(pemData));
		DerValue[] derSequence = derInputStream.getSequence(0);

		// man 3 rsa
		// -- or --
		// http://linux.die.net/man/3/rsa

		if (derSequence.length < SEQUENCE_LENGTH) {
			throw new GeneralSecurityException("Invalid PKCS1 private key! Missing Private Key Data");
		}

		// We don't need the version data at derSequence[0]
		BigInteger n = derSequence[1].getBigInteger();
		BigInteger e = derSequence[2].getBigInteger();
		BigInteger d = derSequence[3].getBigInteger();
		BigInteger p = derSequence[4].getBigInteger();
		BigInteger q = derSequence[5].getBigInteger();
		BigInteger dmp1 = derSequence[6].getBigInteger();
		BigInteger dmq1 = derSequence[7].getBigInteger();
		BigInteger iqmp = derSequence[8].getBigInteger();

		return new RSAMultiPrimePrivateCrtKeySpec(n, e, d, p, q, dmp1, dmq1, iqmp, null);
	}
}
