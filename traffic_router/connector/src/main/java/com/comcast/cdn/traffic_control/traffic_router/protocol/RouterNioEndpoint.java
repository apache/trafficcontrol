package com.comcast.cdn.traffic_control.traffic_router.protocol;

import com.comcast.cdn.traffic_control.traffic_router.secure.KeyManager;
import org.apache.tomcat.util.net.NioEndpoint;

import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSessionContext;
import javax.net.ssl.TrustManager;
import java.net.InetSocketAddress;
import java.nio.channels.ServerSocketChannel;
import java.util.concurrent.CountDownLatch;

public class RouterNioEndpoint extends NioEndpoint {

	// The same code as its parent class except for SSL
	@SuppressWarnings({"PMD.SignatureDeclareThrowsException", "PMD.NPathComplexity"})
	public void init() throws Exception {

		if (initialized) {
			return;
		}

		serverSock = ServerSocketChannel.open();
		serverSock.socket().setPerformancePreferences(socketProperties.getPerformanceConnectionTime(),
			socketProperties.getPerformanceLatency(),
			socketProperties.getPerformanceBandwidth());

		final InetSocketAddress addr = (address != null ? new InetSocketAddress(address, port) : new InetSocketAddress(port));

		serverSock.socket().bind(addr, backlog);
		serverSock.configureBlocking(true); //mimic APR behavior
		serverSock.socket().setSoTimeout(getSocketProperties().getSoTimeout());

		// Initialize thread count defaults for acceptor, poller
		if (acceptorThreadCount == 0) {
			// FIXME: Doesn't seem to work that well with multiple accept threads
			acceptorThreadCount = 1;
		}

		if (pollerThreadCount <= 0) {
			//minimum one poller thread
			pollerThreadCount = 1;
		}

		stopLatch = new CountDownLatch(pollerThreadCount);

		// Initialize SSL if needed
		if (isSSLEnabled()) {
			sslContext = SSLContext.getInstance(getSslProtocol());
			sslContext.init(wrap(new javax.net.ssl.KeyManager[]{new KeyManager()}), new TrustManager[]{}, null);

			final SSLSessionContext sessionContext = sslContext.getServerSessionContext();

			if (sessionContext != null) {
				sessionContext.setSessionCacheSize(sessionCacheSize);
				sessionContext.setSessionTimeout(sessionTimeout);
			}
		}

		if (oomParachute > 0) {
			reclaimParachute(true);
		}

		selectorPool.open();
		initialized = true;
	}
}
