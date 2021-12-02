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

package org.apache.traffic_control.traffic_router.secure;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import javax.net.ssl.ExtendedSSLSession;
import javax.net.ssl.SNIServerName;
import javax.net.ssl.SSLEngine;
import javax.net.ssl.SSLSocket;
import javax.net.ssl.X509ExtendedKeyManager;
import javax.net.ssl.X509KeyManager;
import java.net.Socket;
import java.security.Principal;
import java.security.PrivateKey;
import java.security.cert.X509Certificate;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

// Uses the in memory CertificateRegistry to provide dynamic key and certificate management for the router
// The provided default implementation does not allow for the key store to change state
// once the JVM loads the default classes.
public class KeyManager extends X509ExtendedKeyManager implements X509KeyManager {
	private final CertificateRegistry certificateRegistry = CertificateRegistry.getInstance();
	private static final Logger log = LogManager.getLogger(KeyManager.class);
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
	public String chooseEngineServerAlias(final String keyType, final Principal[] issuers, final SSLEngine engine) {
		if (keyType == null) {
			return null;
		}

		final ExtendedSSLSession sslSession = (ExtendedSSLSession) engine.getHandshakeSession();
		return chooseServerAlias(sslSession);
	}

	@Override
	public String chooseServerAlias(final String keyType, final Principal[] principals, final Socket socket) {
		if (keyType == null || socket == null) {
			return null;
		}

		final SSLSocket sslSocket = (SSLSocket) socket;
		final ExtendedSSLSession sslSession = (ExtendedSSLSession) sslSocket.getHandshakeSession();
		return chooseServerAlias(sslSession);
	}

	private String chooseServerAlias(final ExtendedSSLSession sslSession) {
		final List<SNIServerName> requestedNames = sslSession.getRequestedServerNames();

		final StringBuilder stringBuilder = new StringBuilder();
		for (final SNIServerName requestedName : requestedNames) {
			if (stringBuilder.length() > 0) {
				stringBuilder.append(", ");
			}

			final String sniString = new String(requestedName.getEncoded());
			stringBuilder.append(sniString);

			final List<String> partialAliasMatches = certificateRegistry.getAliases().stream().filter(sniString::contains).collect(Collectors.toList());
			Optional<String> alias = partialAliasMatches.stream().filter(sniString::contentEquals).findFirst();
			if (alias.isPresent()) {
			    return alias.get();
			}

			// Not an exact match, some of the aliases may have had the leading zone removed
			final String sniStringTrimmed = sniString.substring(sniString.indexOf('.') + 1);
			alias = partialAliasMatches.stream().filter(sniStringTrimmed::contentEquals).findFirst();
			if (alias.isPresent()) {
			    return alias.get();
			}

		}

		if (stringBuilder.length() > 0) {
			log.warn("KeyManager: No certificate registry aliases matching " + stringBuilder.toString());
		} else {
			log.warn("KeyManager: Client " + sslSession.getPeerHost() + " did not send any Server Name Indicators");
		}
		return null;
	}


	@Override
	public X509Certificate[] getCertificateChain(final String alias) {
		final HandshakeData handshakeData = certificateRegistry.getHandshakeData(alias);
		if (handshakeData != null) {
			return handshakeData.getCertificateChain();
		}

		log.error("KeyManager: No certificate chain for alias " + alias);
		return null;
	}

	@Override
	public PrivateKey getPrivateKey(final String alias) {
		final HandshakeData handshakeData = certificateRegistry.getHandshakeData(alias);
		if (handshakeData != null) {
			return handshakeData.getPrivateKey();
		}

		log.error("KeyManager: No private key for alias " + alias);
		return null;
	}

	public CertificateRegistry getCertificateRegistry() {
		return certificateRegistry;
	}

}
