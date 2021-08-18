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

import java.io.File;
import java.io.IOException;
import java.net.InetAddress;
import java.net.UnknownHostException;

import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import com.maxmind.geoip2.exception.GeoIp2Exception;

public class AnonymousIpDatabaseServiceTest {

	private AnonymousIpDatabaseService anonymousIpService;
	private final static String mmdb = "src/test/resources/GeoIP2-Anonymous-IP.mmdb";

	@Before
	public void setup() throws Exception {
		// ignore the test if there is no mmdb file
		File mmdbFile = new File(mmdb);
		org.junit.Assume.assumeTrue(mmdbFile.exists());

		anonymousIpService = new AnonymousIpDatabaseService();
		File databaseFile = new File(mmdb);
		anonymousIpService.setDatabaseFile(databaseFile);
		anonymousIpService.reloadDatabase();
		assert anonymousIpService.isInitialized();
	}

	@Test
	public void testIpInDatabase() throws Exception {
		assertThat(anonymousIpService.lookupIp(InetAddress.getByName("223.26.48.248")), notNullValue());
		assertThat(anonymousIpService.lookupIp(InetAddress.getByName("223.26.48.248")), notNullValue());
		assertThat(anonymousIpService.lookupIp(InetAddress.getByName("1.1.205.152")), notNullValue());
		assertThat(anonymousIpService.lookupIp(InetAddress.getByName("18.85.22.204")), notNullValue());
	}

	@Test
	public void testIpNotInDatabase() throws Exception {
		assertThat(anonymousIpService.lookupIp(InetAddress.getByName("192.168.0.1")), equalTo(null));
	}

	@Test
	public void testDatabaseNotLoaded() throws UnknownHostException, IOException, GeoIp2Exception {
		AnonymousIpDatabaseService anonymousIpService = new AnonymousIpDatabaseService();
		assertThat(anonymousIpService.isInitialized(), equalTo(false));
		assertThat(anonymousIpService.lookupIp(InetAddress.getByName("223.26.48.248")), equalTo(null));
		assertThat(anonymousIpService.lookupIp(InetAddress.getByName("192.168.0.1")), equalTo(null));
	}

	@Test
	public void testLookupTime() throws IOException {
		final InetAddress ipAddress = InetAddress.getByName("223.26.48.248");
		final long start = System.nanoTime();

		long total = 100000;

		for (int i = 0; i <= total; i++) {
			anonymousIpService.lookupIp(ipAddress);
		}

		long duration = System.nanoTime() - start;
		
		System.out.println(String.format("Anonymous IP database average lookup: %s nanoseconds", Long.toString(duration / total)));
	}

	@After
	public void tearDown() throws Exception {
		try {
			anonymousIpService.finalize();
		} catch (Throwable e) {
		}
	}

}
