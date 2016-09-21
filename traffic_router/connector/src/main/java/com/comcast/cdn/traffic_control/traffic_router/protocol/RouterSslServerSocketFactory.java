package com.comcast.cdn.traffic_control.traffic_router.protocol;

import com.comcast.cdn.traffic_control.traffic_router.secure.KeyManager;
import org.apache.tomcat.util.net.jsse.JSSESocketFactory;

// Wrap JSSEKeyManager with our own key manager so we can control handing out certificates
public class RouterSslServerSocketFactory extends JSSESocketFactory {
	protected static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(RouterSslServerSocketFactory.class);

	public RouterSslServerSocketFactory() {
		super();
	}

	@Override
	@SuppressWarnings("PMD.SignatureDeclareThrowsException")
	public javax.net.ssl.KeyManager[] getKeyManagers(final String keystoreType, final String keystoreProvider, final String algorithm, final String keyAlias) throws Exception {
		return new javax.net.ssl.KeyManager[] { new KeyManager() };
	}
}
