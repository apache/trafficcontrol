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

// Uses the KeyStoreHelper to provide dynamic key and certificate management for the router
// The provided default implementation does not allow for the key store to change state
// once the JVM loads the default classes.
public class KeyManager extends X509ExtendedKeyManager {
	private final static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(KeyManager.class);
	private final X509KeyManager delegate;
	private final KeyStoreHelper keyStoreHelper = KeyStoreHelper.getInstance();

	public KeyManager(final X509KeyManager delegate) {
		this.delegate = delegate;
	}

	// To date this method is not getting exercised while running the router
	@Override
	public String chooseClientAlias(final String[] strings, final Principal[] principals, final Socket socket) {
		return delegate.chooseClientAlias(strings, principals, socket);
	}

	// To date this method is not getting exercised while running the router
	@Override
	public String[] getServerAliases(final String s, final Principal[] principals) {
		return delegate.getServerAliases(s, principals);
	}

	// To date this method is not getting exercised while running the router
	@Override
	public String[] getClientAliases(final String s, final Principal[] principals) {
		return delegate.getClientAliases(s, principals);
	}

	@Override
	public String chooseServerAlias(final String keyType, final Principal[] principals, final Socket socket) {
		if (keyType == null) {
			return null;
		}

		if (keyStoreHelper.getLastModified() > keyStoreHelper.getLastLoaded()) {
			log.warn("Reloading keystore from filesystem");
			keyStoreHelper.reload();
		}

		final SSLSocket sslSocket = (SSLSocket) socket;
		final ExtendedSSLSession sslSession = (ExtendedSSLSession) sslSocket.getHandshakeSession();
		final List<SNIServerName> requestedNames = sslSession.getRequestedServerNames();

		StringBuilder stringBuilder = new StringBuilder();
		for (final SNIServerName requestedName : requestedNames) {
			if (stringBuilder.length() > 0) {
				stringBuilder.append(", ");
			}

			final String sniString = new String(requestedName.getEncoded());
			stringBuilder.append(sniString);

			try {
				final Enumeration<String> aliases = keyStoreHelper.getKeyStore().aliases();

				if (!aliases.hasMoreElements()) {
					log.error("Keystore has NO aliases!");
					return null;
				}

				while (aliases.hasMoreElements()) {
					final String alias = aliases.nextElement();
					if (sniString.contains(alias)) {
						return alias;
					}
				}
			} catch (KeyStoreException e) {
				log.error("Failed getting aliases from keystore: " + e.getMessage());
			}
		}

		if (stringBuilder.length() > 0) {
			log.warn("No keystore aliases matching " + stringBuilder.toString());
		} else {
			log.warn("Client " + sslSocket.getRemoteSocketAddress() + " did not send any Server Name Indicators");
		}
		return null;
	}

	private X509Certificate[] reverse(X509Certificate[] x509Certificates) {
		int low = 0;
		int high = x509Certificates.length - 1;

		while (low < high) {
			final X509Certificate tmp = x509Certificates[low];
			x509Certificates[low] = x509Certificates[high];
			x509Certificates[high] = tmp;
			low++;
			high--;
		}

		return x509Certificates;
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

			return reverse(x509Certificates);
		} catch (Exception e) {
			log.error("Failed retrieving certificate chain from keystore for alias '" + s + "' (" + e.getClass().getCanonicalName() + "): " + e.getMessage());
		}

		return null;
	}

	@Override
	public PrivateKey getPrivateKey(String alias) {
		return keyStoreHelper.getPrivateKey(alias);
	}
}
