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

package secure;

import org.apache.traffic_control.traffic_router.secure.CertificateRegistry;
import org.apache.traffic_control.traffic_router.secure.HandshakeData;
import org.apache.traffic_control.traffic_router.secure.KeyManager;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mockito;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import javax.net.ssl.ExtendedSSLSession;
import javax.net.ssl.SNIServerName;
import javax.net.ssl.SSLSocket;

import java.security.PrivateKey;
import java.security.cert.X509Certificate;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.mockStatic;

@RunWith(PowerMockRunner.class)
@PrepareForTest({CertificateRegistry.class})
@PowerMockIgnore("javax.management.*")
public class KeyManagerTest {
	private KeyManager keyManager;
	private X509Certificate x509Certificate;
	private PrivateKey privateKey;

	@Before
	public void before() throws Exception {
		privateKey = mock(PrivateKey.class);
		x509Certificate = mock(X509Certificate.class);

		X509Certificate[] x509Certificates = new X509Certificate[] {
			x509Certificate
		};

		HandshakeData handshakeData = mock(HandshakeData.class);
		when(handshakeData.getCertificateChain()).thenReturn(x509Certificates);
		when(handshakeData.getPrivateKey()).thenReturn(privateKey);

		CertificateRegistry certificateRegistry = mock(CertificateRegistry.class);

		when(certificateRegistry.getAliases()).thenReturn(Arrays.asList(
			"deliveryservice3.cdn2.example.com",
			"deliveryservice2.cdn2.example.com"
		));

		mockStatic(CertificateRegistry.class);
		when(CertificateRegistry.getInstance()).thenReturn(certificateRegistry);
		when(certificateRegistry.getHandshakeData("deliveryservice2.cdn2.example.com")).thenReturn(handshakeData);

		keyManager = new KeyManager();
	}

	@Test
	public void itSelectsServerAlias() {
		List<SNIServerName> sniServerNames = new ArrayList<>();
		sniServerNames.add(new TestSNIServerName(1, "tr.deliveryservice1.cdn1.example.com"));
		sniServerNames.add(new TestSNIServerName(1, "tr.deliveryservice2.cdn2.example.com"));

		ExtendedSSLSession sslExtendedSession = Mockito.mock(ExtendedSSLSession.class);
		when(sslExtendedSession.getRequestedServerNames()).thenReturn(sniServerNames);

		SSLSocket sslSocket = Mockito.mock(SSLSocket.class);
		when(sslSocket.getHandshakeSession()).thenReturn(sslExtendedSession);

		String serverAlias = keyManager.chooseServerAlias("RSA", null, sslSocket);
		assertThat(serverAlias, equalTo("deliveryservice2.cdn2.example.com"));
	}

	@Test
	public void itGetsCertFromRegistry() {
		assertThat(keyManager.getCertificateChain("deliveryservice2.cdn2.example.com")[0], equalTo(x509Certificate));
	}

	@Test
	public void itGetsKeyFromRegistry() {
		assertThat(keyManager.getPrivateKey("deliveryservice2.cdn2.example.com"), equalTo(privateKey));
	}

	class TestSNIServerName extends SNIServerName {
		public TestSNIServerName(int type, String name) {
			super(type, name.getBytes());
		}
	}
}
