package com.comcast.cdn.traffic_control.traffic_router.keystore;

import javax.net.ssl.ExtendedSSLSession;
import javax.net.ssl.SNIServerName;
import javax.net.ssl.SSLSocket;
import javax.net.ssl.X509ExtendedKeyManager;
import javax.net.ssl.X509KeyManager;
import java.net.Socket;
import java.security.KeyStoreException;
import java.security.Principal;
import java.security.PrivateKey;
import java.security.cert.Certificate;
import java.security.cert.X509Certificate;
import java.util.Enumeration;
import java.util.List;

// Does exactly like JSSEKeyManager except for returning certificates
// would have extended it but the class is final
public class KeyManager extends X509ExtendedKeyManager {
	private final static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(KeyManager.class);
	private final X509KeyManager delegate;
	private final KeyStoreHelper keyStoreHelper = KeyStoreHelper.getInstance();

	public KeyManager(final X509KeyManager delegate) {
		this.delegate = delegate;
	}

	@Override
	public String[] getClientAliases(final String s, final Principal[] principals) {
		return delegate.getClientAliases(s, principals);
	}

	@Override
	public String chooseClientAlias(final String[] strings, final Principal[] principals, final Socket socket) {
		return delegate.chooseClientAlias(strings, principals, socket);
	}

	@Override
	public String[] getServerAliases(final String s, final Principal[] principals) {
		return delegate.getServerAliases(s, principals);
	}

	@Override
	public String chooseServerAlias(final String keyType, final Principal[] principals, final Socket socket) {
		if (keyType == null) {
			return null;
		}

		final SSLSocket sslSocket = (SSLSocket) socket;
		final ExtendedSSLSession sslSession = (ExtendedSSLSession) sslSocket.getHandshakeSession();
		final List<SNIServerName> requestedNames = sslSession.getRequestedServerNames();

		for (final SNIServerName requestedName : requestedNames) {
			try {
				final Enumeration<String> aliases = keyStoreHelper.getKeyStore().aliases();

				while (aliases.hasMoreElements()) {
					final String alias = aliases.nextElement();
					final String sniString = new String(requestedName.getEncoded());

					if (sniString.contains(alias)) {
						return alias;
					}
				}
			} catch (KeyStoreException e) {
				log.error("Failed getting aliases from keystore: " + e.getMessage());
			}
		}

		return delegate.chooseServerAlias(keyType, principals, socket);
	}

	@Override
	public X509Certificate[] getCertificateChain(final String s) {
		try {
			final Certificate[] certificates = keyStoreHelper.getKeyStore().getCertificateChain(s);
			final X509Certificate[] x509Certificates = new X509Certificate[certificates.length];
			int i = 0;

			for (final Certificate certificate : certificates) {
				if (certificate instanceof X509Certificate) {
					x509Certificates[i] = (X509Certificate) certificate;
					i++;
				}
			}

			return x509Certificates;
		} catch (Exception e) {
			log.error("Failed retrieving certificate chain from keystore for alias '" + s + "' : " + e.getMessage());
		}

		return null;
	}

	@Override
	public PrivateKey getPrivateKey(final String s) {
		return delegate.getPrivateKey(s);
	}
}
