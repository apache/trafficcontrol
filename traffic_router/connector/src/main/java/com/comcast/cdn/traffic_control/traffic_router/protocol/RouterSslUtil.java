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

import org.apache.juli.logging.Log;
import org.apache.juli.logging.LogFactory;
import org.apache.tomcat.util.net.*;
import org.apache.tomcat.util.net.jsse.JSSEUtil;
import org.apache.tomcat.util.net.openssl.OpenSSLContext;
import org.apache.tomcat.util.net.openssl.OpenSSLEngine;
import org.apache.tomcat.util.net.openssl.OpenSSLUtil;

import javax.net.ssl.KeyManager;
import javax.net.ssl.SSLSessionContext;
import javax.net.ssl.TrustManager;
import java.util.HashSet;
import java.util.List;
import java.util.Locale;
import java.util.Set;

public class RouterSslUtil extends SSLUtilBase {

    private static final Log log = LogFactory.getLog(RouterSslUtil.class);

    public RouterSslUtil(SSLHostConfigCertificate certificate) {
        super(certificate);
    }

    @Override
    protected Log getLog() {
        return log;
    }


    @Override
    protected Set<String> getImplementedProtocols() {
        return OpenSSLEngine.IMPLEMENTED_PROTOCOLS_SET;
    }


    @Override
    protected Set<String> getImplementedCiphers() {
        return OpenSSLEngine.AVAILABLE_CIPHER_SUITES;
    }


    @Override
    public SSLContext createSSLContext(List<String> negotiableProtocols) throws Exception {
        return new OpenSSLContext(certificate, negotiableProtocols);
    }

    @Override
    @SuppressWarnings("PMD.SignatureDeclareThrowsException")
    public javax.net.ssl.KeyManager[] getKeyManagers() throws Exception {
        return new javax.net.ssl.KeyManager[] { new com.comcast.cdn.traffic_control.traffic_router.secure.KeyManager() };
    }

    @Override
    public TrustManager[] getTrustManagers() throws Exception {
            return null;
    }

    @Override
    public void configureSessionContext(SSLSessionContext sslSessionContext) {
    }

}
