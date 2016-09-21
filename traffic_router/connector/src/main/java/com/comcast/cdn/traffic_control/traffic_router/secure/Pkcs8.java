package com.comcast.cdn.traffic_control.traffic_router.secure;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.security.spec.KeySpec;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.Base64;

public class Pkcs8 extends Pkcs {
	private final static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(Pkcs8.class);
	public static final String HEADER = "-----BEGIN PRIVATE KEY-----";
	public static final String FOOTER = "-----END PRIVATE KEY-----";

	public Pkcs8(final String data) throws IOException, GeneralSecurityException {
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
		try {
			return new PKCS8EncodedKeySpec(Base64.getDecoder().decode((data.getBytes())));
		} catch (Exception e) {
			log.error("Failed to create PKCS8 Encoded Key Spec " + e.getClass().getCanonicalName() + ": " + e.getMessage(), e);
		}
		return null;
	}
}
