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

import org.apache.tomcat.util.net.SSLHostConfigCertificate;
import org.apache.tomcat.util.net.jsse.JSSEUtil;
//import org.apache.tomcat.util.net.openssl.OpenSSLUtil;

//public class RouterSslUtil extends OpenSSLUtil {
public class RouterSslUtil extends JSSEUtil{

    public RouterSslUtil(SSLHostConfigCertificate certificate) {
        super(certificate);
    }

    @Override
    @SuppressWarnings("PMD.SignatureDeclareThrowsException")
    public javax.net.ssl.KeyManager[] getKeyManagers() throws Exception {
        return new javax.net.ssl.KeyManager[] { new com.comcast.cdn.traffic_control.traffic_router.secure.KeyManager() };
    }

}
