package com.comcast.cdn.traffic_control.traffic_router.core.ds;

import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.greaterThan;
import static org.hamcrest.Matchers.lessThan;
import static org.junit.Assert.assertThat;
import static org.mockito.Mockito.mock;
import static com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.HOST;
import static com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.PATH;

@PrepareForTest(DeliveryService.class)
@RunWith(PowerMockRunner.class)
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
}
