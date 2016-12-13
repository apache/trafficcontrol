/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
