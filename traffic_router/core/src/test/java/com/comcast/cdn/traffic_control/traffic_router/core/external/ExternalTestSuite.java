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
	public static String FAKE_SERVER;
	private static CatalinaTrafficRouter catalinaTrafficRouter;
	private static HttpDataServer httpDataServer;
	private static File tmpDeployDir;
	private static int testHttpServerPort;

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

		httpDataServer = new HttpDataServer(testHttpServerPort);
		httpDataServer.start(testHttpServerPort);
	}

	@BeforeClass
	public static void beforeClass() throws Exception {
		testHttpServerPort = findAvailableTcpPort();

		System.setProperty("testHttpServerPort", "" + testHttpServerPort);
		System.setProperty("routerHttpPort", "" + findAvailableTcpPort());
		System.setProperty("routerSecurePort", "" + findAvailableTcpPort());

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

		KeyFactory keyFactory = KeyFactory.getInstance("RSA");

		byte[] keyBytes = Base64.getDecoder().decode(HTTPS_TEST_KEY.getBytes());
		PKCS8EncodedKeySpec keySpec = new PKCS8EncodedKeySpec(keyBytes);
		PrivateKey privateKey = keyFactory.generatePrivate(keySpec);

		keyStore.setKeyEntry("https-only-test.thecdn.example.com", privateKey, "testing-testing".toCharArray(), new X509Certificate[] {x509Certificate});

		x509Certificate = (X509Certificate) CertificateFactory.getInstance("X.509")
			.generateCertificate(new ByteArrayInputStream(Base64.getDecoder().decode(HTTP_AND_HTTPS_TEST_CERT)));

		keyBytes = Base64.getDecoder().decode(HTTP_AND_HTTPS_TEST_KEY.getBytes());
		keySpec = new PKCS8EncodedKeySpec(keyBytes);
		privateKey = keyFactory.generatePrivate(keySpec);

		keyStore.setKeyEntry("http-and-https-test.thecdn.example.com", privateKey, "testing-testing".toCharArray(), new X509Certificate[] {x509Certificate});

		x509Certificate = (X509Certificate) CertificateFactory.getInstance("X.509")
			.generateCertificate(new ByteArrayInputStream(Base64.getDecoder().decode(HTTP_TO_HTTPS_TEST_CERT)));

		keyBytes = Base64.getDecoder().decode(HTTP_TO_HTTPS_TEST_KEY.getBytes());
		keySpec = new PKCS8EncodedKeySpec(keyBytes);
		privateKey = keyFactory.generatePrivate(keySpec);

		keyStore.setKeyEntry("http-to-https-test.thecdn.example.com", privateKey, "testing-testing".toCharArray(), new X509Certificate[] {x509Certificate});

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

	// The following is just a self signed certificate and key to use for testing
	// this is NOT private data from a CA
	// *.http-and-https-test.thecdn.example.com
	//
	// openssl req -nodes -newkey rsa:2048 -keyout ~/tmp/http-and-https.key -out ~/tmp/http-and-https.csr
	//
	// openssl x509 -req -days 3650 -in ~/tmp/http-and-https.csr -signkey ~/tmp/http-and-https.key -out ~/tmp/http-and-https.crt

	final static String HTTP_AND_HTTPS_TEST_KEY =
		"MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDjhGSGLMVtaY32" +
		"aS7aBogJCVmWNb6esx7+W6ug/wwYgwrsCL0nl+J6snPBOG4HoHHU5pKisHVYAbUi" +
		"3TUBgjeP/uCGxKIonjru5cbS7tRIoTqSX/PJOlm35i0sJ55A3UHZafH8C0xnSzQj" +
		"Ti9Evot58wCza/zfHKp/01Ig0M38+BhU7jpDzEgNEbYfUVpgAkES0JoJBAvrzqKB" +
		"IZHXtp9hMy6uG8cTOybFDZZGJfwfUpngWqQGfT2h91ah49gVbAlef//647EuW1Dw" +
		"X7PQKb0rc1bPYtlaVfeEZNb+e+52hs/B/hL2rXJ++G2t+2skahi5Kq7bcycQGjr+" +
		"OROkA5UbAgMBAAECggEBALqsif49Bc/752rianqhGUSw0zyX5Es6FJgGhw+VtEr4" +
		"WiHIGcs+p6icerVyo3TGhB92/6FUvzLyU7jDXxZZzVTsfzSUaaiCC0Cwby3qn2ro" +
		"PrKS3+efZLWqui2cZBA8eib08oMmkg2+eoztPYNeA/qPE2gjlltJnes7bAtYx2pi" +
		"aQK5stsynOWyxYdDdmqC7VJxCtC9sTmubaAQXHBBzrQnmMt7XzYsc6X/WhbWIFbs" +
		"CxvC6K0CSUD0DzPUz2eN8QERMQQZQMysq+DA73/G9Zrl4LFL4qLTooviVgzq1gjx" +
		"5yzKYR93CNk0HFIJzREw/FPm9xKa6tH5ZnwpV6FVRSECgYEA+EfzX/vWGqTWH05X" +
		"y1p1iNmNst9PbPVdJx/SqiitU0RsF+A2g1y/THEt67d6l4CD6OlsK/yTNHRr/dMm" +
		"L9B0SLVYUNGNnphl3WF10VQ2Y+80dRyWSbJpqty7L3P3uLOmkZRX1/9X254ryVEj" +
		"n6KmkKK056u7RZoQw1ZXMK8RkVUCgYEA6pcvCwJBgzIwc9LufIpLa7NurX249Bbp" +
		"9B9LYS1vfS2GYAypZfvIwqUi8jAjH2SIaVzI7q3mzocn+lV93ZuvU/dHjYs1VTC3" +
		"nW2G1sTk3nkSlrWnH0mDpkta9UK3/nD4gOmZmHD2rPyAvzj+RE3EAB0lzQGV9Squ" +
		"aztpg7BsTK8CgYEA68xZxhUFmLRob78V/png+qGzw+f2JQM6/0dn6hdL1cMr7dkR" +
		"rNzPCiiLdk0BbxWtMe1OwM/WdoEDd0OsBskxR0SDpe3/VFpklEZVgQM7zNmHtpn5" +
		"2fBKDu4oEL9Qy+hDEAwVCZ0GshucdkxLSvdMvhzpNwWQjF/v/7TmheQfCSkCgYBM" +
		"hdiAnNHF/B82CP5mfa4wia12xmYIqVjTm0m5f1q42JrWxgqUC9fnNnr5yZ4LZX3h" +
		"8LRSt0Ns50WxMSYHnftJRoZ+s4RIL8YVgl7TvBJ0R8Y6hzLmz9Iz8qzPCF6Aj1Vg" +
		"p9LEmUS+FPfiaLL4kO14pAlqoDPMb4nJzO2UWX5aXQKBgAmnvhj/aLcJnCJM0YnC" +
		"/aRWTF5Q3HQmPOHx5fQlw9+hCjQUkaoPL5JVs4/Z/dOj1RsWYHg23fGy78zNHkQi" +
		"zL6P2WpZ7pEpJbK4wobpfitzczKfNZROAzdJPDV4+ebtPHMkGA/ibN04AM2SWKTH" +
		"UoGXvsZbRbb+j3EptEHBiNiN";

	final static String HTTP_AND_HTTPS_TEST_CERT =
		"MIIDqjCCApICCQDx6373gd/QFDANBgkqhkiG9w0BAQsFADCBljELMAkGA1UEBhMC" +
		"VVMxETAPBgNVBAgMCENvbG9yYWRvMQ8wDQYDVQQHDAZEZW52ZXIxFTATBgNVBAoM" +
		"DEh1YmNhcHMgUiBVczEZMBcGA1UECwwQSHViY2FwIFBvbGlzaGVyczExMC8GA1UE" +
		"AwwoKi5odHRwLWFuZC1odHRwcy10ZXN0LnRoZWNkbi5leGFtcGxlLmNvbTAeFw0x" +
		"NjA4MDkyMTI0NDdaFw0yNjA4MDcyMTI0NDdaMIGWMQswCQYDVQQGEwJVUzERMA8G" +
		"A1UECAwIQ29sb3JhZG8xDzANBgNVBAcMBkRlbnZlcjEVMBMGA1UECgwMSHViY2Fw" +
		"cyBSIFVzMRkwFwYDVQQLDBBIdWJjYXAgUG9saXNoZXJzMTEwLwYDVQQDDCgqLmh0" +
		"dHAtYW5kLWh0dHBzLXRlc3QudGhlY2RuLmV4YW1wbGUuY29tMIIBIjANBgkqhkiG" +
		"9w0BAQEFAAOCAQ8AMIIBCgKCAQEA44RkhizFbWmN9mku2gaICQlZljW+nrMe/lur" +
		"oP8MGIMK7Ai9J5fierJzwThuB6Bx1OaSorB1WAG1It01AYI3j/7ghsSiKJ467uXG" +
		"0u7USKE6kl/zyTpZt+YtLCeeQN1B2Wnx/AtMZ0s0I04vRL6LefMAs2v83xyqf9NS" +
		"INDN/PgYVO46Q8xIDRG2H1FaYAJBEtCaCQQL686igSGR17afYTMurhvHEzsmxQ2W" +
		"RiX8H1KZ4FqkBn09ofdWoePYFWwJXn//+uOxLltQ8F+z0Cm9K3NWz2LZWlX3hGTW" +
		"/nvudobPwf4S9q1yfvhtrftrJGoYuSqu23MnEBo6/jkTpAOVGwIDAQABMA0GCSqG" +
		"SIb3DQEBCwUAA4IBAQBpO3jPVhDvFPJZJmzFbaC2vT/yq1oPtn9Z29bvkz9UTOc8" +
		"aItDK84KjbuUZ3i9ol1AWu6tWQRitfnxpkhKDEMXaOZq/HBMrz4XPHC+2Ez/+lOU" +
		"SmwAQHaaQMS20/9TAtNjIBvwphFpXXeT6Iz2NZl2EYEVdIfbQkTW0UsoFzBZGn3S" +
		"/0OXhd1lRXt0lH8glYEkL35FQJ0PCIM5W4mRJ50FKTI1x52xFY44ctEtGYkrGeWZ" +
		"4xYU0pTLKEYET0vKBHkjcvevI7dTd7caaWIXu4WG6ToVz8suTiKH49dMd3ev0qM7" +
		"qnx67ypmcnGqRRLxC6F2gnMx8B8sJ37QQlEYBMQo";

	// The following is just a self signed certificate and key to use for testing
	// this is NOT private data from a CA
	// *.http-to-https-test.thecdn.example.com
	//
	// openssl req -nodes -newkey rsa:2048 -keyout ~/tmp/http-to-https.key -out ~/tmp/http-to-https.csr
	//
	// openssl x509 -req -days 3650 -in ~/tmp/http-to-https.csr -signkey ~/tmp/http-to-https.key -out ~/tmp/http-to-https.crt

	final static String HTTP_TO_HTTPS_TEST_KEY =
		"MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDKu4aRsrexzyck" +
		"+IxuotTSfvt2YLhrhFRscSLQdo56+RE402e8FuJ4DONwPRxbNdjL+E7elLg+onOf" +
		"gB9sYkzzFIy4CDp5qoGkencFNJrDwJ+KGxWzxTyPsy3rTGSx4hHj0o1AuE0k1Tgx" +
		"A8XWtG7UE20iV6Gj5bzSqDLOR/fvsQESSdFFHKArO0fQcd8Z3LVdczv+3To10Mhu" +
		"+zdltzjE7v4A38ewKVgFBk1BAOxc7y/ytP3iUH2YS99H3jB61Ej18e+VYe5oaw60" +
		"/4r2PG0FlXORxJuB9rC7OydXvu6fvOE6lmcImnyXUrBf8C+bHgO4CWG6thOVW4el" +
		"vOo8dFP3AgMBAAECggEAEFxr8swyiPYH2bL5WmBnvoki8B3EJGEskwfaYGqA+ymo" +
		"myZsg8BxDHE11bQI2s+QrH1gmBP2fo+Ltz6Wyp9wSFnLNXrshS8egVCk1FW3e77K" +
		"4VFoQfbT+WDjfs7OfZCaEwHGBogZKbTPcR011SsAmrrqns/lqp16zKFoYD9sofpJ" +
		"AZHL6Biu9PTfob0W8Co6thiii1xn+TEdc1ESDYdkYM5xsphrLoYyM7n1VyXRl31g" +
		"sewofW/ArF4K0Vl5iGygRKPw+Izqq4iCSqTzVr1T4Eh56k0cW0opOgww/LdybyGq" +
		"EqvczqHkj0sjHX9WKbTkNGAcymCUAVyaCf4g8Upn0QKBgQD7c8zBPhj1NO42I+yJ" +
		"+SkZKg24zudJb+ztjeBFg28Vg8n13xQIgHHMtIDaC8G/5vgrS0WFFVZuYTSLU2R3" +
		"b954H65c+L5N1mDAD3EDE73+xHfbm8dEeJVeGK59x1CgkGbLtgiaf/d466KhOiQj" +
		"xlsBEkByLIXfmrxXYZH54xD1GQKBgQDOZihiRKZ9oGUlGh4CWO/gh3RjrXhqxDES" +
		"9OzMrpEJQLe3Af1rHUHkL1ugjkykYwqD8AvKnsoJ+2Bbri5dtmTE0f6R3K5QP8vH" +
		"pShnFTxU6Q3/meDxwnIX6a5AfLXJsyxmVV1fmsD3UayN7lrAtWT4CNlicFrHUZJL" +
		"S18epmEjjwKBgGZyRpDQyQBWQVtjhYKtNfZfsNmDyq2b4U7jx+TqaL6+Q/Fdot7X" +
		"3gWF4R11Psn9w0x4TWmsSNuN1QeSwVL8DAqq9bJBUd+KoT5+zA9x4q3CxAaAUE5w" +
		"RoLg0W7DXvEcBBWpI5Y23s+wSUEg3AqLTRaBpioeQ6jXdTawtPW3cng5AoGBAK2X" +
		"nj+IHb9rN6aM4NB4nMfrJSjwrWaeu+eFt+Quri1qERoKwmlkohaY/id7h1p7Mkzl" +
		"iAVSp/rdQZ3aUYTf8sDXHZTwVmuIPIwdjG2mnqeLnApuEZNER1F1aOkz+nE6EQ3A" +
		"nlfagJGCT+7PmeSaq+ExECSK+s7I/JH3Qnk01l5hAoGAWSa7fzLS57XFTHTTYddt" +
		"tK5W6hlKwEb/tBnI8iLnWa+KhmTo/VPsc1C4rV3FqVFfMN6ZHMCEdG/Hq9vQdkQZ" +
		"35crLobjKIk5tlVzEbWxwl8EQez180r0O1VsRIiceIlzXRl3I17GeEKQHaORx/wS" +
		"PkZQNvkYw/OLHPViXWBGCsQ=";

	final static String HTTP_TO_HTTPS_TEST_CERT =
		"MIIDpjCCAo4CCQDLCWeLJrqPvDANBgkqhkiG9w0BAQsFADCBlDELMAkGA1UEBhMC" +
		"VVMxETAPBgNVBAgMCENvbG9yYWRvMQ8wDQYDVQQHDAZEZW52ZXIxGDAWBgNVBAoM" +
		"D1Ntb2tlIEFuZCBGbGFtZTEaMBgGA1UECwwRU3BsaW50ZXIgVHJpbW1lcnMxKzAp" +
		"BgNVBAMMIiouaHR0cC10by1odHRwcy50aGVjZG4uZXhhbXBsZS5jb20wHhcNMTYw" +
		"ODA5MTgyMDA1WhcNMjYwODA3MTgyMDA1WjCBlDELMAkGA1UEBhMCVVMxETAPBgNV" +
		"BAgMCENvbG9yYWRvMQ8wDQYDVQQHDAZEZW52ZXIxGDAWBgNVBAoMD1Ntb2tlIEFu" +
		"ZCBGbGFtZTEaMBgGA1UECwwRU3BsaW50ZXIgVHJpbW1lcnMxKzApBgNVBAMMIiou" +
		"aHR0cC10by1odHRwcy50aGVjZG4uZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEB" +
		"AQUAA4IBDwAwggEKAoIBAQDKu4aRsrexzyck+IxuotTSfvt2YLhrhFRscSLQdo56" +
		"+RE402e8FuJ4DONwPRxbNdjL+E7elLg+onOfgB9sYkzzFIy4CDp5qoGkencFNJrD" +
		"wJ+KGxWzxTyPsy3rTGSx4hHj0o1AuE0k1TgxA8XWtG7UE20iV6Gj5bzSqDLOR/fv" +
		"sQESSdFFHKArO0fQcd8Z3LVdczv+3To10Mhu+zdltzjE7v4A38ewKVgFBk1BAOxc" +
		"7y/ytP3iUH2YS99H3jB61Ej18e+VYe5oaw60/4r2PG0FlXORxJuB9rC7OydXvu6f" +
		"vOE6lmcImnyXUrBf8C+bHgO4CWG6thOVW4elvOo8dFP3AgMBAAEwDQYJKoZIhvcN" +
		"AQELBQADggEBAKREwCYFiz858Iqsf+m/rkQErTVeSUPg6KSlDDknVI/x+x0uCwXN" +
		"OgGo5s2S6Ec0V8hd9PrADasCDtAGaLJ2giNEyv/0iZRcTfR2mfnKClZcVbEgvhqt" +
		"1e6oQ1ybKw+fsvSWOu8h30CiKjct4+gWjoSbyVgmHFSBdKvZJwice2ewi2SE+H4y" +
		"ekPD6BptIJQc6UfFE4ZuO7S7ajroWU7dVeI495Q8BQ89LWPgUwc/a90VrICAT9bB" +
		"VM+SiCpEFStMFlz/bEkm9goZJKroaPwXVf75hEicAaAPFs5zlQpfh4LOF+Gk0P/G" +
		"WNmQ5qbTdEyM1vxgM/4anoOFfaHhB4Pk8T4=";
}
