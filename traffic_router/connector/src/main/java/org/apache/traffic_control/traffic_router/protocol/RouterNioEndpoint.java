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

package org.apache.traffic_control.traffic_router.protocol;

import org.apache.traffic_control.traffic_router.secure.CertificateRegistry;
import org.apache.traffic_control.traffic_router.secure.HandshakeData;
import org.apache.traffic_control.traffic_router.secure.KeyManager;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.tomcat.util.net.NioEndpoint;
import org.apache.tomcat.util.net.SSLHostConfig;
import org.apache.tomcat.util.net.SSLHostConfigCertificate;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Set;

public class RouterNioEndpoint extends NioEndpoint {
	private static final Logger LOGGER = LogManager.getLogger(RouterNioEndpoint.class);
	private String protocols;

	// Grabs the aliases from our custom certificate registry, creates a sslHostConfig for them
	// and adds the newly created config to the list of sslHostConfigs.  We also remove the default config
	// since it won't be found in our registry.  This allows OpenSSL to start successfully and serve our
	// certificates.  When we are done we call the parent classes initialiseSsl.
	@SuppressWarnings({"PMD.SignatureDeclareThrowsException"})
	@Override
	protected void initialiseSsl() throws Exception {
		if (isSSLEnabled()) {
			destroySsl();
			sslHostConfigs.clear();
			final KeyManager keyManager = new KeyManager();
			final CertificateRegistry certificateRegistry = keyManager.getCertificateRegistry();
			replaceSSLHosts(certificateRegistry.getHandshakeData());

			//Now let initialiseSsl do it's thing.
			super.initialiseSsl();
			certificateRegistry.setEndPoint(this);
		}
	}

	synchronized private List<String> replaceSSLHosts(final Map<String, HandshakeData> sslHostsData) {
		final Set<String> aliases = sslHostsData.keySet();
		String lastHostName = "";
		final List<String> failedUpdates = new ArrayList<>();

		for (final String alias : aliases) {
			final SSLHostConfig sslHostConfig = new SSLHostConfig();
			final SSLHostConfigCertificate cert = new SSLHostConfigCertificate(sslHostConfig, SSLHostConfigCertificate.Type.RSA);
			sslHostConfig.setHostName(sslHostsData.get(alias).getHostname());
			cert.setCertificateKeyAlias(alias);
			sslHostConfig.addCertificate(cert);
			sslHostConfig.setProtocols(protocols != null ? protocols : "all");
			sslHostConfig.setSslProtocol(sslHostConfig.getSslProtocol());
			sslHostConfig.setCertificateVerification("none");
			LOGGER.info("sslHostConfig: "+sslHostConfig.getHostName() + " " + sslHostConfig.getTruststoreAlgorithm());

		if (!sslHostConfig.getHostName().equals(lastHostName)){
			try{
				addSslHostConfig(sslHostConfig, true);
			} catch (Exception fubar){
				LOGGER.error("In RouterNioEndpoint.replaceSSLHosts, sslHostConfig and certs did not get replaced " +
				  "for host: " + sslHostConfig.getHostName() + ", because of execption - " + fubar.toString());
				failedUpdates.add(alias);
			}
			lastHostName = sslHostConfig.getHostName();
		}

			if (CertificateRegistry.DEFAULT_SSL_KEY.equals(alias) && !failedUpdates.contains(alias)){
				// One of the configs must be set as the default
				setDefaultSSLHostConfigName(sslHostConfig.getHostName());
			}
		}
		return failedUpdates;
	}

	synchronized public List<String> reloadSSLHosts(final Map<String, HandshakeData> cr) {
		final List<String> failedUpdates = replaceSSLHosts(cr);
		if (!failedUpdates.isEmpty()) {
			failedUpdates.forEach(alias-> {
				cr.remove(alias);
			});
		}

		final List<String> failedContextUpdates = new ArrayList<>();
		for (final String alias : cr.keySet()) {
			try{
				final HandshakeData data = cr.get(alias);
				final SSLHostConfig sslHostConfig = sslHostConfigs.get(data.getHostname());
				sslHostConfig.setSslProtocol(sslHostConfig.getSslProtocol());
				createSSLContext(sslHostConfig);
			}
			catch (Exception rfubar) {
				LOGGER.error("In RouterNioEndpoint could not create new SSLContext for cert " + alias +
						" because of exception: "+rfubar.toString());
				failedContextUpdates.add(alias);
			}
		}

		if (!failedContextUpdates.isEmpty()) {
			failedUpdates.addAll(failedContextUpdates);
		}

		return failedUpdates;
	}

	@Override
	protected SSLHostConfig getSSLHostConfig(final String sniHostName){
		return super.getSSLHostConfig(sniHostName == null ? null : sniHostName.toLowerCase());
	}

	public String getProtocols() {
		return protocols;
	}

	public void setProtocols(final String protocols) {
		this.protocols = protocols;
	}

}
