package com.comcast.cdn.traffic_control.traffic_router.tomcat;

import com.comcast.cdn.traffic_control.traffic_router.secure.CertificateDataListener;
import com.comcast.cdn.traffic_control.traffic_router.shared.DeliveryServiceCertificates;
import com.comcast.cdn.traffic_control.traffic_router.shared.DeliveryServiceCertificatesMBean;
import org.apache.catalina.Lifecycle;
import org.apache.catalina.LifecycleEvent;
import org.apache.catalina.LifecycleListener;

import javax.management.MBeanServer;
import javax.management.ObjectName;
import java.lang.management.ManagementFactory;

public class TomcatLifecycleListener implements LifecycleListener {
	protected static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(TomcatLifecycleListener.class);
	private CertificateDataListener certificateDataListener = new CertificateDataListener();

	@Override
	@SuppressWarnings("PMD.AvoidThrowingRawExceptionTypes")
	public void lifecycleEvent(final LifecycleEvent event) {
		if (!Lifecycle.INIT_EVENT.equals(event.getType())) {
			return;
		}

		try {
			log.info("Registering delivery service certifcates mbean");
			final ObjectName objectName = new ObjectName(DeliveryServiceCertificatesMBean.OBJECT_NAME);

			final MBeanServer platformMBeanServer = ManagementFactory.getPlatformMBeanServer();
			platformMBeanServer.registerMBean(new DeliveryServiceCertificates(), objectName);
			platformMBeanServer.addNotificationListener(objectName, certificateDataListener, null, null);

		} catch (Exception e) {
			throw new RuntimeException("Failed to register MBean " + DeliveryServiceCertificatesMBean.OBJECT_NAME + " " + e.getClass().getSimpleName() + ": " + e.getMessage(), e);
		}
	}

	public CertificateDataListener getCertificateDataListener() {
		return certificateDataListener;
	}

	public void setCertificateDataListener(final CertificateDataListener certificateDataListener) {
		this.certificateDataListener = certificateDataListener;
	}
}
