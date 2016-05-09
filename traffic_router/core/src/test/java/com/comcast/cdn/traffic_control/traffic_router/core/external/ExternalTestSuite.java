/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.external;

import com.comcast.cdn.traffic_control.traffic_router.core.CatalinaTrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.util.ExternalTest;
import org.apache.catalina.LifecycleException;
import org.apache.log4j.ConsoleAppender;
import org.apache.log4j.Level;
import org.apache.log4j.LogManager;
import org.apache.log4j.PatternLayout;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.experimental.categories.Category;
import org.junit.runner.RunWith;
import org.junit.runners.Suite;

import java.io.File;
import java.io.IOException;
import java.lang.reflect.Field;
import java.nio.file.FileVisitResult;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.SimpleFileVisitor;
import java.nio.file.attribute.BasicFileAttributes;
import java.util.HashMap;
import java.util.Map;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.springframework.util.SocketUtils.findAvailableTcpPort;
import static org.springframework.util.SocketUtils.findAvailableUdpPort;

@Category(ExternalTest.class)
@RunWith(Suite.class)
@Suite.SuiteClasses({
	SteeringTest.class,
	ConsistentHashTest.class,
	CoverageZoneTest.class,
	DeliveryServicesTest.class,
	LocationsTest.class,
	RouterTest.class,
	StatsTest.class,
	ZonesTest.class
})
public class ExternalTestSuite {
	public static final String TRAFFIC_MONITOR_BOOTSTRAP_LOCAL = "TRAFFIC_MONITOR_BOOTSTRAP_LOCAL";
	public static final String TRAFFIC_MONITOR_HOSTS = "TRAFFIC_MONITOR_HOSTS";
	public static final String FAKE_SERVER = "localhost:8889;";
	private static CatalinaTrafficRouter catalinaTrafficRouter;
	private static HttpDataServer httpDataServer;
	private static File tmpDeployDir;

	@SuppressWarnings("unchecked")
	public static void addToEnv(Map<String, String> envVars) throws Exception {
		Map<String, String> envMap = System.getenv();
		Class<?> clazz = envMap.getClass();
		Field m = clazz.getDeclaredField("m");
		m.setAccessible(true);

		Map<String, String> mutableEnvMap = (Map<String, String>) m.get(envMap);
		mutableEnvMap.putAll(envVars);
	}

	public static void setupFakeServers() throws Exception {
		// Set up a local server that can hand out
		// cr-config and cr-states (i.e. fake traffic monitor endpoints)
		// czmap
		// federations
		// steering
		// fake setting a cookie

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

		httpDataServer = new HttpDataServer();
		httpDataServer.start(8889);
	}

	@BeforeClass
	public static void beforeClass() throws Exception {
		setupFakeServers();

		tmpDeployDir = Files.createTempDirectory("ext-test-").toFile();
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

		LogManager.getLogger("org.eclipse.jetty").setLevel(Level.WARN);
		LogManager.getLogger("org.springframework").setLevel(Level.WARN);

		ConsoleAppender consoleAppender = new ConsoleAppender(new PatternLayout("%d{ISO8601} [%-5p] %c{4}: %m%n"));
		LogManager.getRootLogger().addAppender(consoleAppender);
		LogManager.getRootLogger().setLevel(Level.WARN);

		catalinaTrafficRouter = new CatalinaTrafficRouter("src/main/opt/tomcat/conf/server.xml", "src/main/webapp");
		catalinaTrafficRouter.start();
	}

	@AfterClass
	public static void afterClass() throws LifecycleException, IOException {
		catalinaTrafficRouter.stop();
		httpDataServer.stop();
		tmpDeployDir.deleteOnExit();

		Files.walkFileTree(tmpDeployDir.toPath(), new SimpleFileVisitor<Path>() {
			@Override
			public FileVisitResult visitFile(Path path, BasicFileAttributes attrs) {
				path.toFile().delete();
				return FileVisitResult.CONTINUE;
			}

			@Override
			public FileVisitResult postVisitDirectory(Path path, IOException e) {
				path.toFile().delete();
				return FileVisitResult.CONTINUE;
			}
		});
	}
}
