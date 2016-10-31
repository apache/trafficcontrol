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
