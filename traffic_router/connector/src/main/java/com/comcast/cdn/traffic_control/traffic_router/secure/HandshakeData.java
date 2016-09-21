package com.comcast.cdn.traffic_control.traffic_router.secure;

import java.security.PrivateKey;
import java.security.cert.X509Certificate;

public class HandshakeData {
	protected static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(HandshakeData.class);

	private final String deliveryService;
	private final String hostname;
	private final X509Certificate[] certificateChain;
	private PrivateKey privateKey;

	public HandshakeData(final String deliveryService, final String hostname, final X509Certificate[] certificateChain, final PrivateKey privateKey) {
		this.deliveryService = deliveryService;
		this.hostname = hostname;
		this.certificateChain = certificateChain;
		this.privateKey = privateKey;
	}

	public String getDeliveryService() {
		return deliveryService;
	}

	public String getHostname() {
		return hostname;
	}

	public X509Certificate[] getCertificateChain() {
		return certificateChain;
	}

	public PrivateKey getPrivateKey() {
		return privateKey;
	}

	public void setPrivateKey(final PrivateKey privateKey) {
		this.privateKey = privateKey;
	}
}
