package com.comcast.cdn.traffic_control.traffic_router.secure;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.security.PrivateKey;
import java.util.Base64;

public class PrivateKeyDecoder {
	public PrivateKey decode(final String data) throws IOException, GeneralSecurityException {
		final String decodedData = new String(Base64.getMimeDecoder().decode(data.getBytes()));
		return decodedData.contains(Pkcs1.HEADER) ? new Pkcs1(decodedData).getPrivateKey() : new Pkcs8(decodedData).getPrivateKey();
	}
}
