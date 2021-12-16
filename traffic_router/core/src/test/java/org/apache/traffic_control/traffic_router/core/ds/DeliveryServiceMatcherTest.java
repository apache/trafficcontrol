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

package org.apache.traffic_control.traffic_router.core.ds;

import org.apache.traffic_control.traffic_router.core.request.HTTPRequest;
import org.apache.traffic_control.traffic_router.core.request.Request;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.greaterThan;
import static org.hamcrest.Matchers.lessThan;
import static org.junit.Assert.assertThat;
import static org.mockito.Mockito.mock;
import static org.apache.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.HOST;
import static org.apache.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.PATH;
import static org.mockito.Mockito.when;

@PrepareForTest(DeliveryService.class)
@RunWith(PowerMockRunner.class)
@PowerMockIgnore("javax.management.*")
public class DeliveryServiceMatcherTest {
	@Test
	public void itReturnsTrueWhenAllMatchersPass() {
		DeliveryServiceMatcher deliveryServiceMatcher = new DeliveryServiceMatcher(mock(DeliveryService.class));
		deliveryServiceMatcher.addMatch(HOST, ".*\\.service01-kabletown.com\\..*", "");
		deliveryServiceMatcher.addMatch(PATH, ".*abcd.*", "");

		HTTPRequest httpRequest = new HTTPRequest();
		httpRequest.setHostname("foo.service01-kabletown.com.bar");
		httpRequest.setPath("foo/abcd/bar");
		assertThat(deliveryServiceMatcher.matches(httpRequest), equalTo(true));
	}

	@Test
	public void itReturnsFalseWhenAnyMatcherFails() {
		DeliveryServiceMatcher deliveryServiceMatcher = new DeliveryServiceMatcher(mock(DeliveryService.class));
		deliveryServiceMatcher.addMatch(HOST, ".*\\.service01-kabletown.com\\..*", "");
		deliveryServiceMatcher.addMatch(PATH, ".*abcd.*", "");

		HTTPRequest httpRequest = new HTTPRequest();
		httpRequest.setHostname("foo.serviceZZ-kabletown.com.bar");
		httpRequest.setPath("foo/abcd/bar");
		assertThat(deliveryServiceMatcher.matches(httpRequest), equalTo(false));
	}

	@Test
	public void itReturnsFalseWhenItHasNoMatchers() {
		DeliveryServiceMatcher deliveryServiceMatcher = new DeliveryServiceMatcher(mock(DeliveryService.class));

		Request request = new Request();
		assertThat(deliveryServiceMatcher.matches(request), equalTo(false));

		HTTPRequest httpRequest = new HTTPRequest();
		assertThat(deliveryServiceMatcher.matches(httpRequest), equalTo(false));
	}

	@Test
	public void itComparesByMatchRegexes() {
		DeliveryService deliveryService = mock(DeliveryService.class);

		DeliveryServiceMatcher deliveryServiceMatcher1 = new DeliveryServiceMatcher(deliveryService);
		deliveryServiceMatcher1.addMatch(HOST, ".*\\.service01-kabletown.com\\..*", "");
		deliveryServiceMatcher1.addMatch(PATH, ".*abc.*", "");

		DeliveryServiceMatcher deliveryServiceMatcher1a = new DeliveryServiceMatcher(deliveryService);
		deliveryServiceMatcher1a.addMatch(HOST, ".*\\.service01-kabletown.com\\..*", "");
		deliveryServiceMatcher1a.addMatch(PATH, ".*abc.*", "");

		DeliveryServiceMatcher deliveryServiceMatcher2 = new DeliveryServiceMatcher(deliveryService);
		deliveryServiceMatcher2.addMatch(HOST, ".*\\.service01-kabletown.com\\..*", "");
		deliveryServiceMatcher2.addMatch(PATH, ".*abcde.*", "");

		assertThat(deliveryServiceMatcher1.equals(deliveryServiceMatcher1a), equalTo(true));
		assertThat(deliveryServiceMatcher1a.equals(deliveryServiceMatcher1), equalTo(true));

		assertThat(deliveryServiceMatcher1.compareTo(deliveryServiceMatcher1), equalTo(0));
		assertThat(deliveryServiceMatcher1.compareTo(deliveryServiceMatcher1a), equalTo(0));

		assertThat(deliveryServiceMatcher1.compareTo(deliveryServiceMatcher2), greaterThan(0));
		assertThat(deliveryServiceMatcher2.compareTo(deliveryServiceMatcher1), lessThan(0));
	}

	@Test
	public void itHandlesMatcherWithoutRequestMatchers() {
		DeliveryService deliveryService = mock(DeliveryService.class);

		DeliveryServiceMatcher deliveryServiceMatcher1 = new DeliveryServiceMatcher(deliveryService);
		deliveryServiceMatcher1.addMatch(HOST, ".*\\.service01-kabletown.com\\..*", "");
		deliveryServiceMatcher1.addMatch(PATH, ".*abc.*", "");

		DeliveryServiceMatcher deliveryServiceMatcher2 = new DeliveryServiceMatcher(deliveryService);

		assertThat(deliveryServiceMatcher1.compareTo(deliveryServiceMatcher2), equalTo(-1));
		assertThat(deliveryServiceMatcher2.compareTo(deliveryServiceMatcher1), equalTo(1));
	}

	@Test
	public void compareToReturns0WhenSameMatchersDifferentDeliveryServices() {
		DeliveryService deliveryService1 = mock(DeliveryService.class);
		when(deliveryService1.getId()).thenReturn("delivery service 1");

		DeliveryService deliveryService2 = mock(DeliveryService.class);
		when(deliveryService2.getId()).thenReturn("delivery service 2");

		DeliveryServiceMatcher deliveryServiceMatcher1 = new DeliveryServiceMatcher(deliveryService1);
		deliveryServiceMatcher1.addMatch(HOST, ".*\\.service01-kabletown.com\\..*", "");
		deliveryServiceMatcher1.addMatch(PATH, ".*abc.*", "");

		DeliveryServiceMatcher deliveryServiceMatcher2 = new DeliveryServiceMatcher(deliveryService2);
		deliveryServiceMatcher2.addMatch(HOST, ".*\\.service01-kabletown.com\\..*", "");
		deliveryServiceMatcher2.addMatch(PATH, ".*abc.*", "");

		assertThat(deliveryServiceMatcher1.equals(deliveryServiceMatcher2), equalTo(false));
		assertThat(deliveryServiceMatcher2.equals(deliveryServiceMatcher1), equalTo(false));

		assertThat(deliveryServiceMatcher1.compareTo(deliveryServiceMatcher2), equalTo(0));
	}
}
