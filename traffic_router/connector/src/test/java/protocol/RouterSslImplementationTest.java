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

package protocol;

import org.apache.traffic_control.traffic_router.protocol.RouterSslImplementation;
import org.apache.traffic_control.traffic_router.protocol.RouterSslUtil;
import org.apache.tomcat.util.net.SSLHostConfig;
import org.apache.tomcat.util.net.SSLHostConfigCertificate;
import org.apache.tomcat.util.net.SSLSupport;
import org.apache.tomcat.util.net.SSLUtil;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.powermock.api.mockito.PowerMockito;

import javax.net.ssl.SSLSession;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.instanceOf;

@RunWith(PowerMockRunner.class)
@PrepareForTest({RouterSslImplementation.class, SSLHostConfigCertificate.class, RouterSslUtil.class})
@PowerMockIgnore("javax.management.*")
public class RouterSslImplementationTest {
	SSLSession sslSession = PowerMockito.mock(SSLSession.class);
	SSLHostConfig sslHostConfig = PowerMockito.mock(SSLHostConfig.class);
	SSLHostConfigCertificate.Type type = PowerMockito.mock(SSLHostConfigCertificate.Type.class);
	SSLHostConfigCertificate sslHostConfigCertificate = new SSLHostConfigCertificate(sslHostConfig, type);
	RouterSslUtil sslutil = PowerMockito.mock(RouterSslUtil.class);

	@Test
	public void itReturnsSSLSupport() throws Exception {
		assertThat(new RouterSslImplementation().getSSLSupport(sslSession), instanceOf(SSLSupport.class));
	}

	@Test
	public void itReturnsSSLUtil() throws Exception {
		PowerMockito.whenNew(RouterSslUtil.class).withArguments(sslHostConfigCertificate).thenReturn(sslutil);
		assertThat(new RouterSslImplementation().getSSLUtil(sslHostConfigCertificate), instanceOf(SSLUtil.class));
	}

	@Test
	public void itRegistersSSLHostConfigs() throws Exception {

	}

	@Before
	public void before() {
	}
}
