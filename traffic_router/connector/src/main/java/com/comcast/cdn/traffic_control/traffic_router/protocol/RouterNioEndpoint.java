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

import com.comcast.cdn.traffic_control.traffic_router.secure.CertificateRegistry;
import com.comcast.cdn.traffic_control.traffic_router.secure.HandshakeData;
import com.comcast.cdn.traffic_control.traffic_router.secure.KeyManager;
import org.apache.log4j.Logger;
import org.apache.tomcat.util.net.NioEndpoint;
import org.apache.tomcat.util.net.SSLHostConfig;
import org.apache.tomcat.util.net.SSLHostConfigCertificate;
import java.util.Map;
import java.util.Set;

public class RouterNioEndpoint extends NioEndpoint {
    private static final Logger LOGGER = Logger.getLogger(RouterNioEndpoint.class);
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
            final CertificateRegistry certificateRegistry =  keyManager.getCertificateRegistry();
            replaceSSLHosts(certificateRegistry.getHandshakeData());

            //Now let initialiseSsl do it's thing.
            super.initialiseSsl();
            certificateRegistry.setEndPoint(this);
        }
    }

    synchronized private void replaceSSLHosts(final Map<String, HandshakeData> sslHostsData) {
        final Set<String> aliases = sslHostsData.keySet();
        boolean firstAlias = true;
        String lastHostName = "";

        for (final String alias : aliases) {
            final SSLHostConfig sslHostConfig = new SSLHostConfig();
            final SSLHostConfigCertificate cert = new SSLHostConfigCertificate(sslHostConfig, SSLHostConfigCertificate.Type.RSA);
            cert.setCertificateKeyAlias(alias);
            sslHostConfig.addCertificate(cert);
            sslHostConfig.setCertificateKeyAlias(alias);
            sslHostConfig.setHostName(sslHostsData.get(alias).getHostname());
            sslHostConfig.setProtocols("all");
            sslHostConfig.setConfigType(getSslConfigType());
            sslHostConfig.setCertificateVerification("none");
            LOGGER.info("sslHostConfig: "+sslHostConfig.getHostName()+" "+sslHostConfig.getTruststoreAlgorithm());

            if (!sslHostConfig.getHostName().equals(lastHostName)) {
                addSslHostConfig(sslHostConfig, true);
                lastHostName = sslHostConfig.getHostName();
            }

            if (firstAlias && ! "".equals(alias)) {
                // One of the configs must be set as the default
                setDefaultSSLHostConfigName(sslHostConfig.getHostName());
                firstAlias = false;
            }
        }

    }

    synchronized public void reloadSSLHosts(final Map<String, HandshakeData> cr) {
        replaceSSLHosts(cr);

        for (final HandshakeData data : cr.values()) {
            final SSLHostConfig sslHostConfig = sslHostConfigs.get(data.getHostname());
            sslHostConfig.setConfigType(getSslConfigType());
            createSSLContext(sslHostConfig);
        }
    }
}
