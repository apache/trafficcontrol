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

import org.apache.tomcat.util.net.SSLHostConfigCertificate;
import org.apache.tomcat.util.net.SSLSupport;
import org.apache.tomcat.util.net.jsse.JSSESupport;
import org.apache.tomcat.util.net.SSLImplementation;
import org.apache.tomcat.util.net.SSLUtil;

import javax.net.ssl.SSLSession;

public class RouterSslImplementation extends SSLImplementation {
    @Override
	public SSLSupport getSSLSupport(final SSLSession session) {
        return new JSSESupport(session);
	}

    @Override
        public SSLUtil getSSLUtil(final SSLHostConfigCertificate certificate) {
            return new RouterSslUtil(certificate);
    }

    @Override
    public boolean isAlpnSupported() {
        return true;
    }
}
