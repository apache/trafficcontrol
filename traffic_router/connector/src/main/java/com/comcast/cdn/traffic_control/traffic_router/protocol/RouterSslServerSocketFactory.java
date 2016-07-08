package com.comcast.cdn.traffic_control.traffic_router.protocol;

import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyManager;
import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyStoreHelper;
import org.apache.tomcat.util.net.jsse.JSSESocketFactory;

import javax.net.ssl.X509KeyManager;
import java.io.IOException;
import java.security.KeyStore;

// Wrap JSSEKeyManager with our own key manager so we can control handing out certificates
public class RouterSslServerSocketFactory extends JSSESocketFactory {
	protected static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(RouterSslServerSocketFactory.class);

	@Override
	@SuppressWarnings("PMD.SignatureDeclareThrowsException")
	public javax.net.ssl.KeyManager[] getKeyManagers(final String keystoreType, final String keystoreProvider, final String algorithm, final String keyAlias) throws Exception {
		final javax.net.ssl.KeyManager[] keyManagers = super.getKeyManagers(keystoreType, keystoreProvider, algorithm, keyAlias);

		for (int i = 0; i < keyManagers.length; i++) {
			keyManagers[i] = new KeyManager((X509KeyManager) keyManagers[i]);
		}

		return keyManagers;
	}

	@Override
	protected String getKeystorePassword() {
		return new String(KeyStoreHelper.getInstance().getKeyPass());
	}

	@Override
	protected KeyStore getKeystore(final String type, final String provider, final String pass) throws IOException {
		final String keyStorePath = KeyStoreHelper.getInstance().getKeystorePath();
		setAttribute("keystore", keyStorePath);
		System.setProperty("javax.net.ssl.keyStore", keyStorePath);
		return KeyStoreHelper.getInstance().getKeyStore();
	}
}
