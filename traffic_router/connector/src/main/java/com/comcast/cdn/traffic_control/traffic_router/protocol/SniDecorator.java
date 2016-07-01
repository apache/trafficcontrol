package com.comcast.cdn.traffic_control.traffic_router.protocol;

import javax.net.ssl.SNIMatcher;
import javax.net.ssl.SSLParameters;
import javax.net.ssl.SSLServerSocket;
import java.net.ServerSocket;
import java.util.ArrayList;
import java.util.Collection;

import static javax.net.ssl.SNIHostName.createSNIMatcher;

public class SniDecorator {
	public ServerSocket addSni(final ServerSocket socket) {
		if (!(socket instanceof SSLServerSocket)) {
			return socket;
		}

		final Collection<SNIMatcher> matchers = new ArrayList<>(1);
		matchers.add(createSNIMatcher("www\\.example\\.(com|org)"));

		final SSLParameters sslParameters = ((SSLServerSocket) socket).getSSLParameters();
		sslParameters.setSNIMatchers(matchers);
		return socket;
	}
}
