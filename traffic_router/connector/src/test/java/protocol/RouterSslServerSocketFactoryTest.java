package protocol;

import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyStoreHelper;
import com.comcast.cdn.traffic_control.traffic_router.protocol.RouterSslServerSocketFactory;
import org.apache.tomcat.util.net.jsse.JSSESocketFactory;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.invocation.InvocationOnMock;
import org.mockito.stubbing.Answer;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import javax.net.ssl.KeyManagerFactory;
import javax.net.ssl.X509KeyManager;

import java.security.KeyStore;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Matchers.anyObject;
import static org.mockito.Matchers.anyString;
import static org.mockito.Matchers.eq;
import static org.mockito.Mockito.doNothing;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.spy;
import static org.mockito.Mockito.verify;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.when;

@RunWith(PowerMockRunner.class)
@PrepareForTest({
	RouterSslServerSocketFactory.class,
	JSSESocketFactory.class,
	KeyStore.class,
	KeyManagerFactory.class,
	KeyStoreHelper.class,
	System.class})
public class RouterSslServerSocketFactoryTest {

	String keyStoreValue = null;

	@Test
	public void itAddSniData() throws Exception {
		mockStatic(System.class);

		when(System.setProperty(eq("javax.net.ssl.keyStore"), anyString())).thenAnswer(new Answer<Void>() {
			@Override
			public Void answer(InvocationOnMock invocation) throws Throwable {
				keyStoreValue = (String) invocation.getArguments()[1];
				return null;
			}
		});

		when(System.getProperty(anyString(), anyString())).thenCallRealMethod();

		KeyStoreHelper keyStoreHelper = PowerMockito.mock(KeyStoreHelper.class);
		PowerMockito.when(keyStoreHelper.getKeyPass()).thenReturn("super-secret".toCharArray());

		mockStatic(KeyStoreHelper.class);
		PowerMockito.when(KeyStoreHelper.getInstance()).thenReturn(keyStoreHelper);

		KeyManagerFactory keyManagerFactory = PowerMockito.mock(KeyManagerFactory.class);

		when(keyManagerFactory.getKeyManagers()).thenReturn(new X509KeyManager[] {
			mock(X509KeyManager.class), mock(X509KeyManager.class)
		});

		mockStatic(KeyManagerFactory.class);
		PowerMockito.when(KeyManagerFactory.getInstance("SunX509")).thenReturn(keyManagerFactory);

		KeyStore keyStore = PowerMockito.mock(KeyStore.class);
		when(keyStore.isKeyEntry("some-alias")).thenReturn(true);
		mockStatic(KeyStore.class);
		when(KeyStore.getInstance("JKS")).thenReturn(keyStore);
		doNothing().when(keyStore).load(anyObject(), anyObject());

		when(keyStoreHelper.getKeyStore()).thenReturn(keyStore);
		when(keyStoreHelper.getKeystorePath()).thenReturn("/opt/traffic_router/db/.keystore");

		RouterSslServerSocketFactory socketFactory = spy(new RouterSslServerSocketFactory());
		assertThat(socketFactory.getKeyManagers("JKS", null, "SunX509", "some-alias").length, equalTo(2));
		verify(keyStoreHelper).getKeyPass();
		assertThat(keyStoreValue, equalTo("/opt/traffic_router/db/.keystore"));
		verify(socketFactory).setAttribute("keystore", "/opt/traffic_router/db/.keystore");
		verify(keyStoreHelper).getKeyStore();
	}
}
