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

import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.lang.reflect.Field;
import java.nio.file.FileVisitResult;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.SimpleFileVisitor;
import java.nio.file.attribute.BasicFileAttributes;
import java.security.KeyFactory;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.Base64;
import java.util.HashMap;
import java.util.Map;
import java.util.Properties;

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

		System.setProperty("traffic_monitor.properties", "src/test/conf/traffic_monitor.properties");

		X509Certificate x509Certificate = (X509Certificate) CertificateFactory.getInstance("X.509")
			.generateCertificate(new ByteArrayInputStream(Base64.getDecoder().decode(HTTPS_TEST_CERT)));

		final Properties properties = new Properties();
		properties.setProperty("keypass", "testing-testing");
		properties.store(new FileOutputStream(tmpDeployDir.getAbsolutePath() + "/conf/keystore.properties"), null);

		KeyStore keyStore = KeyStore.getInstance(KeyStore.getDefaultType());
		keyStore.load(null, "testing-testing".toCharArray());

		byte[] keyBytes = Base64.getDecoder().decode(HTTPS_TEST_KEY.getBytes());
		PKCS8EncodedKeySpec keySpec = new PKCS8EncodedKeySpec(keyBytes);
		KeyFactory fact = KeyFactory.getInstance("RSA");
		PrivateKey priv = fact.generatePrivate(keySpec);

		keyStore.setKeyEntry("https-test.thecdn.example.com", priv, "testing-testing".toCharArray(), new X509Certificate[] {x509Certificate});

		File dbDirectory = new File(tmpDeployDir, "db");
		dbDirectory.mkdir();
		File keystoreFile = new File(dbDirectory, ".keystore");
		keyStore.store(new FileOutputStream(keystoreFile), "testing-testing".toCharArray());

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

	// The following is just a self signed certificate and key to use for testing
	// this is NOT private data from a CA
	final static String HTTPS_TEST_CERT =
		"MIID8jCCAtoCCQCnNzY6/zYPcDANBgkqhkiG9w0BAQsFADCBujELMAkGA1UEBhMC" +
		"VVMxETAPBgNVBAgMCENvbG9yYWRvMQ8wDQYDVQQHDAZEZW52ZXIxFDASBgNVBAoM" +
		"C0V4YW1wbGUgSW5jMRkwFwYDVQQLDBBVbmljb3JuIFRyYWluZXJzMSgwJgYDVQQD" +
		"DB8qLmh0dHBzLXRlc3QudGhlY2RuLmV4YW1wbGUuY29tMSwwKgYJKoZIhvcNAQkB" +
		"Fh1vcGVyYXRpb25zQHRoZWNkbi5leGFtcGxlLmNvbTAeFw0xNjA2MzAyMDU3MTha" +
		"Fw0yNjA2MjgyMDU3MThaMIG6MQswCQYDVQQGEwJVUzERMA8GA1UECAwIQ29sb3Jh" +
		"ZG8xDzANBgNVBAcMBkRlbnZlcjEUMBIGA1UECgwLRXhhbXBsZSBJbmMxGTAXBgNV" +
		"BAsMEFVuaWNvcm4gVHJhaW5lcnMxKDAmBgNVBAMMHyouaHR0cHMtdGVzdC50aGVj" +
		"ZG4uZXhhbXBsZS5jb20xLDAqBgkqhkiG9w0BCQEWHW9wZXJhdGlvbnNAdGhlY2Ru" +
		"LmV4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAseev" +
		"nM90V6cz7/SI3U59EG23DoVFYoHA6Hrd98xBcipq1qyB0w5/IfYVwI94Pw03CJ+S" +
		"kui2Id4Zp3+aRT6gNxx61kdp1cUhpt1Y8lcUVIJrDQ3C3vjB6GwcKSXsFloKPNRM" +
		"HTw0xFw4gIB9Gx+TztabPs/R+hp/o9kyeuoD2foG7sPnDV9O2d5A8bPX8XZcrgN9" +
		"JtpTlwQ475sPoYDBVU2ncWmYyu7bPnfvPRPBncMRHjwK/bFlmK/+QNbPL0BA5zhd" +
		"qDHDNF05qgehB7EPge8nYkpFtKAjbiMZA1T/4Y1G/bLA9lWWqBEEP86bX2EqRg9v" +
		"aIOT2TO0H+Pc3t8qUwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCr82hSXa4Od5dd" +
		"AzU18s3wg/5eZhq+pqmlv71CujH/P54u0RkS4U0/4K8QBO3uFttymrPjn4oYH2sk" +
		"bKXGx8tvEnHCqLCOVY+G0CMAdbGJE1ztUuLByD0UC64Njat7pxyX1I5B8PTafvW3" +
		"7aqMlQQMGRAk/bQZaPtoEmArgUeE+fQoACFCjYKtTaP/p81wrefWZs4INKo1bATO" +
		"j30a3JwcWQunZDM2IMQwUPcjaXyfz3CRjTtZZ1ICGIyjKau8HXvYUnnoW63l1XYW" +
		"OIErXkbyeY4CS14wDbKAEffOPTdBWlpxbWI8m/Ox0y8mTAJQWcbSGuTiTzFSV3Bi" +
		"p4+a2PGN";

	static final String HTTPS_TEST_KEY =
		"MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCx56+cz3RXpzPv" +
		"9IjdTn0QbbcOhUVigcDoet33zEFyKmrWrIHTDn8h9hXAj3g/DTcIn5KS6LYh3hmn" +
		"f5pFPqA3HHrWR2nVxSGm3VjyVxRUgmsNDcLe+MHobBwpJewWWgo81EwdPDTEXDiA" +
		"gH0bH5PO1ps+z9H6Gn+j2TJ66gPZ+gbuw+cNX07Z3kDxs9fxdlyuA30m2lOXBDjv" +
		"mw+hgMFVTadxaZjK7ts+d+89E8GdwxEePAr9sWWYr/5A1s8vQEDnOF2oMcM0XTmq" +
		"B6EHsQ+B7ydiSkW0oCNuIxkDVP/hjUb9ssD2VZaoEQQ/zptfYSpGD29og5PZM7Qf" +
		"49ze3ypTAgMBAAECggEANPd92XoKcW5ekDqF5R3RLmr77V7QYZuwC4dJPtWZRpxK" +
		"Ys/Jd0UBpOLXZxVP/7W3hOG9ie+vCjZN/QiIrcUPflkEWXe5kuO2OS/9o2k5rE+H" +
		"/8LxGeGGGeTWHYok2CLGmYW7g5jBPRUX3Wpj1Qd5wkMyxWiqY4QwLGAmH2I881B1" +
		"6DC3v7XeFF5eG5WlOQdg+4l7D283Zp4Jlq5K31WNyRarplXvEE+DOpd3iGIzGI+w" +
		"0V3DdcbKT2pnqP4mpz6dVFC9xTBONX6nZHmPuIoZcEql49kKD833fU9WtnJ3bllM" +
		"HRHnQDtyaCzVXOMCA2hHzNfZH60iLB2Wo19niSPMgQKBgQDoN+PBhLlBxamo6mi0" +
		"qDtMpAq+n9PzMctYDyxt23xFxZy+wB2htAGxhnIKXcFZpE0Tbg8Y+HUFiag9sDh1" +
		"3bhhLioki0QL1DPpVdlOdADkqnOXCcvepmYMZNeUpHA3R34lil3hXWx167LWxJ80" +
		"rAPlPJq2HaTXsZzOXx4o1ee3mQKBgQDEH9rLXdsV3cb7ZqdjVGazNxKVRAemv2rW" +
		"UZTsHfgikndu3z+4B++H1176ZC4LQfPSOG8B7ZcPeN/GMDAk+iJRCth0pzAB0mZx" +
		"I/3GiAN7ACdosKzSoSwI0tXobXfBxMTJMx8ToUWrPKPncxkp8vDNhid/x7ENe6qY" +
		"EbkR3PG0ywKBgQCs2Z0wWKjE6mqlDwatImQxYhGVXsaXSUNA4tqBU1SnYraPzdTA" +
		"nop8J8UPLkZTgVbV1aBrR9VjL9oJQPhl04oA3CoGVZtq6qNRVdOQ8AwSKUYs8N/N" +
		"dTKUmyNUwym8G/0r2FiU/cNT6wONlYGj5T5pDbljQaGH4+8CNg7u+nmmUQKBgFvw" +
		"Gt7utn8/ocHEU3+K10H39Swn4fZXETw6rjcprWJ3iqlc2j/o6G6jlZCHWdZJKoVH" +
		"kzIyMHg+T5hWipsq7t9S2DmHDkgsW316Q8LHi+ojHlZDTCDJER1pyIDWoCcjmKRA" +
		"5LaNCV3GZYdgO1Gg4yVVWDrcX7FUYZo75KftDRmVAoGAa8OJerRIBCW6NVcVzILW" +
		"IZaqRiHZaBEl3fy1tQ8QznMdwr5Z/RqD/R5uss0H8rmMU9TqZU7864nMYKaRAb3/" +
		"pzEyKUZJcQO+o0rWMVM3GcTmVcg/z3kl7eo1ouYlvPkpWxbFQlvXbA6WMPmwqq1O" +
		"tTHrGbnpc7Ey4W2MvXSvLMo=";
}
