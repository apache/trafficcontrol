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
import static org.hamcrest.MatcherAssert.assertThat;

import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.Before;
import org.junit.Test;

import java.io.IOException;

public class AnonymousIpWhitelistTest {

	AnonymousIpWhitelist ip4whitelist;
	AnonymousIpWhitelist ip6whitelist;

	@Before
	public void setup() throws IOException, JsonUtilsException, NetworkNodeException {
		final ObjectMapper mapper = new ObjectMapper();

		ip4whitelist = new AnonymousIpWhitelist();
		ip4whitelist.init(mapper.readTree("[\"192.168.30.0/24\", \"10.0.2.0/24\", \"10.0.0.0/16\"]"));

		ip6whitelist = new AnonymousIpWhitelist();
		ip6whitelist.init(mapper.readTree("[\"::1/32\", \"2001::/64\"]"));
	}

	@Test
	public void testAnonymousIpWhitelistConstructor() {
		// final InetAddress address = InetAddresses.forString("192.168.30.1");
		assertThat(ip4whitelist.contains("192.168.30.1"), equalTo(true));
	}

	@Test
	public void testIPsInWhitelist() {
		assertThat(ip4whitelist.contains("192.168.30.1"), equalTo(true));

		assertThat(ip4whitelist.contains("192.168.30.254"), equalTo(true));

		assertThat(ip4whitelist.contains("10.0.2.1"), equalTo(true));

		assertThat(ip4whitelist.contains("10.0.2.254"), equalTo(true));

		assertThat(ip4whitelist.contains("10.0.1.1"), equalTo(true));

		assertThat(ip4whitelist.contains("10.0.254.254"), equalTo(true));
	}

	@Test
	public void testIPsNotInWhitelist() {
		assertThat(ip4whitelist.contains("192.168.31.1"), equalTo(false));

		assertThat(ip4whitelist.contains("192.167.30.1"), equalTo(false));

		assertThat(ip4whitelist.contains("10.1.1.1"), equalTo(false));

		assertThat(ip4whitelist.contains("10.10.1.1"), equalTo(false));
	}

	/* IPv6 Testing */

	@Test
	public void testIPv6AddressInWhitelist() {
		assertThat(ip6whitelist.contains("::1"), equalTo(true));
	}

	@Test
	public void testIPv6AddressInWhitelistInSubnet() {
		assertThat(ip6whitelist.contains("2001::"), equalTo(true));

		assertThat(ip6whitelist.contains("2001:0:0:0:0:0:0:1"), equalTo(true));

		assertThat(ip6whitelist.contains("2001:0:0:0:0:0:1:1"), equalTo(true));

		assertThat(ip6whitelist.contains("2001:0:0:0:a:a:a:a"), equalTo(true));

		assertThat(ip6whitelist.contains("2001:0:0:0:ffff:ffff:ffff:ffff"), equalTo(true));
	}

	@Test
	public void testIpv6AddressNotInWhitelist() {
		assertThat(ip6whitelist.contains("2001:1:0:0:0:0:0:0"), equalTo(false));

		assertThat(ip6whitelist.contains("2001:0:1::"), equalTo(false));

		assertThat(ip6whitelist.contains("2002:0:0:0:0:0:0:1"), equalTo(false));

		assertThat(ip6whitelist.contains("2001:0:0:1:ffff:ffff:ffff:ffff"), equalTo(false));
	}

	@Test
	public void testWhitelistCreationLeafFirst() throws IOException, JsonUtilsException, NetworkNodeException {
		final ObjectMapper mapper = new ObjectMapper();

		ip4whitelist.init(mapper.readTree("[\"10.0.2.0/24\", \"10.0.0.0/16\"]"));

		assertThat(ip4whitelist.contains("10.0.2.1"), equalTo(true));

		assertThat(ip4whitelist.contains("10.0.10.1"), equalTo(true));
	}

	@Test
	public void testWhitelistCreationParentFirst() throws IOException, JsonUtilsException, NetworkNodeException {
		final ObjectMapper mapper = new ObjectMapper();

		ip4whitelist.init(mapper.readTree("[\"10.0.0.0/16\"], \"10.0.2.0/24\""));

		assertThat(ip4whitelist.contains("10.0.2.1"), equalTo(true));

		assertThat(ip4whitelist.contains("10.0.10.1"), equalTo(true));
	}

	/* IPv4 validation */

	@Test(expected = IOException.class)
	public void badIPv4Input1() throws IOException, JsonUtilsException, NetworkNodeException {
		final ObjectMapper mapper = new ObjectMapper();
		AnonymousIpWhitelist badlist = new AnonymousIpWhitelist();
		badlist.init(mapper.readTree("[\"\"192.168.1/24\"]"));
		assertThat(badlist.contains("192.168.0.1"), equalTo(false));
	}

	@Test(expected = IOException.class)
	public void badIPv4Input2() throws IOException, JsonUtilsException, NetworkNodeException {
		final ObjectMapper mapper = new ObjectMapper();
		AnonymousIpWhitelist badlist = new AnonymousIpWhitelist();
		badlist.init(mapper.readTree("[\"\"256.168.0.1/24\"]"));
		assertThat(badlist.contains("192.168.0.1"), equalTo(false));
	}

