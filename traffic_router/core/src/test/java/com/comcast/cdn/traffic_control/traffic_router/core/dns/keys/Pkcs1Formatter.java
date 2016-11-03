package com.comcast.cdn.traffic_control.traffic_router.core.dns.keys;

import sun.security.rsa.RSAPrivateCrtKeyImpl;
import sun.security.util.DerOutputStream;
import sun.security.util.DerValue;

import java.io.IOException;
import java.security.interfaces.RSAPublicKey;

public class Pkcs1Formatter {

	// https://tools.ietf.org/html/rfc3447#appendix-A.1.1

	public byte[] toBytes(RSAPrivateCrtKeyImpl key) throws IOException {
		byte tag = 2;
		DerValue[] outputSequence = new DerValue[] {
			new DerValue(tag, new byte[]{0}),
			new DerValue(tag, key.getModulus().toByteArray()),
			new DerValue(tag, key.getPublicExponent().toByteArray()),
			new DerValue(tag, key.getPrivateExponent().toByteArray()),
			new DerValue(tag, key.getPrimeP().toByteArray()),
			new DerValue(tag, key.getPrimeQ().toByteArray()),
			new DerValue(tag, key.getPrimeExponentP().toByteArray()),
			new DerValue(tag, key.getPrimeExponentQ().toByteArray()),
			new DerValue(tag, key.getCrtCoefficient().toByteArray()),
		};

		DerOutputStream outputStream = new DerOutputStream();

		outputStream.putSequence(outputSequence);
		outputStream.flush();

		return outputStream.toByteArray();
	}

	public byte[] toBytes(RSAPublicKey key) throws IOException {
		byte tag = 2;
		DerValue[] outputSequence = new DerValue[] {
			new DerValue(tag, key.getModulus().toByteArray()),
			new DerValue(tag, key.getPublicExponent().toByteArray())
		};

		DerOutputStream outputStream = new DerOutputStream();

		outputStream.putSequence(outputSequence);
		outputStream.flush();

		return outputStream.toByteArray();
	}
}
