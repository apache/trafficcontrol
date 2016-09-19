package com.comcast.cdn.traffic_control.traffic_router.keystore;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.security.spec.KeySpec;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.Base64;

public class Pkcs8 extends Pkcs {
	public final String HEADER = "-----BEGIN PRIVATE KEY-----";
	public final String FOOTER = "-----END PRIVATE KEY-----";

	public Pkcs8(String data) throws IOException, GeneralSecurityException {
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
		return new PKCS8EncodedKeySpec(Base64.getDecoder().decode((data.getBytes())));
	}
}
