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

package org.apache.traffic_control.traffic_router.core;

import org.apache.traffic_control.traffic_router.shared.DeliveryServiceCertificates;
import org.apache.traffic_control.traffic_router.shared.DeliveryServiceCertificatesMBean;
import org.apache.logging.log4j.core.appender.ConsoleAppender;
import org.apache.logging.log4j.core.layout.PatternLayout;
import org.apache.logging.log4j.core.LoggerContext;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.context.ApplicationContext;
import org.springframework.context.support.FileSystemXmlApplicationContext;

import javax.management.MBeanServer;
import javax.management.ObjectName;
import java.lang.management.ManagementFactory;

import static org.springframework.util.SocketUtils.findAvailableTcpPort;
import static org.springframework.util.SocketUtils.findAvailableUdpPort;

public class TestBase {
	private static final Logger LOGGER = LogManager.getLogger(TestBase.class);
	public static final String monitorPropertiesPath = "src/test/conf/traffic_monitor.properties";
	private static ApplicationContext context;

	public static ApplicationContext getContext() {

		System.setProperty("deploy.dir", "src/test");
		System.setProperty("dns.zones.dir", "src/test/var/auto-zones");

		System.setProperty("dns.tcp.port", String.valueOf(findAvailableTcpPort()));
		System.setProperty("dns.udp.port", String.valueOf(findAvailableUdpPort()));

		if (context != null) {
			return context;
		}

		final MBeanServer platformMBeanServer = ManagementFactory.getPlatformMBeanServer();
		try {
			final ObjectName objectName = new ObjectName(DeliveryServiceCertificatesMBean.OBJECT_NAME);
			platformMBeanServer.registerMBean(new DeliveryServiceCertificates(), objectName);
		} catch (Exception e) {
			e.printStackTrace();
		}

		ConsoleAppender consoleAppender = ConsoleAppender.newBuilder().setLayout(PatternLayout.newBuilder().withPattern("%d{ISO8601} [%-5p] %c{4}: %m%n").build()).build();
		LoggerContext.getContext().getRootLogger().addAppender(consoleAppender);
		LoggerContext.getContext().getRootLogger().setLevel(Level.INFO);

		LOGGER.warn("Initializing context before running integration tests");
		context = new FileSystemXmlApplicationContext("src/main/webapp/WEB-INF/applicationContext.xml");
		LOGGER.warn("Context initialized integration tests will now start running");
		return context;
	}

}
