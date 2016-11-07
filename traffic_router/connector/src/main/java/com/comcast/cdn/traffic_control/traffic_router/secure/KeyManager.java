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

package com.comcast.cdn.traffic_control.traffic_router.secure;

import javax.net.ssl.ExtendedSSLSession;
import javax.net.ssl.SNIServerName;
import javax.net.ssl.SSLSocket;
import javax.net.ssl.X509KeyManager;
import java.net.Socket;
import java.security.Principal;
import java.security.PrivateKey;
import java.security.cert.X509Certificate;
import java.util.List;
import java.util.Optional;

// Uses the in memory CertificateRegistry to provide dynamic key and certificate management for the router
// The provided default implementation does not allow for the key store to change state
// once the JVM loads the default classes.
public class KeyManager implements X509KeyManager {
	private final static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(KeyManager.class);
	private final CertificateRegistry certificateRegistry = CertificateRegistry.getInstance();

	// To date this method is not getting exercised while running the router
	@Override
	public String chooseClientAlias(final String[] strings, final Principal[] principals, final Socket socket) {
		throw new UnsupportedOperationException("Traffic Router KeyManager does not support choosing Client Alias");
	}

	// To date this method is not getting exercised while running the router
	@Override
	public String[] getServerAliases(final String s, final Principal[] principals) {
		return certificateRegistry.getAliases().toArray(new String[certificateRegistry.getAliases().size()]);
	}

	// To date this method is not getting exercised while running the router
	@Override
	public String[] getClientAliases(final String s, final Principal[] principals) {
		throw new UnsupportedOperationException("Traffic Router KeyManager does not support getting a list of Client Aliases");
	}

	@Override
	public String chooseServerAlias(final String keyType, final Principal[] principals, final Socket socket) {
		if (keyType == null) {
			return null;
		}

		final SSLSocket sslSocket = (SSLSocket) socket;
		final ExtendedSSLSession sslSession = (ExtendedSSLSession) sslSocket.getHandshakeSession();
		final List<SNIServerName> requestedNames = sslSession.getRequestedServerNames();

		final StringBuilder stringBuilder = new StringBuilder();
		for (final SNIServerName requestedName : requestedNames) {
			if (stringBuilder.length() > 0) {
				stringBuilder.append(", ");
			}

			final String sniString = new String(requestedName.getEncoded());
			stringBuilder.append(sniString);

			final Optional<String> optionalAlias = certificateRegistry.getAliases().stream().filter(sniString::contains).findFirst();
			if (optionalAlias.isPresent()) {
				return optionalAlias.get();
			}
		}

		if (stringBuilder.length() > 0) {
			log.warn("No certificate registry aliases matching " + stringBuilder.toString());
		} else {
			log.warn("Client " + sslSocket.getRemoteSocketAddress() + " did not send any Server Name Indicators");
		}
		return null;
	}

	@Override
	public X509Certificate[] getCertificateChain(final String alias) {
		if (certificateRegistry.getAliases().contains(alias)) {
			return certificateRegistry.getHandshakeData(alias).getCertificateChain();
		}

		log.error("No certificate chain for alias " + alias);
		return null;
	}

	@Override
	public PrivateKey getPrivateKey(final String alias) {
		if (certificateRegistry.getAliases().contains(alias)) {
			return certificateRegistry.getHandshakeData(alias).getPrivateKey();
		}

		log.error("No private key for alias " + alias);
		return null;
	}
}
