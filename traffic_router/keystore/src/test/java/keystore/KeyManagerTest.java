package keystore;

import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyManager;
import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyStoreHelper;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mockito;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import javax.net.ssl.ExtendedSSLSession;
import javax.net.ssl.SNIServerName;
import javax.net.ssl.SSLSocket;
import javax.net.ssl.X509KeyManager;

import java.net.Socket;
import java.security.KeyStore;
import java.security.Principal;
import java.security.cert.X509Certificate;
import java.util.ArrayList;
import java.util.List;
import java.util.Vector;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest({KeyManager.class, KeyStoreHelper.class})
public class KeyManagerTest {
	private KeyManager keyManager;
	private X509KeyManager delegate;
	private X509Certificate x509Certificate;

	@Before
	public void before() throws Exception {
		x509Certificate = mock(X509Certificate.class);

		KeyStore keyStore = PowerMockito.mock(KeyStore.class);
		Mockito.when(keyStore.aliases()).thenAnswer(invocation -> {
			Vector<String> vector = new Vector<>();
			vector.add("deliveryservice3.cdn2.example.com");
			vector.add("deliveryservice2.cdn2.example.com");
			return vector.elements();
		});

		when(keyStore.getCertificateChain("deliveryService2.cdn2.example.com")).thenReturn(new X509Certificate[] {x509Certificate});

		KeyStoreHelper keyStoreHelper = Mockito.mock(KeyStoreHelper.class);
		when(keyStoreHelper.getKeyStore()).thenReturn(keyStore);

		mockStatic(KeyStoreHelper.class);
		PowerMockito.when(KeyStoreHelper.getInstance()).thenReturn(keyStoreHelper);

		whenNew(KeyStoreHelper.class).withNoArguments().thenReturn(keyStoreHelper);
		delegate = mock(X509KeyManager.class);
		keyManager = new KeyManager(delegate);
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
	public void itGetsCertFromKeyStore() {
		assertThat(keyManager.getCertificateChain("deliveryService2.cdn2.example.com")[0], equalTo(x509Certificate));
	}

	@Test
	public void itUsesDelegate() {
		Principal[] principals = {Mockito.mock(Principal.class)};

		keyManager.getClientAliases("foo", principals);
		verify(delegate).getClientAliases("foo", principals);

		String[] strings = new String[] { "foo" };
		Socket socket = Mockito.mock(SSLSocket.class);

		keyManager.chooseClientAlias(strings, principals, socket);
		verify(delegate).chooseClientAlias(strings, principals, socket);

		keyManager.getServerAliases("foo", principals);
		verify(delegate).getServerAliases("foo", principals);
	}

	class TestSNIServerName extends SNIServerName {
		public TestSNIServerName(int type, String name) {
			super(type, name.getBytes());
		}
	}
}
