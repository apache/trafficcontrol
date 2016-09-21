package protocol;

import com.comcast.cdn.traffic_control.traffic_router.protocol.RouterSslServerSocketFactory;
import org.apache.tomcat.util.net.jsse.JSSESocketFactory;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import javax.net.ssl.KeyManagerFactory;
import javax.net.ssl.X509KeyManager;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.spy;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.when;

@RunWith(PowerMockRunner.class)
@PrepareForTest({RouterSslServerSocketFactory.class, JSSESocketFactory.class, KeyManagerFactory.class})
public class RouterSslServerSocketFactoryTest {

	@Test
	public void itAddSniData() throws Exception {
		KeyManagerFactory keyManagerFactory = PowerMockito.mock(KeyManagerFactory.class);

		when(keyManagerFactory.getKeyManagers()).thenReturn(new X509KeyManager[] {
			mock(X509KeyManager.class), mock(X509KeyManager.class)
		});

		mockStatic(KeyManagerFactory.class);
		PowerMockito.when(KeyManagerFactory.getInstance("SunX509")).thenReturn(keyManagerFactory);

		RouterSslServerSocketFactory socketFactory = spy(new RouterSslServerSocketFactory());
		assertThat(socketFactory.getKeyManagers("JKS", null, "SunX509", "some-alias").length, equalTo(1));
	}
}
