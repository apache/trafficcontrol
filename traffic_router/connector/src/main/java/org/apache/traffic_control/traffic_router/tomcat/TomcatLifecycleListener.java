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

package org.apache.traffic_control.traffic_router.tomcat;

import org.apache.traffic_control.traffic_router.secure.CertificateDataListener;
import org.apache.traffic_control.traffic_router.shared.DeliveryServiceCertificates;
import org.apache.traffic_control.traffic_router.shared.DeliveryServiceCertificatesMBean;
import org.apache.catalina.Lifecycle;
import org.apache.catalina.LifecycleEvent;
import org.apache.catalina.LifecycleListener;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;


import javax.management.MBeanServer;
import javax.management.ObjectName;
import java.lang.management.ManagementFactory;

public class TomcatLifecycleListener implements LifecycleListener {
	private static final Logger log = LogManager.getLogger(LifecycleListener.class);
	private CertificateDataListener certificateDataListener = new CertificateDataListener();

	@Override
	@SuppressWarnings("PMD.AvoidThrowingRawExceptionTypes")
	public void lifecycleEvent(final LifecycleEvent event) {
		if (!Lifecycle.AFTER_INIT_EVENT.equals(event.getType())) {
			return;
		}

		try {
			log.info("Registering delivery service certificates mbean");
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
