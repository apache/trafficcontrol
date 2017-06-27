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
import com.comcast.cdn.traffic_control.traffic_router.secure.KeyManager;
import org.apache.tomcat.util.net.NioEndpoint;
import org.apache.tomcat.util.net.SSLHostConfig;
import java.util.List;



public class RouterNioEndpoint extends NioEndpoint {

    // Grabs the aliases from our custom certificate registry, creates a sslHostConfig for them
    // and adds the newly created config to the list of sslHostConfigs.  We also remove the default config
    // since it won't be found in our registry.  This allows OpenSSL to start successfully and serve our
    // certificates.  When we are done we call the parent classes initialiseSsl.
    @SuppressWarnings({"PMD.SignatureDeclareThrowsException"})
    @Override
    protected void initialiseSsl() throws Exception {
        if (isSSLEnabled()) {
            final KeyManager keyManager = new KeyManager();
            final CertificateRegistry certificateRegistry =  keyManager.getCertificateRegistry();
            final List<String> aliases = certificateRegistry.getAliases();

            //remove default config since it won't be found in our keystore
            sslHostConfigs.clear();

            for (final String alias : aliases) {
                final SSLHostConfig sslHostConfig = new SSLHostConfig();
                sslHostConfig.setCertificateKeyAlias(alias);
                addSslHostConfig(sslHostConfig);
            }

            //Now let initialiseSsl do it's thing.
            super.initialiseSsl();

        }
    }
}
