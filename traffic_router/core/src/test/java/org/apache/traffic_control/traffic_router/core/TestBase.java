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

import org.apache.traffic_control.traffic_router.core.external.HttpDataServer;
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
import java.io.File;
import java.lang.management.ManagementFactory;
import java.lang.reflect.Field;
import java.nio.file.Files;
import java.util.HashMap;
import java.util.Map;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.springframework.util.SocketUtils.findAvailableTcpPort;
import static org.springframework.util.SocketUtils.findAvailableUdpPort;

public class TestBase {
	private static final Logger LOGGER = LogManager.getLogger(TestBase.class);

	private static int testHttpServerPort = findAvailableTcpPort();
	private static HttpDataServer httpDataServer = new HttpDataServer(testHttpServerPort);
	private static File tmpDeployDir;
	private static ApplicationContext context;

	public static void setupFakeServers() throws Exception {
		// Set up a local server that can hand out
		// cr-config and cr-states (i.e. fake traffic monitor endpoints)
		// czmap
		// federations
		// steering
		// fake setting a cookie

		if (tmpDeployDir == null) {
			tmpDeployDir = Files.createTempDirectory("ext-test-").toFile();
		}

		final String TRAFFIC_MONITOR_BOOTSTRAP_LOCAL = "TRAFFIC_MONITOR_BOOTSTRAP_LOCAL";
		final String TRAFFIC_MONITOR_HOSTS = "TRAFFIC_MONITOR_HOSTS";
		String FAKE_SERVER;

		FAKE_SERVER = "localhost:" + testHttpServerPort + ";";

		Map<String, String> additionalEnvironment = new HashMap<>();

		additionalEnvironment.put(TRAFFIC_MONITOR_BOOTSTRAP_LOCAL, "true");
		additionalEnvironment.put(TRAFFIC_MONITOR_HOSTS, FAKE_SERVER);

		if (System.getenv(TRAFFIC_MONITOR_HOSTS) != null) {
			System.out.println("External Test Suite overriding env var [" + TRAFFIC_MONITOR_HOSTS + "] to " + FAKE_SERVER);
		}

		if (System.getenv(TRAFFIC_MONITOR_BOOTSTRAP_LOCAL) != null) {
			System.out.println("External Test Suite overriding env var [" + TRAFFIC_MONITOR_BOOTSTRAP_LOCAL + "] to true");
		}

		addToEnv(additionalEnvironment);

		assertThat(System.getenv(TRAFFIC_MONITOR_BOOTSTRAP_LOCAL), equalTo("true"));
		assertThat(System.getenv(TRAFFIC_MONITOR_HOSTS), equalTo(FAKE_SERVER));

		httpDataServer.start(testHttpServerPort);

		System.setProperty("testHttpServerPort", "" + testHttpServerPort);
		System.setProperty("routerHttpPort", "" + findAvailableTcpPort());
		System.setProperty("routerSecurePort", "" + findAvailableTcpPort());

		new File(tmpDeployDir,"conf").mkdirs();
		System.out.println();
		System.out.println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>");
		System.out.println(">>>>>>>> Going to use tmp directory '" + tmpDeployDir + "' as traffic router deploy directory");
		System.out.println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>");
		System.out.println();
		System.setProperty("deploy.dir", tmpDeployDir.getAbsolutePath());
		System.setProperty("dns.zones.dir", "src/test/var/auto-zones");

		System.setProperty("cache.health.json.refresh.period", "10000");
		System.setProperty("cache.config.json.refresh.period", "10000");
		System.setProperty("dns.tcp.port", "" + findAvailableTcpPort());
		System.setProperty("dns.udp.port", "" + findAvailableUdpPort());
		System.setProperty("traffic_monitor.properties", "not_needed");

		File dbDirectory = new File(tmpDeployDir, "db");
		dbDirectory.mkdir();

		LoggerContext.getContext().getLogger("org.eclipse.jetty").setLevel(Level.WARN);
		LoggerContext.getContext().getLogger("org.springframework").setLevel(Level.WARN);
		LoggerContext.getContext().getLogger("").setLevel(Level.WARN);

		final MBeanServer platformMBeanServer = ManagementFactory.getPlatformMBeanServer();
		try {
			final ObjectName objectName = new ObjectName(DeliveryServiceCertificatesMBean.OBJECT_NAME);
			platformMBeanServer.registerMBean(new DeliveryServiceCertificates(), objectName);
		} catch (Exception e) {
			e.printStackTrace();
		}

		ConsoleAppender consoleAppender = ConsoleAppender.newBuilder().setName("TestBase").setLayout(PatternLayout.newBuilder().withPattern("%d{ISO8601} [%-5p] %c{4}: %m%n").build()).build();
		LoggerContext.getContext().getRootLogger().addAppender(consoleAppender);
		LoggerContext.getContext().getRootLogger().setLevel(Level.INFO);
	}

	public static void addToEnv(Map<String, String> envVars) throws Exception {
		Map<String, String> envMap = System.getenv();
		Class<?> clazz = envMap.getClass();
		Field m = clazz.getDeclaredField("m");
		m.setAccessible(true);

		Map<String, String> mutableEnvMap = (Map<String, String>) m.get(envMap);
		mutableEnvMap.putAll(envVars);
	}

	public static void tearDownFakeServers() throws Exception {
		httpDataServer.stop();
		tmpDeployDir.deleteOnExit();
		final MBeanServer platformMBeanServer = ManagementFactory.getPlatformMBeanServer();
		try {
			final ObjectName objectName = new ObjectName(DeliveryServiceCertificatesMBean.OBJECT_NAME);
			platformMBeanServer.unregisterMBean(objectName);
		} catch (Exception e) {
			e.printStackTrace();
		}
	}

	public static ApplicationContext getContext() {
		if (context == null) {
			LOGGER.warn("Initializing context before running integration tests");
			context = new FileSystemXmlApplicationContext("src/main/webapp/WEB-INF/applicationContext.xml");
			LOGGER.warn("Context initialized integration tests will now start running");
		}
		return context;
	}
}
