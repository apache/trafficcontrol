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

import org.apache.traffic_control.traffic_router.secure.CertificateDataListener;
import org.apache.traffic_control.traffic_router.shared.DeliveryServiceCertificates;
import org.apache.traffic_control.traffic_router.shared.DeliveryServiceCertificatesMBean;
import org.apache.traffic_control.traffic_router.tomcat.TomcatLifecycleListener;
import org.apache.catalina.Lifecycle;
import org.apache.catalina.LifecycleEvent;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mockito;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import javax.management.MBeanServer;
import javax.management.ObjectName;
import java.lang.management.ManagementFactory;
import java.util.Arrays;

import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;

@RunWith(PowerMockRunner.class)
@PrepareForTest({TomcatLifecycleListener.class, ManagementFactory.class, LifecycleEvent.class})
@PowerMockIgnore("javax.management.*")
public class TomcatLifecycleListenerTest {
	@Before
	public void before() {
		PowerMockito.mockStatic(ManagementFactory.class);
	}

	@Test
	public void itIgnoresNonInitEvents() {
		Mockito.when(ManagementFactory.getPlatformMBeanServer()).thenThrow(new RuntimeException("invoked getPlatformMBeanServer"));

		Lifecycle lifecycle = Mockito.mock(Lifecycle.class);
		TomcatLifecycleListener tomcatLifecycleListener = new TomcatLifecycleListener();

		Arrays.asList(
			Lifecycle.AFTER_START_EVENT,
			Lifecycle.AFTER_STOP_EVENT,
			Lifecycle.BEFORE_START_EVENT,
			Lifecycle.BEFORE_STOP_EVENT,
			Lifecycle.PERIODIC_EVENT,
			Lifecycle.START_EVENT,
			Lifecycle.STOP_EVENT
		).forEach(s -> {
			tomcatLifecycleListener.lifecycleEvent(new LifecycleEvent(lifecycle,s, new Object()));
		});
	}

	@Test
	public void itRegistersBeanAndAddsListenerOnInit() throws Exception {
		MBeanServer mBeanServer = mock(MBeanServer.class);
		Mockito.when(ManagementFactory.getPlatformMBeanServer()).thenAnswer(invocationOnMock -> mBeanServer);

		CertificateDataListener certificateDataListener = mock(CertificateDataListener.class);

		TomcatLifecycleListener tomcatLifecycleListener = new TomcatLifecycleListener();
		tomcatLifecycleListener.setCertificateDataListener(certificateDataListener);

		LifecycleEvent lifecycleEvent = PowerMockito.mock(LifecycleEvent.class);
		PowerMockito.when(lifecycleEvent.getType()).thenReturn(Lifecycle.AFTER_INIT_EVENT);

		tomcatLifecycleListener.lifecycleEvent(lifecycleEvent);
		ObjectName name = new ObjectName(DeliveryServiceCertificatesMBean.OBJECT_NAME);
		verify(mBeanServer).registerMBean(new DeliveryServiceCertificates(), name);
		verify(mBeanServer).addNotificationListener(name,certificateDataListener, null, null);
	}
}
