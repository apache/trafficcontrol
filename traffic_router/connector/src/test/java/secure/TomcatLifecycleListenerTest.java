package secure;

import com.comcast.cdn.traffic_control.traffic_router.secure.CertificateDataListener;
import com.comcast.cdn.traffic_control.traffic_router.shared.DeliveryServiceCertificates;
import com.comcast.cdn.traffic_control.traffic_router.shared.DeliveryServiceCertificatesMBean;
import com.comcast.cdn.traffic_control.traffic_router.tomcat.TomcatLifecycleListener;
import org.apache.catalina.Lifecycle;
import org.apache.catalina.LifecycleEvent;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mockito;
import org.powermock.api.mockito.PowerMockito;
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
			Lifecycle.DESTROY_EVENT,
			Lifecycle.PERIODIC_EVENT,
			Lifecycle.START_EVENT,
			Lifecycle.STOP_EVENT
		).forEach(s -> {
			tomcatLifecycleListener.lifecycleEvent(new LifecycleEvent(lifecycle,s));
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
		PowerMockito.when(lifecycleEvent.getType()).thenReturn(Lifecycle.INIT_EVENT);

		tomcatLifecycleListener.lifecycleEvent(lifecycleEvent);
		ObjectName name = new ObjectName(DeliveryServiceCertificatesMBean.OBJECT_NAME);
		verify(mBeanServer).registerMBean(new DeliveryServiceCertificates(), name);
		verify(mBeanServer).addNotificationListener(name,certificateDataListener, null, null);
	}
}
