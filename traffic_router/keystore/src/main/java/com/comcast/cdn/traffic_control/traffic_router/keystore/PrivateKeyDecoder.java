package com.comcast.cdn.traffic_control.traffic_router.keystore;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.security.PrivateKey;

public class PrivateKeyDecoder {
	public PrivateKey decode(final String data) throws IOException, GeneralSecurityException {
		return data.contains(Pkcs1.HEADER) ? new Pkcs1(data).getPrivateKey() : new Pkcs8(data).getPrivateKey();
	}
}
