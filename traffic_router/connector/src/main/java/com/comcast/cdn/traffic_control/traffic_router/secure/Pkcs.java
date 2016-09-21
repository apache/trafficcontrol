package com.comcast.cdn.traffic_control.traffic_router.secure;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.security.KeyFactory;
import java.security.PrivateKey;
import java.security.spec.KeySpec;

@SuppressWarnings("PMD.AbstractNaming")
public abstract class Pkcs {
	private final String data;
	private final PrivateKey privateKey;
	private KeySpec keySpec;

	public Pkcs(final String data) throws IOException, GeneralSecurityException {
		this.data = data;
		keySpec = toKeySpec(data);
		privateKey = KeyFactory.getInstance("RSA").generatePrivate(keySpec);
	}

	public String getData() {
		return data;
	}

	public KeySpec getKeySpec() {
		return keySpec;
	}

	public void setKeySpec(final KeySpec keySpec) {
		this.keySpec = keySpec;
	}

	public PrivateKey getPrivateKey() {
		return privateKey;
	}

	public abstract String getHeader();

	public abstract String getFooter();

	private String stripHeaderAndFooter(final String data) {
		return data.replaceAll(getHeader(), "").replaceAll(getFooter(), "").replaceAll("\\s", "");
	}

	protected abstract KeySpec decodeKeySpec(final String data) throws IOException, GeneralSecurityException;

	private KeySpec toKeySpec(final String data) throws IOException, GeneralSecurityException {
		return decodeKeySpec(stripHeaderAndFooter(data));
	}
}