	@Test(expected = IOException.class)
	public void badNetmaskInput1() throws IOException, JsonUtilsException, NetworkNodeException {
		final ObjectMapper mapper = new ObjectMapper();
		AnonymousIpWhitelist badlist = new AnonymousIpWhitelist();
		badlist.init(mapper.readTree("[\"\"192.168.0.1/33\"]"));
		assertThat(badlist.contains("192.168.0.1"), equalTo(false));
	}

	@Test(expected = IOException.class)
	public void badNetmaskInput2() throws IOException, JsonUtilsException, NetworkNodeException {
		final ObjectMapper mapper = new ObjectMapper();
		AnonymousIpWhitelist badlist = new AnonymousIpWhitelist();
		badlist.init(mapper.readTree("[\"\"::1/129\"]"));
		assertThat(badlist.contains("::1"), equalTo(false));
	}

	@Test(expected = IOException.class)
	public void badNetmaskInput3() throws IOException, JsonUtilsException, NetworkNodeException {
		final ObjectMapper mapper = new ObjectMapper();
		AnonymousIpWhitelist badlist = new AnonymousIpWhitelist();
		badlist.init(mapper.readTree("[\"\"192.168.0.1/-1\"]"));
		assertThat(badlist.contains("192.168.0.1"), equalTo(false));
	}

	@Test(expected = IOException.class)
	public void validIPv4Input() throws IOException, JsonUtilsException, NetworkNodeException {
		final ObjectMapper mapper = new ObjectMapper();
		AnonymousIpWhitelist badlist = new AnonymousIpWhitelist();
		badlist.init(mapper.readTree("[\"\"192.168.0.1/32\"]"));
		assertThat(badlist.contains("192.168.0.1"), equalTo(false));
	}

	@Test(expected = IOException.class)
	public void validIPv6Input() throws IOException, JsonUtilsException, NetworkNodeException {
		final ObjectMapper mapper = new ObjectMapper();
		AnonymousIpWhitelist badlist = new AnonymousIpWhitelist();
		badlist.init(mapper.readTree("[\"\"::1/128\"]"));
		assertThat(badlist.contains("::1"), equalTo(false));
	}

	/* NetworkNode takes forever to create Tree - commented out until it is needed
	@Test
	public void testAnonymousIpWhitelistPerformance65000() throws NetworkNodeException {
		AnonymousIpWhitelist whitelist = new AnonymousIpWhitelist();
		List<String> tempList = new ArrayList<>();
		// add a bunch of ips to the whitelist

		for (int i = 0; i < 255; i++) {
			for (int j = 0; j < 255; j++) {
				int a = ThreadLocalRandom.current().nextInt(1, 254 + 1);
				int b = ThreadLocalRandom.current().nextInt(1, 254 + 1);
				int c = ThreadLocalRandom.current().nextInt(1, 254 + 1);
				int d = ThreadLocalRandom.current().nextInt(1, 254 + 1);
				tempList.add(String.format("%s.%s.%s.%s", a, b, c, d));
			}
		}

		long startTime = System.nanoTime();

		for (int i = 0; i < tempList.size(); i++) {
			whitelist.add(tempList.get(i) + "/32");
		}

		long durationTime = System.nanoTime() - startTime;

		System.out.println(String.format("Anonymous IP Whitelist creation took %s nanoseconds to create tree of %d subnets", Long.toString(durationTime),
				tempList.size()));

		int total = 1000;

		long start = System.nanoTime();

		for (int i = 0; i <= total; i++) {
			whitelist.contains("192.168.30.1");
		}

		long duration = System.nanoTime() - start;

		System.out.println(
				String.format("Anonymous IP Whitelist average lookup took %s nanoseconds for %d ips", Long.toString(duration / total), tempList.size()));
	}
	*/
	@Test
	public void testAddSubnets() throws NetworkNodeException {
		AnonymousIpWhitelist whitelist = new AnonymousIpWhitelist();

		whitelist.add("192.168.1.1/32");
		assertThat(whitelist.contains("192.168.1.1"), equalTo(true));

		whitelist.add("192.168.1.0/24");
		assertThat(whitelist.contains("192.168.1.255"), equalTo(true));
		assertThat(whitelist.contains("192.168.1.167"), equalTo(true));

		whitelist.add("192.168.1.0/27");
		assertThat(whitelist.contains("192.168.1.255"), equalTo(true));
		assertThat(whitelist.contains("192.168.1.167"), equalTo(true));

		whitelist.add("10.0.0.1/32");
		assertThat(whitelist.contains("10.0.0.1"), equalTo(true));
		assertThat(whitelist.contains("10.0.0.2"), equalTo(false));
		assertThat(whitelist.contains("192.168.2.1"), equalTo(false));
		assertThat(whitelist.contains("192.168.2.255"), equalTo(false));
		assertThat(whitelist.contains("192.167.1.1"), equalTo(false));
		assertThat(whitelist.contains("192.169.1.1"), equalTo(false));
		assertThat(whitelist.contains("10.0.0.0"), equalTo(false));
	}
}
