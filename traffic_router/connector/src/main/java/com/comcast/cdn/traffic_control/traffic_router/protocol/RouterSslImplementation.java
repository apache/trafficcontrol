package com.comcast.cdn.traffic_control.traffic_router.protocol;

import org.apache.tomcat.util.net.SSLImplementation;
import org.apache.tomcat.util.net.SSLSupport;
import org.apache.tomcat.util.net.ServerSocketFactory;

import javax.net.ssl.SSLSession;
import java.net.Socket;

public class RouterSslImplementation extends SSLImplementation {
	RouterJsseFactory factory = new RouterJsseFactory();

	@Override
	public String getImplementationName() {
		return getClass().getSimpleName();
	}

	@Override
	public ServerSocketFactory getServerSocketFactory() {
		return factory.getSocketFactory();
	}

	@Override
	public SSLSupport getSSLSupport(final Socket sock) {
		return factory.getSSLSupport(sock);
	}

	@Override
	public SSLSupport getSSLSupport(final SSLSession session) {
		return factory.getSSLSupport(session);
	}
}
