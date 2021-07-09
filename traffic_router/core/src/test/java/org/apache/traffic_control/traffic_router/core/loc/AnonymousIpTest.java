/*
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

package org.apache.traffic_control.traffic_router.core.loc;

import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

import java.io.File;
import java.io.IOException;

import org.junit.Before;
import org.junit.Test;

import org.apache.traffic_control.traffic_router.core.router.TrafficRouter;

public class AnonymousIpTest {
	private TrafficRouter trafficRouter;

	final File configFile = new File("src/test/resources/anonymous_ip.json");
	final File configNoWhitelist = new File("src/test/resources/anonymous_ip_no_whitelist.json");

	final String mmdb = "src/test/resources/GeoIP2-Anonymous-IP.mmdb";
	File databaseFile = new File(mmdb);

	@Before
	public void setUp() throws Exception {
		// ignore the test if there is no mmdb file
		File mmdbFile = new File(mmdb);
		org.junit.Assume.assumeTrue(mmdbFile.exists());

		AnonymousIp.parseConfigFile(configFile, false);
		assert (AnonymousIp.getCurrentConfig().getIPv4Whitelist() != null);
		assert (AnonymousIp.getCurrentConfig().getIPv6Whitelist() != null);

		// Set up a mock traffic router with real database
		AnonymousIpDatabaseService anonymousIpService = new AnonymousIpDatabaseService();
		anonymousIpService.setDatabaseFile(databaseFile);
		anonymousIpService.reloadDatabase();
		assert anonymousIpService.isInitialized();
		trafficRouter = mock(TrafficRouter.class);
		when(trafficRouter.getAnonymousIpDatabaseService()).thenReturn(anonymousIpService);
		assert (trafficRouter.getAnonymousIpDatabaseService() != null);
	}

	@Test
	public void testConfigFileParsingIpv4() {
		AnonymousIp currentConfig = AnonymousIp.getCurrentConfig();
		assertThat(currentConfig, notNullValue());
		AnonymousIpWhitelist whitelist = currentConfig.getIPv4Whitelist();
		assertThat(whitelist, notNullValue());
	}

	@Test
	public void testConfigFileParsingIpv6() {
		AnonymousIp currentConfig = AnonymousIp.getCurrentConfig();
		assertThat(currentConfig, notNullValue());
		AnonymousIpWhitelist whitelist = currentConfig.getIPv6Whitelist();
		assertThat(whitelist, notNullValue());
	}

	@Test
	public void testIpInWhitelistIsAllowed() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "5.34.32.79";

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);

		assertThat(result, equalTo(false));
	}

	@Test
	public void testFallsUnderManyPolicies() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "2.38.158.142";

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);

		assertThat(result, equalTo(true));
	}

	@Test
	public void testAllowNotCheckingPolicy() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "2.36.248.52";

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);

		assertThat(result, equalTo(false));
	}

	@Test
	public void testEnforceAllowed() throws IOException {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "10.0.0.1";

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);

		assertThat(result, equalTo(false));
	}

	@Test
	public void testEnforceAllowedIpInWhitelist() throws IOException {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "10.0.2.1";

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);

		assertThat(result, equalTo(false));
	}

	@Test
	public void testEnforceBlocked() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "223.26.48.248";

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);

		assertThat(result, equalTo(true));
	}

	@Test
	public void testEnforceNotInWhitelistNotInDB() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "192.168.0.1";

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);

		assertThat(result, equalTo(false));
	}

	/* IPv4 no whitelist */

	@Test
	public void testEnforceNoWhitelistAllowed() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "192.168.0.1";
		AnonymousIp.parseConfigFile(configNoWhitelist, false);

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);
		assertThat(result, equalTo(false));
	}

	@Test
	public void testEnforceNoWhitelistBlocked() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "223.26.48.248";
		AnonymousIp.parseConfigFile(configNoWhitelist, false);

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);
		assertThat(result, equalTo(true));
	}

	@Test
	public void testEnforceNoWhitelistNotEnforcePolicy() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "2.36.248.52";
		AnonymousIp.parseConfigFile(configNoWhitelist, false);

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);
		assertThat(result, equalTo(false));
	}

	/* IPv6 Testing */

	@Test
	public void testIpv6EnforceBlock() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "2001:418:9807::1";

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);
		assertThat(result, equalTo(true));
	}

	@Test
	public void testIpv6EnforceNotBlock() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "2001:418::1";

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);
		assertThat(result, equalTo(false));
	}

	@Test
	public void testIpv6EnforceNotBlockWhitelisted() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "2001:550:90a:0:0:0:0:1";

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);
		assertThat(result, equalTo(false));
	}

	@Test
	public void testIpv6EnforceNotBlockOnWhitelist() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "::1";

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);
		assertThat(result, equalTo(false));
	}

	/* IPv6 tests no whitelist */

	@Test
	public void testIpv6NoWhitelistEnforceBlock() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "2001:418:9807::1";
		AnonymousIp.parseConfigFile(configNoWhitelist, false);

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);
		assertThat(result, equalTo(true));
	}

	@Test
	public void testIpv6NoWhitelistNoBlock() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "::1";
		AnonymousIp.parseConfigFile(configNoWhitelist, false);

		final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);
		assertThat(result, equalTo(false));
	}

	@Test
	public void testAnonymousIpPerformance() {
		final String dsvcId = "dsID";
		final String url = "http://ds1.example.com/live1";
		final String ip = "2.36.248.52";

		long total = 100000;

		long start = System.nanoTime();

		for (int i = 0; i <= total; i++) {
			final boolean result = AnonymousIp.enforce(trafficRouter, dsvcId, url, ip);
		}

		long duration = System.nanoTime() - start;

		System.out.println(String.format("Anonymous IP blocking average took %s nanoseconds", Long.toString(duration / total)));
	}
}
